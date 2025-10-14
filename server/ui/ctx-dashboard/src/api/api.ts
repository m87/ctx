import axios, { AxiosRequestTransformer } from "axios";
import {ContextApi} from "@/api/api-context";
import { DateTime, Zone } from "luxon";
import { SummaryApi } from "./api-summary";
import {IntervalsApi} from "@/api/api-intervals";
import { QueryClient } from "@tanstack/react-query";


export class ZonedDateTime {
  constructor(public time: string | null, public timezone: string | null) {}

  public static fromDateTime(dt: DateTime): ZonedDateTime {
    return new ZonedDateTime(dt.toISO(), dt.zoneName);
  }

  public static fromtTimeString(time: string, timezone: string): ZonedDateTime {
    return new ZonedDateTime(DateTime.fromISO(time).toISO(), timezone)
  }

  public toDateTime(): DateTime {
    return DateTime.fromISO(this.time ?? '', { zone: this.timezone ?? "utc" });
  }

  public toInputValue(): string {
    return this.toDateTime().toFormat("yyyy-MM-dd'T'HH:mm");
  }

  public toString(): string {
    return this.toDateTime().toFormat("yyyy-MM-dd HH:mm");
  }
}

export interface Interval {
    id: string
    start: ZonedDateTime,
    end: ZonedDateTime,
    duration: number,
}

export interface Context {
    id: string,
    description: string,
    intervals: { [key: string]: Interval },
    duration: number,
    labels: string[],
}

export function mapZoned(obj: any): ZonedDateTime {
  return new ZonedDateTime(obj.time ?? null, obj.timezone ?? null);
}

export function mapInterval(obj: any): Interval {
  return {
    id: obj.id,
    duration: obj.duration,
    start: mapZoned(obj.start),
    end: mapZoned(obj.end),
  };
}

export function mapContext(obj: any): Context {
  return {
    id: obj.id,
    description: obj.description,
    duration: obj.duration,
    intervals: Object.fromEntries(Object.entries(obj.intervals).map(([key, value]) => [key, mapInterval(value)])),
    labels: obj.labels ?? [],
  };
}

export const httpConfig = {
    baseURL: "/api",
    withCredentials: true,
    withXSRFToken: true,
    timeout: 6000,
    headers: {
        Accept: "application/json",
    },
};


export function invalidateQueriesByDate(queryClient: QueryClient, variables: any) {
      const today = DateTime.now().toFormat("yyyy-MM-dd")
      queryClient.invalidateQueries({ queryKey: ["intervals", variables.day ?? today] })
      queryClient.invalidateQueries({ queryKey: ["currentContext"] })
      queryClient.invalidateQueries({ queryKey: ["todaySummary"] })
      queryClient.invalidateQueries({ queryKey: ["daySummary", variables.day ?? today, true] })
      queryClient.invalidateQueries({ queryKey: ["daySummary", variables.day ?? today, false] })
}
 

export const http = axios.create(httpConfig);

export class Api {
    context = new ContextApi();
    summary = new SummaryApi();
    intervals = new IntervalsApi();

    version = () => http.get<string>("/version").then(response => response.data["version"]);
    versionQuery = { queryKey: ["version"], queryFn: this.version};
}

export const api = new Api();
