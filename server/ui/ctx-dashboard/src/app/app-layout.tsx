import {AppSidebar} from "@/components/app-sidebar";
import {SiteHeader} from "@/components/site-header";
import {SidebarInset, SidebarProvider} from "@/components/ui/sidebar";
import {Outlet} from "react-router-dom";


export function AppLayout() {

    return (
        <SidebarProvider>
            <AppSidebar variant="inset"/>
            <SidebarInset>
                <SiteHeader/>
                <div className="flex flex-col h-full flex-1 min-h-0min-h-0">
                <Outlet/>
                </div>

            </SidebarInset>
        </SidebarProvider>
    )
}

export default AppLayout;
