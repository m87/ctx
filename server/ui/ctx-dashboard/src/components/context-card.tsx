import { ArrowDown, ChevronDown, ChevronUp, PlayCircleIcon } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "./ui/card";
import { useState } from "react";
import { api } from "@/api/api";
import IntervalComponent from "./interval-component";


export function ContextCard({ context }) {
    const [hovered, setHovered] = useState(false);
    const [expanded, setExpand] = useState(false);
    const cardClick = (id: string) => {
        api.context.switch(id)
    };
    return (
        <Card key={context.id} className="@container/card"
            onMouseEnter={() => setHovered(true)}
            onMouseLeave={() => setHovered(false)}
        >
            <CardHeader className="relative">
                <CardTitle className="@[250px]/card:text-3xl text-2xl font-semibold tabular-nums flex justify-between w-full items-center">
                    <div className="flex w-full items-center">
                        {hovered && <div className="cursor-pointer"><PlayCircleIcon size={30} onClick={() => cardClick(context.id)} /></div>}
                        <div>    {context.description} </div>
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
                        <IntervalComponent key={interval.id} interval={interval} />
                    ))}
                </div>
            </CardContent>
            }
        </Card>
    )
}

export default ContextCard;