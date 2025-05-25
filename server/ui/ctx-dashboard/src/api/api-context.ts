import { Context, http, mapContext, ZonedDateTime } from "@/api/api";


export class ContextApi {
    list = () => http.get<Context[]>("/context/list").then(response => response.data.map(mapContext));
    listQuery = { queryKey: ["contextList"], queryFn: this.list, select: (data: Context[]) => data.sort((a, b) => a.id.localeCompare(b.id)) };

    current = () => http.get<Context>("/context/current").then(response => response.data ? mapContext(response.data) : null);
    currentQuery = { queryKey: ["currentContext"], queryFn: this.current };

    free = () => http.post<void>("/context/free").then(response => response)

    switch = (id: string) => http.post<void>("/context/switch", { id: id }).then(response => response)

    createAndSwitch = (description: string) =>
        http.post<Context>("/context/createAndSwitch", { description: description }).then(response => response.data);

    updateInterval = (contextId: string, intervalId: string, start: ZonedDateTime, end: ZonedDateTime) =>
        http.put<void>("/context/interval", { contextId: contextId, intervalId: intervalId, start: start, end: end }).then(response => response);
}
