import { VERSION } from "@/app/version"
import {
  SidebarGroup,
  SidebarGroupContent,
  SidebarMenu,
} from "@/components/ui/sidebar"


export function NavBottom() {
  return (
    <SidebarGroup>
      <SidebarGroupContent className="flex flex-col gap-2 opacity-50">
        <SidebarMenu>
          {VERSION}
        </SidebarMenu>
      </SidebarGroupContent>
    </SidebarGroup>
  )
}
