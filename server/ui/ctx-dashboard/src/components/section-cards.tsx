import { api } from "@/api/api";
import {
  Card,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { useQuery } from "@tanstack/react-query";
import { AnchorIcon, ArrowDown, ArrowDown01Icon, PlayCircleIcon, PlusIcon } from "lucide-react";
import { useState } from "react";
import ContextCard from "./context-card";
import { Input } from "./ui/input";
import { ScrollArea } from "@radix-ui/react-scroll-area";

export function SectionCards() {

  const { data: contextList } = useQuery({ ...api.context.listQuery });

  const [searchTerm, setSearchTerm] = useState('');
  const filteredList = contextList?.filter((context) =>
    context.description.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const createNewContext = (description: string) => {
    api.context.createAndSwitch(description);
  };

  const cardClick = (id) => {
    api.context.switch(id)
  };
  return (<div>
    <div className="pt-3 pb-2 pr-6 pl-6 flex items-center">
      <Input type="text" value={searchTerm} onChange={(e) => setSearchTerm(e.target.value)} onKeyDown={(e) => {
        if (e.key === 'Enter' && searchTerm.trim() !== '' && filteredList && filteredList?.length > 0) {
          api.context.switch(filteredList[0].id);
          setSearchTerm('');
        }
      
        if (e.key === 'Enter' && filteredList?.length === 0) {
          createNewContext(searchTerm);
          setSearchTerm('');
        }
      }}
        placeholder="Search or create new..."></Input>
      {filteredList?.length == 0 && <div>
        <PlusIcon className="cursor-pointer" onClick={() => {
          if (searchTerm.trim() !== '') {
          createNewContext(searchTerm)
          setSearchTerm('');
          }}}></PlusIcon>
      </div>}
    </div>
    <ScrollArea className="h-full flex-2 overflow-auto">
      <div className="*:data-[slot=card]:shadow-xs @xl/main:grid-cols-2 @5xl/main:grid-cols-4 grid grid-cols-1 gap-4 px-4 *:data-[slot=card]:bg-gradient-to-t *:data-[slot=card]:from-primary/5 *:data-[slot=card]:to-card dark:*:data-[slot=card]:bg-card lg:px-6">
        {filteredList?.map((context) => (
          <ContextCard key={context.id} context={context} > </ContextCard>
        ))}
      </div>
    </ScrollArea>
  </div>
  )
}
