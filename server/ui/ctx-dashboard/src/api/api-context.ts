import { Context, http, invalidateQueriesByDate, mapContext, ZonedDateTime } from "@/api/api";
import { QueryClient } from "@tanstack/react-query";

export interface CurrentContext {
  context: Context,
  currentDuration: number
}

export class ContextApi {
  list = () => http.get<Context[]>("/context/list").then(response => response.data.map(mapContext));
  listQuery = { queryKey: ["contextList"], queryFn: this.list, select: (data: Context[]) => data.sort((a, b) => a.id.localeCompare(b.id)) };
  listNamesQuery = { queryKey: ["contextListNames"], queryFn: this.list, select: (data: Context[]) => data.sort((a, b) => a.id.localeCompare(b.id)).map((ctx) => ({ id: ctx.id, description: ctx.description })) };

  current = () => http.get<CurrentContext>("/context/current").then(response => response.data ? {
    context: mapContext(response.data.context),
    currentDuration: response.data.currentDuration
  }: null);
  currentQuery = { queryKey: ["currentContext"], queryFn: this.current };

  free = () => http.post<void>("/context/free").then(response => response)
  freeMutaiton = (queryClient: QueryClient) => ({
    mutationFn: (data: {day?: string}) => this.free(),
    onSuccess: (_, variables) => invalidateQueriesByDate(queryClient, variables),
  })

  switch = (id: string) => http.post<void>("/context/switch", { id: id }).then(response => response)
  switchMutation = (queryClient: QueryClient) => ({
    mutationFn: (data: {id: string, day?: string}) => this.switch(data.id),
    onSuccess: (_, variables) => invalidateQueriesByDate(queryClient, variables),
  }) 

  createAndSwitch = (description: string) =>
    http.post<Context>("/context/createAndSwitch", { description: description }).then(response => response.data);
  createAndSwitchMutation= (queryClient: QueryClient) => ({
    mutationFn: (data: {description: string, day?: string}) => this.createAndSwitch(data.description),
    onSuccess: (_, variables) => invalidateQueriesByDate(queryClient, variables),
  })


  updateInterval = (contextId: string, intervalId: string, start: ZonedDateTime, end: ZonedDateTime) =>
    http.put<void>("/context/interval", { contextId: contextId, intervalId: intervalId, start: start, end: end }).then(response => response);
  updateIntervalMutation = (queryClient: QueryClient) => ({
    mutationFn: (data: {contextId:string, intervalId: string, start: ZonedDateTime, end: ZonedDateTime, day?: string}) => this.updateInterval(data.contextId, data.intervalId, data.start, data.end),
    onSuccess: (_, variables) => invalidateQueriesByDate(queryClient, variables),
  })
}
