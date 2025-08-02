"use client"

import {
    SidebarGroup,
    SidebarGroupContent,
    SidebarMenu,
} from "@/components/ui/sidebar"
import { useState } from "react";
import { ContextCalendar } from "./ui/context-calendar";
import { useQuery } from "@tanstack/react-query";
import { api } from "@/api/api";

export function NavMain() {
    const [selected, setSelected] = useState<Date | undefined>(new Date());
    const { data } = useQuery(api.summary.dayListSummaryQuery);
    return (
        <SidebarGroup>
            <SidebarGroupContent className="flex flex-col gap-2">
                <SidebarMenu>
                    <ContextCalendar 
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
