import {http} from "@/api/api";

export interface Context {
    id: string,
    description: string
}

export class ContextApi {
    list = () => http.get<Context[]>("/context/list").then(response => response.data);
    listQuery = {queryKey: ["contextList"], queryFn: this.list, select: (data: Context[]) => data.sort((a,b) => a.id.localeCompare(b.id))};

    current = () => http.get<Context>("/context/current").then(response => response.data);
    currentQuery = {queryKey: ["currentContext"], queryFn: this.current};

    free = () => http.post<void>("/context/free").then(response => response)

    switch = (id: string) => http.post<void>("/context/switch", {id: id}).then(response => response)
}
