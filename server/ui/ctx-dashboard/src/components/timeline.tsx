import {clsx} from "clsx";

export interface TimeInterval {
  start: string; // format: "HH:MM:SS"
  end: string;   // format: "HH:MM:SS"
  color?: string; // Optional color for the block
}

export interface TimelineProps {
  data: Record<string, TimeInterval[]>; // Keyed by date
  hideDates: boolean;
}

function timeToDecimal(time: string): number {
  const [h, m, s] = time.split(":").map(Number);
  return (h * 3600 + m * 60 + s) / 86400;
}

function Timeline({ data, hideDates }: TimelineProps) {
  const hours = Array.from({ length: 24 }, (_, i) => i);
  const dates = Object.keys(data);
  const BLOCK = 100;
  return (
    <div className="font-sans bg-gray-100">
      <div className="relative p-4 overflow-x-auto">
        <div className="min-w-[1000px]">
          <div className={clsx("relative h-6 mb-2", hideDates ? "" : "ml-24")}>
            {hours.map((hour) => (
              <div
                key={`hour-${hour}`}
                className="absolute top-0 text-xs font-semibold text-center border-l border-black/30"
                style={{
                  left: `${(hour / 24) * BLOCK}%`,
                  width: `${BLOCK / 24}%`,
                }}
              >
                {hour}
              </div>
            ))}
          </div>

          {dates.map((date, rowIndex) => (
            <div key={`row-${rowIndex}`} className="relative w-full h-10 mb-1 flex items-center">
              {!hideDates &&
              <div className="w-24 pr-2 text-sm font-medium text-right text-gray-700">
                {date}
              </div>
              }
              <div className="relative flex-1 h-full">
                {data[date].map((interval, idx) => {
                  const start = timeToDecimal(interval.start);
                  const end = timeToDecimal(interval.end);
                  const left = `${start * BLOCK}%`;
                  const width = `${(end - start) * BLOCK}%`;

                  return (
                    <div
                      key={`interval-${rowIndex}-${idx}`}
                      className={`absolute top-0 h-full rounded text-xs text-white flex items-center justify-center px-1 ${interval.color || "bg-blue-500"}`}
                      style={{ left, width }}
                    >
                      {interval.start} - {interval.end}
                    </div>
                  );
                })}
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}

export default Timeline;
