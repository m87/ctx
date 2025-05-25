import { CheckIcon, EditIcon } from "lucide-react";
import { useState } from "react";
import { Input } from "./ui/input";
import { Interval } from "@/api/api";


export function IntervalComponent({ interval }: Readonly<{ interval: Interval }>) {
    const [edited, setEdited] = useState(false);
    return (
        <div className="flex flex-col">
            {edited && <div className="flex gap-2 p-2 justify-between">
                <div>-</div>
                <div>
                    ({interval.duration / 1000000000} seconds)
                </div>
                <div><CheckIcon className="cursor-pointer" onClick={() => setEdited(false)}></CheckIcon></div>
            </div>}
            {!edited && <div className="flex gap-2 p-2">
                <div>
                    {interval.start.toString()}
                </div>
                <div>-</div>
                <div>
                    {interval.end.toString()}
                </div>
                <div>
                    ({interval.duration / 1000000000} seconds)
                </div>
                <div><EditIcon className="cursor-pointer" onClick={() => setEdited(true)}></EditIcon></div>
            </div>}
        </div>
    );
}

export default IntervalComponent;