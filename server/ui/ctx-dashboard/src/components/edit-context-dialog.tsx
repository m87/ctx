import { useState } from "react";
import { Button } from "./ui/button";
import { Dialog, DialogClose, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from "./ui/dialog";
import { Input } from "./ui/input";
import { Label } from "./ui/label";
import { Context } from "@/api/api";
import { set } from "date-fns";




export function EditContextDialog({ children, context, onChange }: { children: any, context: Context, onChange: (context: Context) => void }) {
    const [open, setOpen] = useState(false)

    return <Dialog open={open} onOpenChange={setOpen}>
        <DialogTrigger asChild>
            {children}
        </DialogTrigger>
        <DialogContent className="sm:max-w-[425px]">
            <DialogHeader>
                <DialogTitle>Edit context</DialogTitle>
                <DialogDescription>
                </DialogDescription>
            </DialogHeader>
            <div className="grid gap-4">
                <div className="grid gap-3">
                    <Label htmlFor="name-1">Name</Label>
                    <Input id="name-1" name="name" defaultValue={context.description} onChange={(e) => context.description = e.currentTarget.value} />
                </div>
            </div>
            <DialogFooter>
                <DialogClose asChild>
                    <Button variant="outline">Cancel</Button>
                </DialogClose>
                <Button type="button" onClick={() => { onChange(context); setOpen(false); }}>Save changes</Button>
            </DialogFooter>
        </DialogContent>
    </Dialog>

}