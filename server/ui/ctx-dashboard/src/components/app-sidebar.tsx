"use client"

import * as React from "react"
import {
  ArrowUpCircleIcon,
  BarChartIcon, Calendar1Icon,
  CalendarRange,
  CameraIcon,
  ChartNoAxesGantt,
  ChartNoAxesGanttIcon,
  ClipboardListIcon,
  Clock,
  DatabaseIcon,
  FileCodeIcon,
  FileTextIcon,
  HelpCircleIcon,
  LayoutDashboardIcon,
  Pause,
  SearchIcon,
  SettingsIcon,
} from "lucide-react"

import { NavMain } from "@/components/nav-main"
import { NavSecondary } from "@/components/nav-secondary"
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar"
import { NavBottom } from "./nav-bottom"
import { Card } from "./ui/card"
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { api } from "@/api/api"
import { useParams } from "react-router-dom"
import { isValid, parseISO } from "date-fns"
import { durationAsHM } from "@/lib/utils"
import { IconProgress } from "@tabler/icons-react"
import { Spinner } from "./ui/spinner"
import { Button } from "./ui/button"

const data = {
  user: {
    name: "shadcn",
    email: "m@example.com",
    avatar: "/avatars/shadcn.jpg",
  },
  navMain: [
    {
      title: "Today",
      url: "",
      icon: Calendar1Icon
    },
    // {
    //   title: "Recent",
    //   url: "recent",
    //   icon: CalendarRange
    // },
    //{
    //  title: "Contexts",
    //  url: "contexts",
    //  icon: LayoutDashboardIcon,
    //},
  ],
  navClouds: [
    {
      title: "Capture",
      icon: CameraIcon,
      isActive: true,
      url: "#",
      items: [
        {
          title: "Active Proposals",
          url: "#",
        },
        {
          title: "Archived",
          url: "#",
        },
      ],
    },
    {
      title: "Proposal",
      icon: FileTextIcon,
      url: "#",
      items: [
        {
          title: "Active Proposals",
          url: "#",
        },
        {
          title: "Archived",
          url: "#",
        },
      ],
    },
    {
      title: "Prompts",
      icon: FileCodeIcon,
      url: "#",
      items: [
        {
          title: "Active Proposals",
          url: "#",
        },
        {
          title: "Archived",
          url: "#",
        },
      ],
    },
  ],
  navSecondary: [
    {
      title: "Settings",
      url: "#",
      icon: SettingsIcon,
    },
    {
      title: "Get Help",
      url: "#",
      icon: HelpCircleIcon,
    },
    {
      title: "Search",
      url: "#",
      icon: SearchIcon,
    },
  ],
  documents: [
  ],
}

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {

  const [selectedDate, setSelectedDate] = React.useState<Date>(new Date());
  const { data: currentContext } = useQuery({ ...api.context.currentQuery, refetchInterval: 5000 });
  const { data: version } = useQuery({ ...api.versionQuery });
  const querClient = useQueryClient()
  const freeMutation = useMutation(api.context.freeMutaiton(querClient))
  const { day } = useParams();


  React.useEffect(() => {
    if (day) {
      const date = parseISO(day)
      if (isValid(date)) {
        setSelectedDate(date)
      }
    }
  }, [day])


  return (
    <Sidebar collapsible="offcanvas" {...props}>
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton
              asChild
            >
              <div className="flex justify-between items-center w-full">
                <div>
                  <a href="#" className="gap-0 flex gap-2 items-center">
                    <ChartNoAxesGanttIcon className="h-5 w-5" />
                    <span className="text-base font-semibold">Ctx</span>
                  </a>
                </div>
                <div className="text-muted-foreground">
                  <span>{version}</span>
                </div>
              </div>

            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent>
        <NavMain />
      </SidebarContent>
      <SidebarFooter>
        <NavBottom />
      </SidebarFooter>

    </Sidebar>
  )
}
