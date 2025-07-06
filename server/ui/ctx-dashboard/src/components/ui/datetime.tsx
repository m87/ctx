"use client"

import * as React from "react"
import { ChevronDownIcon } from "lucide-react"

import { Button } from "@/components/ui/button"
import { Calendar } from "@/components/ui/calendar"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover"
import { DateTime } from "luxon"
import clsx from "clsx"


export interface DateTimeInputProperties {
  datetime: DateTime
  label?: string
  editable?: boolean
  onChange?: (dt: DateTime) => void
}

export function DateTimeInput({datetime, label, editable, onChange}: DateTimeInputProperties) {
  const [open, setOpen] = React.useState(false)
  const [date, setDate] = React.useState<DateTime>(datetime)

  React.useEffect(() => {
    if(onChange) {
      onChange(date)
    }
  }, [date, onChange])


  
  const handleDateChange = (selected: Date | undefined) => {
    if (!selected) return;
    const dt = DateTime.fromJSDate(selected);
    setDate((prev) => prev.set({
      year: dt.year,
      month: dt.month,
      day: dt.day
    }));
    setOpen(false)
  };

  const handleTimeChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const [hour, minute] = e.target.value.split(':').map(Number);
    setDate((prev) => prev.set({ hour, minute }));
  };

  return (
    <>
    {label && <Label htmlFor="date-picker" className="px-1">
       {label}
        </Label>
    }

    <div className="flex gap-4">
            <div className="flex flex-col gap-3">
         <Popover open={open} onOpenChange={setOpen}>
          <PopoverTrigger asChild>
            <Button
              variant="outline"
              id="date-picker"
              className={clsx(editable ? "" : "pointer-events-none", "w-32 justify-between font-normal")}
            >
              {date ? date.toFormat("dd-MM-yyyy") : "Select date"}
              { editable && <ChevronDownIcon /> }
            </Button>
          </PopoverTrigger>
          <PopoverContent className="w-auto overflow-hidden p-0" align="start">
            <Calendar
              mode="single"
              selected={date?.toJSDate()}
              captionLayout="dropdown"
              onSelect={handleDateChange}
            />
          </PopoverContent>
        </Popover>
      </div>
      <div className="flex flex-col gap-3">
        <Input
          type="time"
          id="time-picker"
          step="60"
          onChange={handleTimeChange}
          defaultValue={date.toFormat("HH:mm")}
          className={clsx(editable ? "" : "pointer-events-none" ,"bg-background appearance-none [&::-webkit-calendar-picker-indicator]:hidden [&::-webkit-calendar-picker-indicator]:appearance-none")}
        />
      </div>
    </div>
    </>
  )
}


