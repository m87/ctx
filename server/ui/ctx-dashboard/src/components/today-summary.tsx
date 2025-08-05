import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { api, mapZoned, ZonedDateTime } from "@/api/api";
import { SectionCards } from "./section-cards";
import Timeline, { intervalsResponseAsTimelineData, TimeInterval } from "@/components/timeline";
import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { format, isValid, parseISO } from "date-fns";
import { TimeStringAsSplit } from "@/api/api-intervals";
import { DateTime } from "luxon";


export function TodaySummary() {
  const [selectedDate, setSelectedDate] = useState<Date>(new Date());
  const [selectedInterval, setSelectedInterval] = useState(null)
  const { data: summary } = useQuery({ ...api.summary.daySummaryQuery(format(selectedDate, "yyyy-MM-dd")) });
  const { data: intervals } = useQuery({ ...api.intervals.intervalsByDayQuery(format(selectedDate, "yyyy-MM-dd")), select: intervalsResponseAsTimelineData })
  const { data: names } = useQuery({ ...api.context.listNamesQuery })
  const { day } = useParams();

  const querClient = useQueryClient();
  const splitMutation = useMutation(api.intervals.splitMutation(querClient))
  const moveMutation = useMutation(api.intervals.moveMutation(querClient))
  const deleteMutation = useMutation(api.intervals.deleteMutation(querClient))


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
      <Timeline
        data={intervals ?? {}}
        ctxNames={names ?? []}
        hideDates={true}
        hideGuides={true}
        onItemSelect={(interval) => { setSelectedInterval(interval) }}
        onItemSplit={(interval, time) => splitMutation.mutate({ ctxId: interval.ctxId, id: interval.id, split: TimeStringAsSplit(time), day: day ?? DateTime.now().toFormat("yyyy-MM-dd") })}
        onItemMove={(interval, ctx) => moveMutation.mutate({
          src: interval.ctxId,
          target: ctx.id,
          id: interval.id,
          day: day
        })}
        onItemDelete={(interval: TimeInterval) =>  deleteMutation.mutate({ctxId: interval.ctxId, id: interval.id, day: day})}
      />
      <SectionCards contextList={summary?.contexts} term={selectedInterval?.description ?? ''} expandId={selectedInterval?.ctxId ?? ''}></SectionCards>
    </div>
  );
}

export default TodaySummary;
