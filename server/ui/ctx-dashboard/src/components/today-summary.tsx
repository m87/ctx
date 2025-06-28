import { useQuery } from "@tanstack/react-query";
import {api} from "@/api/api";
import { SectionCards } from "./section-cards";
import Timeline, {intervalsResponseAsTimelineData} from "@/components/timeline";
import { useState } from "react";


export function TodaySummary() {
    const [selectedInterval, setSelectedInterval] = useState(null)
    const {data: summary} = useQuery({...api.summary.todaySummaryQuery});
    const {data: intervals} = useQuery({...api.intervals.intervalsQuery, select: intervalsResponseAsTimelineData})
    const {data: names} = useQuery({...api.context.listNamesQuery})

    return (
        <div className="flex flex-col">
            <div className="flex-1 flex items-center justify-center">
            </div>
            <Timeline data={intervals ?? {}} ctxNames={names} hideDates={true} hideGuides={true} onItemSelect={(interval) => {setSelectedInterval(interval)}}/>
            <SectionCards contextList={summary?.contexts} term={selectedInterval?.description ?? ''} expandId={selectedInterval?.ctxId ?? ''}></SectionCards>
        </div>
    );
}

export default TodaySummary;
