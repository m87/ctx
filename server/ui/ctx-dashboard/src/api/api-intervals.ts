import {http, Interval, mapInterval} from "@/api/api";


export interface IntervalEntry {
    id: string,
    ctxId: string,
    description: string,
    interval: Interval
}

export interface IntervalsResponse {
    intervals: IntervalEntry[]
}


export function mapIntervalEntry(obj: any): IntervalEntry {
    return {
        id: obj.id,
        ctxId: obj.ctxId,
        description: obj.description,
        interval: mapInterval(obj.interval)
    };
}


export class IntervalsApi {
    intervalsByDay = (day: string) => http.get<IntervalsResponse>(`/intervals/${day}`).then(response => {
        response.data.intervals = response.data.intervals.map(mapIntervalEntry);
        return response.data;
    });
    intervalsByDayQuery = {queryKey: ["intervals"], queryFn: this.intervalsByDay};

    intervals = () => http.get<IntervalsResponse>(`/intervals`).then(response => {
        response.data.intervals = response.data.intervals.map(mapIntervalEntry);
        return response.data;
    });
    intervalsQuery = {queryKey: ["intervals"], queryFn: this.intervals};
}
