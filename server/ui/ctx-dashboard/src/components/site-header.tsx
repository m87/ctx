import {api} from "@/api/api";
import {Separator} from "@/components/ui/separator"
import {SidebarTrigger} from "@/components/ui/sidebar"
import {useQuery} from "@tanstack/react-query";
import {Route, Routes} from "react-router-dom";

export function SiteHeader() {

    const {data: currentContext} = useQuery({...api.context.currentQuery, refetchInterval: 5000});

    return (
        <header
            className="group-has-data-[collapsible=icon]/sidebar-wrapper:h-12 flex h-12 shrink-0 items-center gap-2 border-b transition-[width,height] ease-linear">
            <div className="flex w-full items-center gap-1 px-4 lg:gap-2 lg:px-6">
                <SidebarTrigger className="-ml-1"/>
                <Separator
                    orientation="vertical"
                    className="mx-2 data-[orientation=vertical]:h-4"
                />
                <h1 className="text-base font-medium">
                    <Routes>
                        <Route path="/contexts" element={"Contexts"}/>
                        <Route path="/today" element={new Date().toLocaleDateString()}/>
                        <Route path="/" element={new Date().toLocaleDateString()}/>
                    </Routes>
                </h1>
                <div className="flex w-full justify-end">
                    <div className="flex rounded-lg p-1 pl-2 pr-2 font-semibold bg-green-200 animate-pulse">
                        <div>{currentContext?.description}</div>
                    </div>
                </div>
            </div>
        </header>
    )
}
