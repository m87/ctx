package core

import (
	"sort"
	"time"
)

type TimeRange struct {
	Start time.Time
	End   time.Time
}

func ClipIntervalRangeToDay(interval *Interval, date time.Time, now time.Time) (TimeRange, bool) {
	if interval == nil {
		return TimeRange{}, false
	}

	location := date.Location()
	if location == nil {
		location = time.UTC
	}
	localDayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, location)
	dayStart := localDayStart.UTC()
	dayEnd := localDayStart.AddDate(0, 0, 1).UTC()

	if !timeIsSet(interval.Start) {
		return TimeRange{}, false
	}
	start := interval.Start.UTC()

	var end time.Time
	if !timeIsSet(interval.End) {
		if interval.Status != "active" {
			return TimeRange{}, false
		}
		end = now
	} else {
		end = interval.End.UTC()
	}

	if end.Before(dayStart) || !start.Before(dayEnd) {
		return TimeRange{}, false
	}

	if start.Before(dayStart) {
		start = dayStart
	}
	if end.After(dayEnd) {
		end = dayEnd
	}

	if !end.After(start) {
		return TimeRange{}, false
	}

	return TimeRange{Start: start, End: end}, true
}

func ClipIntervalDurationToDay(interval *Interval, date time.Time, now time.Time) time.Duration {
	rng, ok := ClipIntervalRangeToDay(interval, date, now)
	if !ok {
		return 0
	}
	return rng.End.Sub(rng.Start)
}

func SumMergedRangesDuration(ranges []TimeRange) time.Duration {
	if len(ranges) == 0 {
		return 0
	}

	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i].Start.Before(ranges[j].Start)
	})

	mergedStart := ranges[0].Start
	mergedEnd := ranges[0].End
	var total time.Duration

	for _, rng := range ranges[1:] {
		if !rng.Start.After(mergedEnd) {
			if rng.End.After(mergedEnd) {
				mergedEnd = rng.End
			}
			continue
		}

		total += mergedEnd.Sub(mergedStart)
		mergedStart = rng.Start
		mergedEnd = rng.End
	}

	total += mergedEnd.Sub(mergedStart)
	return total
}
