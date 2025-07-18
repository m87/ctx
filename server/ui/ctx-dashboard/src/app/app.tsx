import {AppSidebar} from "@/components/app-sidebar";
import {SiteHeader} from "@/components/site-header";
import {SidebarInset, SidebarProvider} from "@/components/ui/sidebar";
import {Route, Routes} from "react-router-dom";
import {Contexts} from "@/components/contexts";
import TodaySummary from "@/components/today-summary";
import Recent from "@/components/recent";


export function App() {

    return (
        <SidebarProvider>
            <AppSidebar variant="inset"/>
            <SidebarInset>
                <SiteHeader/>
                <div className="flex flex-col h-full flex-1 min-h-0min-h-0">
                    <Routes>
                        <Route path="/recent" element={<Recent/>}/>
                        <Route path="/contexts" element={<Contexts/>}/>
                        <Route path="/today" element={<TodaySummary/>}/>
                        <Route path="/" element={<TodaySummary/>}/>
                    </Routes>
                </div>
            </SidebarInset>
        </SidebarProvider>
    )
}

export default App;
