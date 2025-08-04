import {
  SidebarGroup,
  SidebarGroupContent,
  SidebarMenu,
} from "@/components/ui/sidebar"
import { useEffect, useState } from "react";
import { ContextCalendar } from "./ui/context-calendar";
import { useQuery } from "@tanstack/react-query";
import { api } from "@/api/api";
import { useNavigate, useParams } from "react-router-dom";
import { format, isValid, parseISO } from "date-fns";

export function NavMain() {
  const [selected, setSelected] = useState<Date | undefined>(new Date());
  const { data } = useQuery(api.summary.dayListSummaryQuery);
  const navigate = useNavigate();
  const { day } = useParams();

  useEffect(() => {
    if (day) {
      const date = parseISO(day)
      if (isValid(date)) {
        setSelected(date)
      }
    }

  }, [day])

  return (
    <SidebarGroup>
      <SidebarGroupContent className="flex flex-col gap-2">
        <SidebarMenu>
          <ContextCalendar
            onClick={(date) => navigate(`/day/${format(date, 'yyyy-MM-dd')}`)}
            contextMap={data}
            className="p-0" classNames={{
              day_today: "border border-muted-foreground",
              day_selected: "bg-primary text-primary-foreground hover:bg-primary hover:text-primary-foreground focus:bg-primary focus:text-primary-foreground",
            }}
            selected={selected}
            onSelect={setSelected}></ContextCalendar>
        </SidebarMenu>
      </SidebarGroupContent>
    </SidebarGroup>
  )
}
