"use client"

import {
    SidebarGroup,
    SidebarGroupContent,
    SidebarMenu,
} from "@/components/ui/sidebar"
import { Calendar } from "@/components/ui/calendar"
import { useState } from "react";

export function NavMain() {
    const [selected, setSelected] = useState<Date | undefined>(new Date());
    return (
        <SidebarGroup>
            <SidebarGroupContent className="flex flex-col gap-2">
                <SidebarMenu>
                    <Calendar className="p-0" classNames={{
                        day_today: "border border-muted-foreground",
                        day_selected: "bg-primary text-primary-foreground hover:bg-primary hover:text-primary-foreground focus:bg-primary focus:text-primary-foreground",
                    }}
                        selected={selected}
                        onSelect={setSelected}></Calendar>
                </SidebarMenu>
            </SidebarGroupContent>
        </SidebarGroup>
    )
}
