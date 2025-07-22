import {IntervalEntry, IntervalsResponse, IntervalsResponseEntry} from "@/api/api-intervals";
import {colorHash} from "@/lib/utils";
import {clsx} from "clsx";
import {useEffect, useRef, useState} from "react";
import {
    ContextMenu,
    ContextMenuContent,
    ContextMenuItem,
    ContextMenuSub,
    ContextMenuSubContent,
    ContextMenuSubTrigger,
    ContextMenuTrigger
} from "@/components/ui/context-menu";
import { Search} from "lucide-react";
import {Input} from "@/components/ui/input";
import {api} from "@/api/api";
import { Calendar } from "./ui/calendar";
import { Button } from "./ui/button";
import TimelineSplit from "./timeline-split";
import { Interval } from "luxon";


export function intervalsResponseAsTimelineData(data: IntervalsResponse): Record<string, TimeInterval[]> {
    const output = {}

    data.days.filter((entry: IntervalsResponseEntry) => entry.intervals).map((entry: IntervalsResponseEntry) => output[entry.date] = entry.intervals?.map((interval: IntervalEntry) => ({
        id: interval.id,
        start: interval.interval.start.time?.split("T")[1].substring(0, 8),
        end: interval.interval.end.time?.split("T")[1].substring(0, 8),
        color: `${colorHash(interval.ctxId)}`,
        ctxId: interval.ctxId,
        description: interval.description
    })))

    return output;
}


export interface TimeInterval {
    id: string;
    start: string;
    end: string;
    color?: string;
    ctxId: string;
    description?: string;
}

export interface TimelineProps {
    data: Record<string, TimeInterval[]>;
    hideDates: boolean;
    hideGuides: boolean;
    onItemSelect: (interval: TimeInterval | null) => void;
    ctxNames: {description: string, id: string}[];
}

function timeToDecimal(time: string): number {
    const [h, m, s] = time.split(":").map(Number);
    return (h * 3600 + m * 60 + s) / 86400;
}

function Timeline({data, hideDates, hideGuides, onItemSelect, ctxNames}: TimelineProps) {
    const hours = Array.from({length: 24}, (_, i) => i);
    const dates = Object.keys(data);
    const boxRefs = useRef<Map<string, HTMLDivElement>>(new Map());
    const BLOCK = 100;
    const [ctxNameSearchTerm, setCtxNameSearchTerm] = useState('')

    useEffect(() => {
        function handleClick(event: MouseEvent) {
            const target = event.target as Node;
            const clickedOnBox = Array.from(boxRefs.current.values()).some((boxEl) => boxEl.contains(target));

            if (!clickedOnBox) {
                //       setSelected('')
                //       onItemSelect(null)
            }
        }

        function handleKeyDown(event: KeyboardEvent) {
            if (event.key === 'Escape') {
                setSelected('')
                onItemSelect(null)
            }
        }

        document.addEventListener('mousedown', handleClick);
        document.addEventListener('keydown', handleKeyDown);

        return () => {
            document.removeEventListener('mousedown', handleClick);
            document.removeEventListener('keydown', handleKeyDown);
        }

    }, []);

    const [selected, setSelected] = useState('')
    return (
        <div className="font-sans bg-gray-100">
            <div className="relative p-4 overflow-x-auto">
                <div className="min-w-[1000px]">
                    <div className={clsx("relative h-6 mb-2", hideDates ? "" : "ml-24")}>
                        {hours.map((hour) => (
                            <div
                                key={`hour-${hour}`}
                                className="absolute top-0 text-xs font-semibold text-center border-l border-black/30"
                                style={{
                                    left: `${(hour / 24) * BLOCK}%`,
                                    width: `${BLOCK / 24}%`,
                                }}
                            >
                                {hour}
                            </div>
                        ))}
                    </div>

                    {dates.map((date, rowIndex) => (
                        <div key={`row-${rowIndex}-${date}`}
                             className={clsx(rowIndex === 0 && !hideGuides ? "border-t" : "", !hideGuides ? "border-b" : "", "pt-1 pb-1 relative w-full h-10 flex items-center")}>
                            {!hideDates &&
                                <div className="w-24 pr-2 text-sm font-bold text-right text-gray-700 flex">
                                    {date}
                                </div>
                            }
                            <div className="relative flex-1 h-full">
                                {data[date].map((interval, idx) => {
                                    const start = timeToDecimal(interval.start);
                                    const end = timeToDecimal(interval.end);
                                    const left = `${start * BLOCK}%`;
                                    const width = `${(end - start) * BLOCK}%`;

                                    return (
                                        <ContextMenu>
                                            <ContextMenuTrigger>
                                                <div
                                                    onClick={() => {
                                                        if (selected === interval.ctxId) {
                                                            setSelected('')
                                                            onItemSelect(null);
                                                        } else {
                                                            setSelected(interval.ctxId);
                                                            onItemSelect(interval);
                                                        }
                                                    }}
                                                    ref={(el) => el && boxRefs.current.set(interval.ctxId, el)}
                                                    key={`interval-${rowIndex}-${idx}`}
                                                    className={`${selected && interval.ctxId !== selected ? 'opacity-50' : ''} cursor-pointer absolute top-0 h-full rounded text-xs text-white flex items-center justify-center px-1 text-ellipsis overflow-hidden whitespace-nowrap`}
                                                    style={{left, width, backgroundColor: interval.color}}
                                                >
                                                </div>
                                            </ContextMenuTrigger>
                                            <ContextMenuContent>
                                                <ContextMenuItem>Edit</ContextMenuItem>
                                                  <ContextMenuItem onClick={() => api.intervals.delete(interval.ctxId, interval.id)}>Delete</ContextMenuItem>
                                              <ContextMenuSub>
                                                    <ContextMenuSubTrigger inset>Split</ContextMenuSubTrigger>
                                                    <ContextMenuSubContent className="w-44">
                                                    <TimelineSplit interval={interval}></TimelineSplit>
                                                    </ContextMenuSubContent>
                                                </ContextMenuSub>
                                                <ContextMenuSub>
                                                    <ContextMenuSubTrigger inset>Move to..</ContextMenuSubTrigger>
                                                    <ContextMenuSubContent className="w-44">
                                                        <div className="flex">
                                                            <Input className="w-44"
                                                                   onChange={(e) => setCtxNameSearchTerm(e.target.value)}
                                                                   value={ctxNameSearchTerm}></Input><Search></Search>
                                                        </div>
                                                        {
                                                            ctxNames.filter((ctx) => ctx.description.startsWith(ctxNameSearchTerm)).map((ctx) => {
                                                                return (
                                                                    <ContextMenuItem key={ctx.id} onClick={() => api.intervals.move({
                                                                        src: interval.ctxId,
                                                                        target: ctx.id,
                                                                        id: interval.id
                                                                    })}>{ctx.description}</ContextMenuItem>)
                                                            })
                                                        }
                                                    </ContextMenuSubContent>
                                                </ContextMenuSub>
                                            </ContextMenuContent>
                                        </ContextMenu>
                                    );
                                })}
                            </div>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
}

export default Timeline;
