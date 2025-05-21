"use client"

import * as React from "react"
import { ArchiveX, Command, File, Inbox, Send, Trash2 } from "lucide-react"

import { NavUser } from "@/components/nav-user"
import { Label } from "@/components/ui/label"
import {
    Sidebar,
    SidebarContent,
    SidebarFooter,
    SidebarGroup,
    SidebarGroupContent,
    SidebarHeader,
    SidebarInput,
    SidebarMenu,
    SidebarMenuButton,
    SidebarMenuItem,
    useSidebar,
} from "@/components/ui/sidebar"
import { Switch } from "@/components/ui/switch"

// This is sample data
const data = {
    user: {
        name: "shadcn",
        email: "m@example.com",
        avatar: "/avatars/shadcn.jpg",
    },
    navMain: [
    ],
    mails: [
        {
            subject: "Meeting Tomorrow",
        },
        {
            subject: "Re: Project Update",
        },
        {
            subject: "Weekend Plans",
        },
        {
            subject: "Re: Question about Budget",
        },
        {
            subject: "Important Announcement",
        },
        {
            subject: "Re: Feedback on Proposal",
        },
        {
            subject: "New Project Idea",
        },
        {
            subject: "Vacation Plans",
        },
        {
            subject: "Re: Conference Registration",
        },
        {
            subject: "Team Dinner",
        },
        {
            subject: "Meeting Tomorrow",
        },
        {
            subject: "Re: Project Update",
        },
        {
            subject: "Weekend Plans",
        },
        {
            subject: "Re: Question about Budget",
        },
        {
            subject: "Important Announcement",
        },
        {
            subject: "Re: Feedback on Proposal",
        },
        {
            subject: "New Project Idea",
        },
        {
            subject: "Vacation Plans",
        },
        {
            subject: "Re: Conference Registration",
        },
        {
            subject: "Team Dinner",
        },
    ],
}

export function AppSidebar({ contexts }: React.ComponentProps<typeof Sidebar>) {
    // Note: I'm using state to show active item.
    // IRL you should use the url/router.
    const [activeItem, setActiveItem] = React.useState(data.navMain[0])
    const [mails, setMails] = React.useState(data.mails)
    const { setOpen } = useSidebar()

    return (
        <Sidebar
            collapsible=""
            className="overflow-hidden [&>[data-sidebar=sidebar]]:flex-row"
        >
            {/* This is the first sidebar */}
            {/* We disable collapsible and adjust width to icon. */}
            {/* This will make the sidebar appear as icons. */}
            <Sidebar
                collapsible="none"
                className="!w-[calc(var(--sidebar-width-icon)_+_1px)] border-r"
            >
                <SidebarHeader>
                    <SidebarMenu>
                        <SidebarMenuItem>
                            <SidebarMenuButton size="lg" asChild className="md:h-8 md:p-0">
                                <a href="#">
                                    <div className="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground">
                                        <Command className="size-4" />
                                    </div>
                                    <div className="grid flex-1 text-left text-sm leading-tight">
                                        <span className="truncate font-semibold">Acme Inc</span>
                                        <span className="truncate text-xs">Enterprise</span>
                                    </div>
                                </a>
                            </SidebarMenuButton>
                        </SidebarMenuItem>
                    </SidebarMenu>
                </SidebarHeader>
                <SidebarContent>
                    <SidebarGroup>
                        <SidebarGroupContent className="px-1.5 md:px-0">
                            <SidebarMenu>
                                {data.navMain.map((item) => (
                                    <SidebarMenuItem key={item.title}>
                                        <SidebarMenuButton
                                            tooltip={{
                                                children: item.title,
                                                hidden: false,
                                            }}
                                            onClick={() => {
                                                setActiveItem(item)
                                                const mail = data.mails.sort(() => Math.random() - 0.5)
                                                setMails(
                                                    mail.slice(
                                                        0,
                                                        Math.max(5, Math.floor(Math.random() * 10) + 1)
                                                    )
                                                )
                                                setOpen(true)
                                            }}
                                            isActive={activeItem?.title === item.title}
                                            className="px-2.5 md:px-2"
                                        >
                                            <item.icon />
                                            <span>{item.title}</span>
                                        </SidebarMenuButton>
                                    </SidebarMenuItem>
                                ))}
                            </SidebarMenu>
                        </SidebarGroupContent>
                    </SidebarGroup>
                </SidebarContent>
                <SidebarFooter>
                    <NavUser user={data.user} />
                </SidebarFooter>
            </Sidebar>

            {/* This is the second sidebar */}
            {/* We disable collapsible and let it fill remaining space */}
            <Sidebar collapsible="none" className="hidden flex-1 md:flex">
                <SidebarHeader className="gap-3.5 border-b p-4">
                    <SidebarInput placeholder="Type to search..." />
                </SidebarHeader>
                <SidebarContent>
                    <SidebarGroup className="px-0">
                        <SidebarGroupContent>
                            {contexts.map((mail) => (
                                <a
                                    href="#"
                                    key={mail.id}
                                    className="flex flex-col items-start gap-2 whitespace-nowrap border-b p-4 text-sm leading-tight last:border-b-0 hover:bg-sidebar-accent hover:text-sidebar-accent-foreground"
                                >
                                    <span className="font-medium">{mail.description}</span>
                                </a>
                            ))}
                        </SidebarGroupContent>
                    </SidebarGroup>
                </SidebarContent>
            </Sidebar>
        </Sidebar>
    )
}
