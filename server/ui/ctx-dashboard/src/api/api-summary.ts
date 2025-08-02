import { Context, http, mapContext } from "@/api/api";

export interface DaySummary {
    contexts : Context[];
    duration: number;
}

export class SummaryApi {
    daySummary = (day: string) => http.get<DaySummary>(`/summary/day/${day}`).then(response => {response.data.contexts = response.data.contexts.map(mapContext); return response.data; });
    daySummaryQuery = { queryKey: ["todaySummary"], queryFn: this.daySummary};

    todaySummary = () => http.get<DaySummary>(`/summary/day`).then(response => {response.data.contexts = response.data.contexts.map(mapContext); return response.data; });
    todaySummaryQuery = { queryKey: ["todaySummary"], queryFn: this.todaySummary };

    dayListSummary = () => http.get<{[key: string]: number}>("/summary/day/list").then(response => response.data);
    dayListSummaryQuery = { queryKey: ["dayListSummary"], queryFn: this.dayListSummary };
}



