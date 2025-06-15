import { ArrowDown, ChevronDown, ChevronUp, PlayCircleIcon } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "./ui/card";
import { useState } from "react";
import { api, ZonedDateTime } from "@/api/api";
import IntervalComponent from "./interval-component";
import { Badge } from "./ui/badge";


export function ContextCard({ context }) {
    const [hovered, setHovered] = useState(false);
    const [expanded, setExpand] = useState(false);
    const cardClick = (id: string) => {
        api.context.switch(id)
    };

    const updateInterval = (id: string, start: ZonedDateTime, end: ZonedDateTime) => {
        api.context.updateInterval(context.id, id, start, end);
    };

    return (
        <Card key={context.id} className="@container/card"
            onMouseEnter={() => setHovered(true)}
            onMouseLeave={() => setHovered(false)}
        >
            <CardHeader className="relative">
                <CardTitle className="@[250px]/card:text-3xl text-2xl font-semibold tabular-nums flex justify-between w-full items-center">
                    <div className="flex w-full items-center ">
                        {hovered && <div className="cursor-pointer"><PlayCircleIcon size={30} onClick={() => cardClick(context.id)} /></div>}
                        <div className="flex flex-col items-start">
                            <div>    {context.description} </div>
                            <div className="flex">
                                {context.labels?.map((label) => (
                                    <Badge variant={"secondary"}>{label}</Badge>
                                ))}
                            </div>
                        </div>
                    </div>
                    <div>
                        <div onClick={() => setExpand(!expanded)} className="cursor-pointer">
                            <ChevronDown
                                className={`transition-transform duration-200 ${expanded ? "rotate-180" : ""}`}
                            />
                        </div>
                    </div>
                </CardTitle>
            </CardHeader>
            {expanded && <CardContent className="flex flex-col gap-2">
                <div className="flex flex-col justify-center">
                    {context.intervals?.map((interval) => (
                        <IntervalComponent key={interval.id} interval={interval} onChange={updateInterval} />
                    ))}
                </div>
            </CardContent>
            }
        </Card>
    )
}

export default ContextCard;