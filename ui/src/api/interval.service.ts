import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { map, Observable } from 'rxjs';
import { DateTime } from 'luxon';
import type { Context } from './context.service';
import { Store } from '@ngxs/store';
import { WorkspaceState } from '../app/sidebar/workspace.state';

export interface Interval {
  id: string;
  contextId: string;
  start: ZonedDateTime;
  end: ZonedDateTime;
  duration: number;
  workspaceId: string;
}

type RawZonedDateTime = {
  time: string | null;
  timezone: string | null;
  isZero: boolean | null;
};

export type RawInterval = Omit<Interval, 'start' | 'end'> & {
  start: RawZonedDateTime;
  end: RawZonedDateTime;
};

export interface DayIntervalsResponse {
  contexts: Context[];
  intervals: Interval[];
}

type RawDayIntervalsResponse = {
  contexts: Context[];
  intervals: RawInterval[];
};

export class ZonedDateTime {
  constructor(
    public time: string | null,
    public timezone: string | null,
    public isZero: boolean | null,
  ) {}

  public static fromDateTime(dt: DateTime): ZonedDateTime {
    return new ZonedDateTime(dt.toISO(), dt.zoneName, dt.year === 1);
  }

  public static fromTimeString(time: string, timezone: string): ZonedDateTime {
    return new ZonedDateTime(
      DateTime.fromISO(time).toISO(),
      timezone,
      DateTime.fromISO(time).year === 1,
    );
  }

  public toDateTime(): DateTime {
    const time = this.time ?? '';
    const parsed = DateTime.fromISO(time);
    if (!parsed.isValid) {
      return parsed;
    }

    const zone = this.normalizedZone();
    if (!zone) {
      return parsed;
    }

    const zoned = parsed.setZone(zone);
    return zoned.isValid ? zoned : parsed;
  }

  public toInputValue(): string {
    return this.toDateTime().toFormat("yyyy-MM-dd'T'HH:mm");
  }

  public toString(): string {
    if (this.isZero) {
      return '...';
    }
    return this.toDateTime().toFormat('yyyy-MM-dd HH:mm');
  }

  public toTimeString(): string {
    if (this.isZero) {
      return '...';
    }
    return this.toDateTime().toFormat('HH:mm');
  }

  public toDateString(): string {
    if (this.isZero) {
      return '...';
    }
    return this.toDateTime().toFormat('yyyy-MM-dd');
  }

  private normalizedZone(): string | null {
    if (!this.timezone || this.timezone === 'Local') {
      return DateTime.local().zoneName;
    }
    return this.timezone;
  }
}

export function deserializeInterval(interval: RawInterval): Interval {
  return {
    ...interval,
    start: new ZonedDateTime(
      interval.start?.time ?? null,
      interval.start?.timezone ?? null,
      interval.start?.isZero ?? null,
    ),
    end: new ZonedDateTime(
      interval.end?.time ?? null,
      interval.end?.timezone ?? null,
      interval.end?.isZero ?? null,
    ),
  };
}

export function deserializeIntervals(intervals: RawInterval[]): Interval[] {
  return intervals.map((interval) => deserializeInterval(interval));
}

export function parseDuration(duration: number): string {
  const hours = Math.floor(duration / 3600);
  const minutes = Math.floor((duration % 3600) / 60);
  return `${hours}h ${minutes}m`;
}

@Injectable({ providedIn: 'root' })
export class IntervalService {
  http = inject(HttpClient);
  store = inject(Store);

  createInterval(interval: Interval): Observable<Interval> {
    if (interval.workspaceId == null) {
      interval.workspaceId = this.store.selectSnapshot(WorkspaceState.selectedWorkspaceId)!;
    }
    return this.http
      .post<RawInterval>('/api/interval/', interval)
      .pipe(map((response) => deserializeInterval(response)));
  }

  deleteInterval(id: string): Observable<void> {
    return this.http.delete<void>(`/api/interval/${id}`);
  }

  updateInterval(id: string, interval: Interval): Observable<Interval> {
    return this.http
      .put<RawInterval>(`/api/interval/${id}`, interval)
      .pipe(map((response) => deserializeInterval(response)));
  }

  moveInterval(id: string, targetContextId: string): Observable<Interval> {
    return this.http
      .patch<RawInterval>(`/api/interval/${id}/move/${targetContextId}`, {})
      .pipe(map((response) => deserializeInterval(response)));
  }

  getInterval(id: string): Observable<Interval> {
    return this.http
      .get<RawInterval>(`/api/interval/${id}`)
      .pipe(map((response) => deserializeInterval(response)));
  }

  getDayIntervals(workspaceId: string, day: string): Observable<DayIntervalsResponse> {
    return this.http
      .get<RawDayIntervalsResponse>(`/api/interval/day/${day}?workspaceId=${workspaceId}`)
      .pipe(
        map((response) => ({
          contexts: response.contexts,
          intervals: deserializeIntervals(response.intervals),
        })),
      );
  }
}
