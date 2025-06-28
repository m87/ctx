import {http, Interval, mapInterval} from "@/api/api";


export interface IntervalEntry {
    id: string,
    ctxId: string,
    description: string,
    interval: Interval
}

export interface IntervalsResponse {
    days: IntervalsResponseEntry[]
}

export interface IntervalsResponseEntry{
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


export class IntervalsApi {
    intervalsByDay = (day: string) => http.get<IntervalsResponse>(`/intervals/${day}`).then(response => {
        response.data.days = response.data.days.map(mapIntervalResponseEntry)
        return response.data;
    });
    intervalsByDayQuery = {queryKey: ["intervals"], queryFn: this.intervalsByDay};

    intervals = () => http.get<IntervalsResponse>(`/intervals`).then(response => {
        response.data.days = response.data.days.map(mapIntervalResponseEntry)
        return response.data;
    });
    intervalsQuery = {queryKey: ["intervals"], queryFn: this.intervals};

    recentIntervals = () => http.get<IntervalsResponse>(`/intervals/recent/10`).then(response => {
        response.data.days = response.data.days.map(mapIntervalResponseEntry)
        return response.data;
    });
    recentIntervalsQuery = {queryKey: ["recentIntervals-10"], queryFn: this.recentIntervals};

    move = (data: {
        src: string,
        target: string,
        id: string
    }) => http.post<void>("/intervals/move", data).then(response => response)

    delete = (ctxId: string, id: string) => http.delete<void>(`/intervals/${ctxId}/${id}`).then(response => response)



}
