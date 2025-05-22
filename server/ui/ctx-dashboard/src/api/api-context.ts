import {http} from "@/api/api";

export interface Context {
    id: string,
    description: string
}

export class ContextApi {
    contextList = () => http.get<Context[]>("/context/list").then(response => response.data);
    contextListQuery = {queryKey: ["contextList"], queryFn: this.contextList};

    currentContext = () => http.get<Context>("/context/current").then(response => response.data);
    currentContextQuery = {queryKey: ["currentContext"], queryFn: this.currentContext};

}