import {http} from "@/api/api";

export interface Context {
    id: string,
    description: string
}

export class ContextApi {
    contextList = () => http.get<Context[]>("/context/list").then(response => response.data);
    contextListQuery = {queryKey: ["contextList"], queryFn: this.contextList};

}

