import { api, Interval, mapZoned, ZonedDateTime } from "@/api/api"
import { TimeInterval } from "./timeline"
import { useRef, useState } from "react"
import { Input } from "./ui/input"
import { Button } from "./ui/button"



export interface TimelineSplitData {
  interval: TimeInterval
}

export function TimelineSplit({ interval }: TimelineSplitData) {

  const ref = useRef<HTMLInputElement>()

  return (
    <div className="flex">
      <Input className="w-24"  ref={ref}></Input><Button onClick={() => api.intervals.split(interval.ctxId, interval.id, mapZoned({time: ref.current?.value}))}></Button>
    </div>
  )

}

export default TimelineSplit
