import { api } from "@/api/api";
import { Separator } from "@/components/ui/separator"
import { SidebarTrigger } from "@/components/ui/sidebar"
import { durationAsH, durationAsM, durationAsS } from "@/lib/utils";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { format, isValid, parseISO } from "date-fns";
import { Pause } from "lucide-react";
import { useEffect, useState } from "react";
import { Route, Routes, useParams } from "react-router-dom";

export function SiteHeader() {

  const [selectedDate, setSelectedDate] = useState<Date>(new Date());
  const { data: summary } = useQuery({ ...api.summary.daySummaryQuery(format(selectedDate, "yyyy-MM-dd")) });
  const { data: currentContext } = useQuery({ ...api.context.currentQuery, refetchInterval: 5000 });
  const querClient = useQueryClient()
  const freeMutation = useMutation(api.context.freeMutaiton(querClient))
  const { day } = useParams();

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
      <div className="flex w-full items-center gap-1 px-4 lg:gap-2 lg:px-6">
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
                  return <div>{day}</div> ;
                };
                return <DayComponent />;
              })()

            } />
            <Route path="/today" element={new Date().toLocaleDateString()} />
            <Route path="/" element={new Date().toLocaleDateString()} />
          </Routes>
          {summary?.duration && <div className="ml-5">({durationAsH(summary?.duration)} h { durationAsM(summary?.duration) } min)</div>}
        </h1>
        <div className="flex w-full justify-end">
          {currentContext?.context.description &&
            <div className="flex rounded-lg p-1 pl-2 pr-2 font-semibold bg-green-200 animate-pulse items-center">
              <div>{currentContext?.context.description} ({durationAsH(currentContext?.currentDuration)} h {durationAsM(currentContext?.currentDuration) } min)</div>
              <Pause className="cursor-pointer shrink-0" onClick={() => freeMutation.mutate({day: day})}></Pause>
            </div>
          }
        </div>
      </div>
    </header>
  )
}
