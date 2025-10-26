import { api, Context, Interval, ZonedDateTime } from "@/api/api";
import { Children, ReactNode, useCallback, useState } from "react";
import { Dialog, DialogClose, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from "../ui/dialog";
import { Button } from "../ui/button";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useParams } from "react-router-dom";
import EditableBadges from "../ui/labels-editor";


export function LabelsDialog({ context, children }: { context: Context, children: ReactNode }) {
    const [open, setOpen] = useState(false)
    const qc = useQueryClient();
    const editLabelsMutation = useMutation(api.context.editLabelsMutation(qc))
    const { day } = useParams();

    const handleChange = (badges: string[]) => {
        editLabelsMutation.mutate({ id: context.id, labels: badges })
    }

    return <>
        <Dialog open={open} onOpenChange={setOpen}>
            <DialogTrigger asChild>
                {children}
            </DialogTrigger>


            <DialogContent className="grid grid-rows-[auto_1fr_auto] sm:max-w-[520px] p-0">
                <div className="p-6 pb-4">
                    <DialogHeader>
                        <DialogTitle>Labels</DialogTitle>
                    </DialogHeader>
                </div>

                <EditableBadges initBadges={context.labels} onChange={(badges: string[]) => handleChange(badges)}></EditableBadges>
            </DialogContent>
        </Dialog>
    </>
}
