import { api, Interval } from "@/api/api";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "./ui/table";
import { Button } from "./ui/button";
import { Edit, Trash } from "lucide-react";
import { Dialog, DialogClose, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from "./ui/dialog";
import { Label } from "./ui/label";
import { DateTimeInput } from "./ui/datetime";
import { Item, ItemActions, ItemContent, ItemDescription, ItemTitle } from "./ui/item";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { durationAsHM } from "@/lib/utils";
import { EditIntervalDialog } from "./edit-interval-dialog";




export function IntervalTable({ ctxId, intervals }: { ctxId: string, intervals: Interval[] }) {

    return <>
        <h4 className="font-medium mb-2">Intervals</h4>
        <div className="rounded-lg border">
            <Table>
                <TableHeader>
                    <TableRow>
                        <TableHead>Start</TableHead>
                        <TableHead>End</TableHead>
                        <TableHead>Summary</TableHead>
                        <TableHead className="w-[1%]"></TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {(intervals ?? []).length === 0 && (
                        <TableRow>
                            <TableCell colSpan={4} className="text-sm text-muted-foreground">No intervals</TableCell>
                        </TableRow>
                    )}
                    {(intervals ?? []).map((interval) => (
                        <TableRow key={interval.id} className="hover:bg-muted/40">
                            <TableCell>{interval.start.toString()}</TableCell>
                            <TableCell>{interval.end.toString()}</TableCell>
                            <TableCell>{durationAsHM(interval.duration)}</TableCell>
                            <TableCell className="text-right">
                                <div className="flex gap-2 justify-end">
                                    <EditIntervalDialog interval={interval} ctxId={ctxId}>
                                        <Button variant="ghost"><Edit /></Button>
                                    </EditIntervalDialog>
                                </div>
                            </TableCell>
                        </TableRow>
                    ))}
                </TableBody>
            </Table>
        </div>
    </>
}