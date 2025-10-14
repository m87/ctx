import { api, Interval, ZonedDateTime } from "@/api/api";
import { Children, ReactNode, useCallback, useState } from "react";
import { Dialog, DialogClose, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from "./ui/dialog";
import { Label } from "./ui/label";
import { DateTimeInput } from "./ui/datetime";
import { Item, ItemActions, ItemContent, ItemDescription, ItemTitle } from "./ui/item";
import { Button } from "./ui/button";
import { Trash } from "lucide-react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { DateTime } from "luxon";
import { useParams } from "react-router-dom";


export function EditIntervalDialog({ ctxId, interval, children }: { ctxId: string, interval: Interval, children: ReactNode }) {
    const [open, setOpen] = useState(false)
    const qc = useQueryClient();
    const deleteMutation = useMutation(api.intervals.deleteMutation(qc))
    const [start, setStart] = useState(interval.start);
    const [end, setEnd] = useState(interval.end);
    const updateIntervalMutation = useMutation(api.context.updateIntervalMutation(qc))
    const { day } = useParams();


    const handleStartChange = useCallback((dt: DateTime) => setStart(ZonedDateTime.fromDateTime(dt)), [])
    const handleEndChange = useCallback((dt: DateTime) => setEnd(ZonedDateTime.fromDateTime(dt)), [])

    const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        updateIntervalMutation.mutate({ contextId: ctxId, intervalId: interval.id, start, end, day });
        setOpen(false);
    }

    const deleteInterval = (interval: Interval) => {
        console.log(interval)
        deleteMutation.mutate({ ctxId: ctxId, id: interval.id })
        setOpen(false);
    }

    return <>
        <Dialog open={open} onOpenChange={setOpen}>
            <DialogTrigger asChild>
                {children}
            </DialogTrigger>


            <DialogContent className="grid grid-rows-[auto_1fr_auto] sm:max-w-[520px] p-0">
                <form className="contents" onSubmit={handleSubmit}>
                    <div className="p-6 pb-4">
                        <DialogHeader>
                            <DialogTitle>Edit interval</DialogTitle>
                            <DialogDescription>Adjust start and end times.</DialogDescription>
                        </DialogHeader>
                    </div>

                    <div className="p-6 pt-0 overflow-auto">
                        <div className="grid gap-4 sm:grid-cols-2">
                            <div className="grid gap-2">
                                <Label htmlFor="start">Start</Label>
                                <DateTimeInput id="start" datetime={interval.start.toDateTime()} onChange={handleStartChange} editable />
                            </div>
                            <div className="grid gap-2">
                                <Label htmlFor="end">End</Label>
                                <DateTimeInput id="end" datetime={interval.end.toDateTime()} onChange={handleEndChange} editable/>
                            </div>
                        </div>

                        <div className="mt-6">
                            <Item variant="outline" className="border-destructive/30 bg-destructive/5">
                                <ItemContent>
                                    <ItemTitle className="text-destructive">Delete interval</ItemTitle>
                                    <ItemDescription className="text-muted-foreground"> This action is irreversible </ItemDescription>
                                </ItemContent>
                                <ItemActions>
                                    <Button type="button" variant="destructive" size="sm" onClick={() => deleteInterval(interval)}>
                                        <Trash></Trash>
                                    </Button>
                                </ItemActions>
                            </Item>
                        </div>
                    </div>

                    <DialogFooter className="sm:space-x-0 gap-2 flex-row justify-end border-t bg-background p-4">
                        <DialogClose asChild>
                            <Button type="button" variant="outline">Cancel</Button>
                        </DialogClose>
                        <Button type="submit">Save changes</Button>
                    </DialogFooter>
                </form>
            </DialogContent>
        </Dialog>
    </>
}