import { http, Interval, mapInterval, ZonedDateTime } from "@/api/api";
import { QueryClient } from "@tanstack/react-query";


export interface IntervalEntry {
  id: string,
  ctxId: string,
  description: string,
  interval: Interval
}

export interface IntervalsResponse {
  days: IntervalsResponseEntry[]
}

export interface IntervalsResponseEntry {
  date: string,
  intervals?: IntervalEntry[]
}

export function mapIntervalEntry(obj: any): IntervalEntry {
  return {
    id: obj.id,
    ctxId: obj.ctxId,
    description: obj.description,
    interval: mapInterval(obj.interval)
  };
}

export function mapIntervalResponseEntry(obj: any): IntervalsResponseEntry {
  return {
    date: obj.date,
    intervals: obj.intervals?.map(mapIntervalEntry)
  }
}

export interface Split {
  h: number,
  m: number,
  s: number
}

export function TimeStringAsSplit(time: string): Split {
  const [h, m, s] = time.split(":").map(Number)
  return { h, m, s }
}

export interface SplitPayload {
  ctxId: string,
  id: string,
  split: Split,
  day: string
}

export class IntervalsApi {
  intervalsByDay = (day: string) => http.get<IntervalsResponse>(`/intervals/${day}`).then(response => {
    response.data.days = response.data.days.map(mapIntervalResponseEntry)
    return response.data;
  });
  intervalsByDayQuery = (day: string) => ({ queryKey: ["intervals", day], queryFn: () => this.intervalsByDay(day) });

  intervals = () => http.get<IntervalsResponse>(`/intervals`).then(response => {
    response.data.days = response.data.days.map(mapIntervalResponseEntry)
    return response.data;
  });
  intervalsQuery = { queryKey: ["intervals"], queryFn: this.intervals };

  recentIntervals = () => http.get<IntervalsResponse>(`/intervals/recent/10`).then(response => {
    response.data.days = response.data.days.map(mapIntervalResponseEntry)
    return response.data;
  });
  recentIntervalsQuery = { queryKey: ["recentIntervals-10"], queryFn: this.recentIntervals };

  move = (data: {
    src: string,
    target: string,
    id: string
  }) => http.post<void>("/intervals/move", data).then(response => response)

  delete = (ctxId: string, id: string) => http.delete<void>(`/intervals/${ctxId}/${id}`).then(response => response)

  split = (ctxId: string, id: string, split: Split) => http.post<void>(`/intervals/${ctxId}/${id}/split`, { split: split }).then(response => response)
  splitMutation = (queryClient: QueryClient) => ({
    mutationFn: (payload: SplitPayload) => this.split(payload.ctxId, payload.id, payload.split),
    onSuccess: (_, variables: SplitPayload) => {
      queryClient.invalidateQueries({ queryKey: ["intervals", variables.day] })
    },
  })

}
