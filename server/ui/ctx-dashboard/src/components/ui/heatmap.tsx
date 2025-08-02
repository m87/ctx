import { useMemo } from "react";
import {
  format,
  eachDayOfInterval,
  startOfWeek,
  endOfWeek,
  addDays,
  endOfToday,
} from "date-fns";

// Opcjonalnie możesz dodać tooltip np. z radix-ui

const COLOR_SCALE = ["#ebedf0", "#c6e48b", "#7bc96f", "#239a3b", "#196127"];

function getColor(count: number): string {
  if (count === 0) return COLOR_SCALE[0];
  if (count >= 8) return COLOR_SCALE[4];
  if (count >= 5) return COLOR_SCALE[3];
  if (count >= 3) return COLOR_SCALE[2];
  return COLOR_SCALE[1];
}

type ContextData = {
  date: string; // "2025-07-30"
  count: number;
};

type Props = {
  data?: ContextData[];
  weeksToShow?: number;
};

const exampleData: ContextData[] = [
  { date: "2025-06-15", count: 1 },
  { date: "2025-06-16", count: 4 },
  { date: "2025-06-20", count: 7 },
  { date: "2025-06-25", count: 2 },
  { date: "2025-07-01", count: 3 },
  { date: "2025-07-10", count: 5 },
  { date: "2025-07-20", count: 8 },
  { date: "2025-07-30", count: 6 },
];

export default function ContextHeatmap({ data = exampleData, weeksToShow = 53 }: Props) {
  const dataMap = useMemo(() => {
    const map: Record<string, number> = {};
    data?.forEach((d) => {
      if (d?.date) {
        map[d.date] = d.count;
      }
    });
    return map;
  }, [data]);

  const endDate = endOfToday();
  const startDate = addDays(endDate, -weeksToShow * 7);

  const allDays = eachDayOfInterval({
    start: startOfWeek(startDate, { weekStartsOn: 0 }),
    end: endOfWeek(endDate, { weekStartsOn: 0 }),
  });

  const weeks = [];
  for (let i = 0; i < allDays.length; i += 7) {
    weeks.push(allDays.slice(i, i + 7));
  }

  return (
    <div className="flex gap-1 overflow-x-auto text-xs">
      {weeks.map((week, i) => (
        <div key={i} className="flex flex-col gap-1">
          {week.map((day) => {
            const dateStr = format(day, "yyyy-MM-dd");
            const count = dataMap[dateStr] || 0;
            return (
              <div
                key={dateStr}
                className="w-3 h-3 rounded-sm cursor-pointer"
                style={{ backgroundColor: getColor(count) }}
                title={`${count} context${count !== 1 ? "s" : ""} on ${format(
                  day,
                  "MMM d, yyyy"
                )}`}
              />
            );
          })}
        </div>
      ))}
    </div>
  );
}
