import { CheckIcon, EditIcon } from "lucide-react";
import { useCallback, useState } from "react";
import { Interval, ZonedDateTime } from "@/api/api";
import { DateTimeInput } from "./ui/datetime";
import { DateTime } from "luxon";
import { DataTable } from "./ui/interval-table";


export function IntervalComponent({ interval, onChange }: Readonly<{ interval: Interval, onChange: (id: string, start: ZonedDateTime, end: ZonedDateTime) => void }>) {
  const [edited, setEdited] = useState(false);
  const [start, setStart] = useState(interval.start);
  const [end, setEnd] = useState(interval.end);

  const handleStartChange = useCallback((dt: DateTime) => setStart(ZonedDateTime.fromDateTime(dt)), [])
  const handleEndChange = useCallback((dt: DateTime) => setEnd(ZonedDateTime.fromDateTime(dt)), [])

  return (
    <div className="flex flex-col">
      {edited && <div className="flex gap-2 p-2 items-center">
        <DateTimeInput datetime={start.toDateTime()} editable={true} onChange={handleStartChange}></DateTimeInput>
        <div>-</div>
        <DateTimeInput datetime={end.toDateTime()} editable={true} onChange={handleEndChange}></DateTimeInput>
        <div>
          ({Math.floor(interval.duration / 60000000000) } min)
        </div>
        <div><CheckIcon className="cursor-pointer" onClick={() => {
          onChange(interval.id, start, end)
          setEdited(false);
        }}></CheckIcon></div>
      </div>}
      {!edited && <div className="flex gap-2 p-2 items-center">
        <div>
          <DateTimeInput datetime={start.toDateTime()}></DateTimeInput>
        </div>
        <div>-</div>
        <div>
          <DateTimeInput datetime={end.toDateTime()} ></DateTimeInput>
        </div>
        <div>
          ({Math.ceil(interval.duration / 60000000000) } min)
        </div>
        <div><EditIcon className="cursor-pointer" onClick={() => setEdited(true)}></EditIcon></div>
      </div>}
    </div>
  );
}

export default IntervalComponent;
