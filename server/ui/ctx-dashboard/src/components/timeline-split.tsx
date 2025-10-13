import { TimeInterval } from "./timeline"
import { useEffect, useState } from "react"
import { Slider } from "./ui/slider"
import { Button } from "./ui/button"



export interface TimelineSplitData {
  interval: TimeInterval
  onChange: (interval: TimeInterval, n: string) => void
}

function timeToDecimal(time: string): number {
  const [h, m, s] = time.split(":").map(Number);
  return (h * 3600 + m * 60 + s) / 60;
}

function decimalToTime(time: number): string {
  const h = Math.floor(time / 60);
  const m = time % 60;

  const hh = h.toString().padStart(2, '0');
  const mm = m.toString().padStart(2, '0');
  return `${hh}:${mm}:00`
}

export function TimelineSplit({ interval, onChange }: TimelineSplitData) {

  const [start, setStart] = useState(0)
  const [end, setEnd] = useState(0)
  const [value, setValue] = useState(0)

  useEffect(() => {
    const s = timeToDecimal(interval.start)
    const e = timeToDecimal(interval.end)
    setStart(s)
    setEnd(e)
    setValue(Math.floor(s + (e - s) / 2))

  }, [interval.start, interval.end])


  return (
    <div className="flex flex-col items-center gap-2">
      <div className="gap-2">{decimalToTime(value)}</div>
      <Slider className="gap-2" onValueChange={(n) => setValue(n[0])} disabled={end - start <= 1} value={[value]} min={start} max={end} step={1} />
      <Button className="gap-2 w-full" onClick={(e) => {
        onChange(interval, decimalToTime(value));
        e.currentTarget.dispatchEvent(
          new KeyboardEvent("keydown", { key: "Escape", bubbles: true })
        );
      }}>Split</Button>
    </div>
  )

}

export default TimelineSplit
