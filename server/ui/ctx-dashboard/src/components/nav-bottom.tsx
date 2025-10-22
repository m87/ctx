import { api } from "@/api/api";
import {
  SidebarGroup,
  SidebarGroupContent,
  SidebarMenu,
} from "@/components/ui/sidebar"
import { useQuery } from "@tanstack/react-query";
import { Card } from "./ui/card";


export function NavBottom() {
   const { data: version } = useQuery({ ...api.versionQuery }); 

  return (
    <SidebarGroup>
      <SidebarGroupContent className="flex flex-col gap-2 opacity-50">
        {/* <SidebarMenu>
          {version}
        </SidebarMenu> */}
      </SidebarGroupContent>
    </SidebarGroup>
  )
}
