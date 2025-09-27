import { api, Context } from "@/api/api";
import { PlusIcon } from "lucide-react";
import { useEffect, useRef, useState } from "react";
import ContextCard from "./context-card";
import { Input } from "./ui/input";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useParams } from "react-router-dom";
import { ScrollArea } from "./ui/scroll-area";
import { Checkbox } from "./ui/checkbox";
import { Label } from "./ui/label";
import { format } from "date-fns";
import { DaySummary } from "@/api/api-summary";

export interface CardsProps {
    term: string
    expandId?: string
    selectedDate: Date
}

function concatContexts(summary: DaySummary | undefined) {
    return summary ? [...summary.contexts, ...summary.otherContexts] : []
}

export function SectionCards({ term, expandId, selectedDate }: CardsProps) {
    const [searchTerm, setSearchTerm] = useState(term);
    const [showAllContexts, setShowAllContexts] = useState(false);
    const { data: summary } = useQuery({ ...api.summary.daySummaryQuery(format(selectedDate, "yyyy-MM-dd"), showAllContexts) });
    const filteredList = concatContexts(summary).filter((context) =>
        context.description.toLowerCase().includes(searchTerm.toLowerCase())
    );
    const queryClient = useQueryClient();
    const createAndSwitchMutation = useMutation(api.context.createAndSwitchMutation(queryClient))
    const { day } = useParams();

    useEffect(() => {
        setSearchTerm(term);
    }, [term])

    return (<div className="flex flex-col h-screen">
        <div className="pt-3 pr-6 pl-6 flex items-start">
            <div className="flex flex-col w-full">
                <div className="w-full flex">                
                    <Input type="text" value={searchTerm} onChange={(e) => setSearchTerm(e.target.value)} className="rounded-b-none"
                    onKeyDown={(e) => {
                        if (e.key === 'Enter' && searchTerm.trim() !== '' && filteredList && filteredList?.length > 0) {
                            api.context.switch(filteredList[0].id);
                            setSearchTerm('');
                        }

                        if (e.key === 'Enter' && filteredList?.length === 0) {
                            createAndSwitchMutation.mutate({ description: searchTerm, day });
                            setSearchTerm('');
                        }
                    }}
                    placeholder="Search or create new..."></Input>
                </div>

                <div className="flex items-center gap-3 pt-2 pl-1 pr-1 pb-2 mb-3 opacity-70 border border-t-0 rounded-b-md">
                    <Checkbox id="all-ctx" onCheckedChange={(v) => setShowAllContexts(v)}/>
                    <Label htmlFor="all-ctx">Show all contexts</Label>
                </div>
            </div>

        </div>
        <div className="flex-1 min-h-0 mb-44">
            <ScrollArea className="h-full">
                {filteredList?.length > 0 && <div
                    className="*:data-[slot=card]:shadow-xs @xl/main:grid-cols-2 @5xl/main:grid-cols-4 grid grid-cols-1 gap-4 px-4 *:data-[slot=card]:bg-gradient-to-t *:data-[slot=card]:from-primary/5 *:data-[slot=card]:to-card dark:*:data-[slot=card]:bg-card lg:px-6">
                    {filteredList?.map((context) => (
                        <ContextCard key={context.id} context={context} expandCard={expandId === context.id}> </ContextCard>
                    ))}
                </div>}
                {filteredList?.length === 0 && <div className="flex items-center justify-center h-full">
                    <div className="text-muted-foreground">No contexts found</div>
                </div>}
            </ScrollArea>
        </div>
    </div>
    )
} 
