import { AppSidebar } from "@/components/app-sidebar";
import { SectionCards } from "@/components/section-cards";
import { SiteHeader } from "@/components/site-header";
import { Card, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";
import { Badge, PauseCircle, PlayIcon, TrendingUpIcon } from "lucide-react";



export function App() {
    return (
        <SidebarProvider>
            <AppSidebar variant="inset" />
            <SidebarInset>
                <SiteHeader />
                <div className="flex flex-1 flex-col h-full">
                    <div className="@container/main flex flex-1 flex-col gap-2 h-full">
                        <div className="flex flex-col gap-4 py-4 h-full">
                            <div className="*:data-[slot=card]:shadow-xs @xl/main:grid-cols-2 @5xl/main:grid-cols-4 grid grid-cols-1 gap-4 px-4 *:data-[slot=card]:bg-gradient-to-t *:data-[slot=card]:from-primary/5 *:data-[slot=card]:to-card dark:*:data-[slot=card]:bg-card lg:px-6">
                                <Card className="@container/card bg-slate-100 animate-pulse">

                                    <CardHeader className="relative ">
                                        <CardDescription className="flex justify-between"><div className="text-black">Current context</div><div className="text-black">12:23:12</div></CardDescription>
                                        <CardTitle className="@[250px]/card:text-3xl text-2xl font-semibold tabular-nums flex justify-between w-full">
                                            <div>JIRA-1233 elo</div>
                                            <div><PauseCircle className="size-10"></PauseCircle></div>
                                        </CardTitle>
                                    </CardHeader>
                                </Card>
                            </div>
                            <div className="pt-2 pb-2 pr-6 pl-6">
                                <Input type="text"
                                    placeholder="Search or create new..."></Input>
                            </div>
                            <div className="overflow-y-auto h-full">
                                <SectionCards />
                            </div>
                        </div>
                    </div>
                </div>
            </SidebarInset>
        </SidebarProvider>
    )
}

export default App;