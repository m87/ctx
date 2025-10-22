import { api } from "@/api/api";
import { Separator } from "@/components/ui/separator"
import { SidebarTrigger } from "@/components/ui/sidebar"
import { durationAsH, durationAsHM, durationAsM } from "@/lib/utils";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { format, isValid, parseISO } from "date-fns";
import { Clock, Clock10, Pause } from "lucide-react";
import { useEffect, useState } from "react";
import { Route, Routes, useParams } from "react-router-dom";
import { Card } from "./ui/card";
import { Button } from "./ui/button";
import { Spinner } from "./ui/spinner";
import { TickingClock } from "./ui/tickclock";
import { StopIcon } from "@radix-ui/react-icons";

export function SiteHeader() {

  const [selectedDate, setSelectedDate] = useState<Date>(new Date());
  const { data: summary } = useQuery({ ...api.summary.daySummaryQuery(format(selectedDate, "yyyy-MM-dd"), false) });
  const { data: currentContext } = useQuery({ ...api.context.currentQuery, refetchInterval: 5000 });
  const querClient = useQueryClient()
  const freeMutation = useMutation(api.context.freeMutaiton(querClient))
  const { day } = useParams();

  const [hovered, setHovered] = useState(false);

  useEffect(() => {
    if (day) {
      const date = parseISO(day)
      if (isValid(date)) {
        setSelectedDate(date)
      }
    }
  }, [day])




  return (
    <header
      className="group-has-data-[collapsible=icon]/sidebar-wrapper:h-12 flex h-12 shrink-0 items-center gap-2 border-b transition-[width,height] ease-linear">
      <div className="flex w-full items-center gap-1 px-4 lg:gap-2 pr-2 mr-0">
        <SidebarTrigger className="-ml-1" />
        <Separator
          orientation="vertical"
          className="mx-2 data-[orientation=vertical]:h-4"
        />
        <h1 className="text-base font-medium flex w-full justify-start flex-norwap">
          <Routes>
            <Route path="/contexts" element={"Contexts"} />
            <Route path="/day/:day" element={
              (() => {
                const DayComponent = () => {
                  const { day } = useParams();
                  return <div>{day}</div>;
                };
                return <DayComponent />;
              })()

            } />
            <Route path="/today" element={new Date().toLocaleDateString()} />
            <Route path="/" element={new Date().toLocaleDateString()} />
          </Routes>
          {summary?.duration ? <div className="ml-5 flex gap-1 items-center text-muted-foreground text-sm"><Clock size={16}></Clock>{durationAsHM(summary?.duration)} </div> : <div></div>}
        </h1>
        {currentContext?.context.description &&
          <div className="flex rounded-lg font-li items-center justify-between"
            onMouseEnter={() => setHovered(true)}
            onMouseLeave={() => setHovered(false)}
          >
            <div className="flex items-center gap-2 border rounded-md pl-2">
              <div className="flex items-center gap-2"><span className="text-ellipsis overflow-hidden whitespace-nowrap max-w-64 mr-1">{currentContext?.context.description}</span></div>
              <div className="flex gap-1 items-center text-muted-foreground text-sm mt-2 mb-2"><TickingClock size={16}></TickingClock>{durationAsHM(currentContext?.currentDuration)} </div>
              <Button variant="ghost" onClick={() => freeMutation.mutate({ day: day })} className="flex items-center gap-2 border-l-2 rounded-l-none shadow-none p-1 pr-2 pl-2"><Pause size={16} className="cursor-pointer shrink-0" ></Pause> </Button>
            </div>
          </div>
        }
      </div>
    </header>
  )
}
