import { inject, Injectable } from '@angular/core';
import { IntervalService } from './interval.service';
import { lastValueFrom } from 'rxjs';

@Injectable({ providedIn: 'root' })
export class IntervalQueries {
  private intervalService = inject(IntervalService);
  readonly key = ['interval'];

  get(id: string) {
    return {
      queryKey: [...this.key, id],
      queryFn: () => lastValueFrom(this.intervalService.getInterval(id)),
    };
  }

  day(day: string) {
    return {
      queryKey: [...this.key, 'day', day],
      queryFn: () => lastValueFrom(this.intervalService.getDayIntervals(day)),
    };
  }
}
