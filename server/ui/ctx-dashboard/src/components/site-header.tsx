import { Separator } from "@/components/ui/separator"
import { SidebarTrigger } from "@/components/ui/sidebar"
import {PauseIcon} from "lucide-react";

export function SiteHeader() {
  return (
    <header className="group-has-data-[collapsible=icon]/sidebar-wrapper:h-12 flex h-12 shrink-0 items-center gap-2 border-b transition-[width,height] ease-linear">
      <div className="flex w-full items-center gap-1 px-4 lg:gap-2 lg:px-6">
        <SidebarTrigger className="-ml-1" />
        <Separator
          orientation="vertical"
          className="mx-2 data-[orientation=vertical]:h-4"
        />
        <h1 className="text-base font-medium">Contexts</h1>
          <div className="flex w-full justify-end">
              <div className="flex rounded-[300px] p-1 pl-4 pr-4 bg-green-200 animate-pulse"><PauseIcon></PauseIcon><div>elo</div></div>
          </div>
      </div>
    </header>
  )
}
