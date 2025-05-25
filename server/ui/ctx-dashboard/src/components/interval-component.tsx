import { CheckIcon, EditIcon } from "lucide-react";
import { useState } from "react";
import { Input } from "./ui/input";
import { api, Interval, ZonedDateTime } from "@/api/api";
import { DateTime, Zone } from "luxon";


export function IntervalComponent({ interval, onChange }: Readonly<{ interval: Interval, onChange: (id: string, start: ZonedDateTime, end: ZonedDateTime) => void }>) {
    const [edited, setEdited] = useState(false);
    const [start, setStart] = useState(interval.start);
    const [end, setEnd] = useState(interval.end);
    return (
        <div className="flex flex-col">
            {edited && <div className="flex gap-2 p-2 justify-between">
                <Input type="datetime-local" value={start.toDateTime().toLocaleString()} onChange={(e) => {
                    setStart(ZonedDateTime.fromDateTime(DateTime.fromISO(e.target.value, { zone: interval.start.timezone ?? "utc" })));
                }}></Input>
                <div>-</div>
                <Input type="datetime-local" value={end.toDateTime().toLocaleString()} onChange={(e) => {
                    setEnd(ZonedDateTime.fromDateTime(DateTime.fromISO(e.target.value, { zone: interval.end.timezone ?? "utc" })));
                }}></Input>
                <div>
                    ({interval.duration / 1000000000} seconds)
                </div>
                <div><CheckIcon className="cursor-pointer" onClick={() => {
                    onChange(interval.id, start, end)
                    setEdited(false);
                }}></CheckIcon></div>
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