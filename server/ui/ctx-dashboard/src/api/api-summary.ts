import { Context, http, mapContext } from "@/api/api";

export interface DaySummary {
    contexts: Context[];
    otherContexts: Context[];
    duration: number;
}

export class SummaryApi {
    daySummary = (day: string, showAllContexts: boolean) => http.get<DaySummary>(`/summary/day/${day}`, { params: { "showAllContexts": showAllContexts } }).then(response => {response.data.contexts = (response.data.contexts ?? []).map(mapContext); response.data.otherContexts = (response.data.otherContexts ?? []).map(mapContext); return response.data; });
    daySummaryQuery = (day: string, showAllContexts: boolean) => ({ queryKey: ["daySummary", day, showAllContexts], queryFn: () => this.daySummary(day, showAllContexts) });

    todaySummary = () => http.get<DaySummary>(`/summary/day`).then(response => { response.data.contexts = response.data.contexts.map(mapContext); return response.data; });
    todaySummaryQuery = { queryKey: ["todaySummary"], queryFn: this.todaySummary };

    dayListSummary = () => http.get<{ [key: string]: number }>("/summary/day/list").then(response => response.data);
    dayListSummaryQuery = { queryKey: ["dayListSummary"], queryFn: this.dayListSummary };
}



