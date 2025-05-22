
import Timeline from '@/components/ui/timeline';
import {SidebarInset, SidebarProvider} from "@/components/ui/sidebar";
import {AppSidebar} from "@/components/app-sideba2r";
import { api } from '@/api/api';
import { useSuspenseQuery } from '@tanstack/react-query';

export function App2() {
    const {data: contexts} = useSuspenseQuery(api.context.contextListQuery);

const timelineData = 
  {
    date: "2025-05-19",
    blocks: [
      { label: "Spotkanie", start: 9, end: 10.5, color: "#3b82f6" },
      { label: "Lunch", start: 12.5, end: 13.25, color: "#10b981" },
    ],
  }
    return (
        <SidebarProvider
            style={
                {
                    "--sidebar-width": "350px",
                } as React.CSSProperties
            }
        >
            <AppSidebar contexts={contexts}/>
            <SidebarInset>
                {/* <header className="sticky top-0 flex shrink-0 items-center gap-2 border-b bg-background p-4">
                    <SidebarTrigger className="-ml-1" />
                    <Separator orientation="vertical" className="mr-2 h-4" />
                </header> */}
        <div className="p-0 overflow-y-auto overflow-x-auto w-full h-full flex-1">
            <div className="min-w-[200px]">
            <div className="inset-0 flex border-l border-r sticky top-0 bg-background z-10">
                {Array.from({ length: 25 }, (_, h) => (
                    <div
                        key={h}
                          className="flex-1 border-r text-xs text-center text-muted-foreground p-2"
                    >
                      {String(h).padStart(2, "0")}:00
                    </div>
                    
                ))}
              </div>
                <div className="flex flex-1 flex-col">
                    {Array.from({ length: 24 }).map((_, index) => (

                        <div
                            key={index}
                            className="aspect-video h-12 w-full rounded-lg bg-muted/50"
                        > <Timeline date="2025-05-19" blocks={[      { label: "Spotkanie", start: 9, end: 10.5, color: "#3b82f6" },
      { label: "Lunch", start: 12.5, end: 13.25, color: "#10b981" },]} /></div>
                        
                    ))}
                </div>
                </div>
                </div>
            </SidebarInset>
        </SidebarProvider>
    )
}

export default App2;
