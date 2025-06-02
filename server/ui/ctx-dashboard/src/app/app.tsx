import {AppSidebar} from "@/components/app-sidebar";
import {SiteHeader} from "@/components/site-header";
import {SidebarInset, SidebarProvider} from "@/components/ui/sidebar";
import {Route, Routes} from "react-router-dom";
import {Contexts} from "@/components/contexts";


export function App() {

    return (
        <SidebarProvider>
            <AppSidebar variant="inset"/>
            <SidebarInset>
                <SiteHeader/>
                <div className="flex flex-col h-full flex-1 min-h-0min-h-0">
                    <Routes>
                        <Route path="/contexts" element={<Contexts/>}/>
                    </Routes>

                </div>
            </SidebarInset>
        </SidebarProvider>
    )
}

export default App;
