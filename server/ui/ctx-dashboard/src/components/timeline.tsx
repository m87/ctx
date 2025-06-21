import {clsx} from "clsx";
import { useEffect, useRef, useState } from "react";

export interface TimeInterval {
  start: string; 
  end: string;   
  color?: string;
  ctxId: string;
  description?: string; 
}

export interface TimelineProps {
  data: Record<string, TimeInterval[]>; 
  hideDates: boolean;
  onItemSelect: (interval: TimeInterval | null) => void;
}

function timeToDecimal(time: string): number {
  const [h, m, s] = time.split(":").map(Number);
  return (h * 3600 + m * 60 + s) / 86400;
}

function Timeline({ data, hideDates, onItemSelect }: TimelineProps) {
  const hours = Array.from({ length: 24 }, (_, i) => i);
  const dates = Object.keys(data);
  const boxRefs = useRef<Map<string, HTMLDivElement>>(new Map());
  const BLOCK = 100;

  useEffect(() => {
    function handleClick(event: MouseEvent) {
      const target = event.target as Node;
      const clickedOnBox = Array.from(boxRefs.current.values()).some((boxEl) => boxEl.contains(target));

      if(!clickedOnBox) {
 //       setSelected('')
 //       onItemSelect(null)
      }
    }

    function handleKeyDown(event: KeyboardEvent) {
      if(event.key === 'Escape') {
        setSelected('')
        onItemSelect(null)
      }
    }

    document.addEventListener('mousedown', handleClick);
    document.addEventListener('keydown', handleKeyDown);

    return () => {
      document.removeEventListener('mousedown', handleClick);
      document.removeEventListener('keydown', handleKeyDown);
    }

  }, []);

  const [selected, setSelected] = useState('')
  return (
    <div className="font-sans bg-gray-100">
      <div className="relative p-4 overflow-x-auto">
        <div className="min-w-[1000px]" >
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
                      onClick={() => {
                        if (selected === interval.ctxId) {
                          setSelected('') 
                            onItemSelect(null);
                        }
                          else {
                            setSelected(interval.ctxId); 
                            onItemSelect(interval);
                          }}}
                      ref={(el) => el && boxRefs.current.set(interval.ctxId, el) }
                      key={`interval-${rowIndex}-${idx}`}
                      className={`${selected && interval.ctxId !== selected ? 'opacity-50' : ''} cursor-pointer absolute top-0 h-full rounded text-xs text-white flex items-center justify-center px-1 text-ellipsis overflow-hidden whitespace-nowrap`}
                      style={{ left, width, backgroundColor: interval.color }}
                    >
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
