import { useQuery } from "@tanstack/react-query";
import { api } from "@/api/api";
import { SectionCards } from "./section-cards";
import Timeline, { intervalsResponseAsTimelineData } from "@/components/timeline";
import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { format, isValid, parseISO } from "date-fns";


export function TodaySummary() {
  const [selectedDate, setSelectedDate] = useState<Date>(new Date());
  const [selectedInterval, setSelectedInterval] = useState(null)
  const { data: summary } = useQuery({ ...api.summary.daySummaryQuery(format(selectedDate, "yyyy-MM-dd")) });
  const { data: intervals } = useQuery({ ...api.intervals.intervalsByDayQuery(format(selectedDate, "yyyy-MM-dd")), select: intervalsResponseAsTimelineData })
  const { data: names } = useQuery({ ...api.context.listNamesQuery })
  const { day } = useParams();

  useEffect(() => {
    if (day) {
      const date = parseISO(day)
      if (isValid(date)) {
        setSelectedDate(date)
      }
    }
  }, [day])



  return (
    <div className="flex flex-col">
      <div className="flex-1 flex items-center justify-center">
      </div>
      <Timeline data={intervals ?? {}} ctxNames={names ?? []} hideDates={true} hideGuides={true} onItemSelect={(interval) => { setSelectedInterval(interval) }} />
      <SectionCards contextList={summary?.contexts} term={selectedInterval?.description ?? ''} expandId={selectedInterval?.ctxId ?? ''}></SectionCards>
    </div>
  );
}

export default TodaySummary;
