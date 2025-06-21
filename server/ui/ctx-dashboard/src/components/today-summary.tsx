import { useQuery } from "@tanstack/react-query";
import {api} from "@/api/api";
import { SectionCards } from "./section-cards";
import Timeline, {TimeInterval} from "@/components/timeline";
import {IntervalsResponse} from "@/api/api-intervals";
import {colorHash} from "@/lib/utils"
import { useState } from "react";

export function TodaySummary() {
    const [selectedInterval, setSelectedInterval] = useState(null)
    const {data: summary} = useQuery({...api.summary.todaySummaryQuery});
    const {data: intervals} = useQuery({...api.intervals.intervalsQuery, select: (data: IntervalsResponse) => (
        {"2025-06-10": data.intervals.map(interval => ({
            start: interval.interval.start.time?.split("T")[1].substring(0,8),
            end: interval.interval.end.time?.split("T")[1].substring(0,8),
            color: `${colorHash(interval.ctxId)}`,
            ctxId: interval.ctxId,
            description: interval.description
        }))}) as Record<string, TimeInterval[]>});



    return (
        <div className="flex flex-col">
            <div className="flex-1 flex items-center justify-center">
            </div>
            <Timeline data={intervals ?? {}} hideDates={true} onItemSelect={(interval) => {setSelectedInterval(interval)}}/>
            <SectionCards contextList={summary?.contexts} term={selectedInterval?.description ?? ''} expandId={selectedInterval?.ctxId ?? ''}></SectionCards>
        </div>
    );
}

export default TodaySummary;
