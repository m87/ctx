import { api, Comment, Context, Interval, ZonedDateTime } from "@/api/api";
import { Children, ReactNode, useCallback, useState } from "react";
import { Dialog, DialogClose, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from "../ui/dialog";
import { Button } from "../ui/button";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useParams } from "react-router-dom";


export function DeleteContextCommentDialog({ context, comment, children }: { context: Context, comment: Comment, children: ReactNode }) {
    const [open, setOpen] = useState(false)
    const qc = useQueryClient();
    const deleteMutation = useMutation(api.context.deleteCommentMutation(qc))
    const { day } = useParams();

    const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        deleteMutation.mutate({ ctxId: context.id, commentId: comment.id, day })
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
                            <DialogTitle>Delete comment</DialogTitle>
                            <DialogDescription>Are you sure you want to delete comment?</DialogDescription>
                        </DialogHeader>
                    </div>

                    <DialogFooter className="sm:space-x-0 gap-2 flex-row justify-end border-t bg-background p-4">
                        <DialogClose asChild>
                            <Button type="button" variant="outline">Cancel</Button>
                        </DialogClose>
                        <Button type="submit" variant={"destructive"}>Delete</Button>
                    </DialogFooter>
                </form>
            </DialogContent>
        </Dialog>
    </>
}
