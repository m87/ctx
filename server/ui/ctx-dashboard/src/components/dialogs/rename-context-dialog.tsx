import { api, Context } from "@/api/api";
import { ReactNode, useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { Dialog, DialogClose, DialogContent, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from "../ui/dialog";
import { Label } from "../ui/label";
import { Button } from "../ui/button";
import { Input } from "../ui/input";


export function RenameContextDialog({ context, children }: { context: Context, children: ReactNode }) {
    const [open, setOpen] = useState(false)
    const qc = useQueryClient();
    const renameMutation = useMutation(api.context.renameMutation(qc))
    const [name, setName] = useState(context.description);

    const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        renameMutation.mutate({ ctxId: context.id, name })
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
                            <DialogTitle>Rename context</DialogTitle>
                        </DialogHeader>
                    </div>

                    <div className="p-6 pt-0 overflow-auto">
                        <div className="grid gap-4">
                            <div className="grid gap-2">
                                <Label htmlFor="name">Name</Label>
                                <Input id="name" value={name} onChange={(e) => setName(e.currentTarget.value)} />
                            </div>
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