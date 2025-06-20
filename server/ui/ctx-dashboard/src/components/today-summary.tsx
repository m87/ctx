import { useQuery } from "@tanstack/react-query";
import {api} from "@/api/api";
import { SectionCards } from "./section-cards";
import Timeline, {TimeInterval} from "@/components/timeline";
import {IntervalsResponse} from "@/api/api-intervals";
import {colorHash} from "@/lib/utils"

export function TodaySummary() {
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
            <Timeline data={intervals ?? {}} hideDates={true}/>
            <SectionCards contextList={summary?.contexts}></SectionCards>
        </div>
    );
}

export default TodaySummary;
