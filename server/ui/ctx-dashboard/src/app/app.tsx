
import { Button } from '@/components/ui/button';
import Timeline from '@/components/ui/timeline';
import { Link } from 'react-router-dom';
import ContextList from "@/components/context-list";
import {SidebarInset, SidebarProvider, SidebarTrigger} from "@/components/ui/sidebar";
import {AppSidebar} from "@/components/app-sidebar";
import {Separator} from "@/components/ui/separator";

export function App() {
    return (
        <SidebarProvider
            style={
                {
                    "--sidebar-width": "350px",
                } as React.CSSProperties
            }
        >
            <AppSidebar />
            <SidebarInset>
                <header className="sticky top-0 flex shrink-0 items-center gap-2 border-b bg-background p-4">
                    <SidebarTrigger className="-ml-1" />
                    <Separator orientation="vertical" className="mr-2 h-4" />
                </header>
                <div className="flex flex-1 flex-col gap-4 p-4">
                    {Array.from({ length: 24 }).map((_, index) => (
                        <div
                            key={index}
                            className="aspect-video h-12 w-full rounded-lg bg-muted/50"
                        />
                    ))}
                </div>
            </SidebarInset>
        </SidebarProvider>
    )
}

export default App;
