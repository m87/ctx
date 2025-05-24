import { api } from "@/api/api";
import {
  Card,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { useQuery } from "@tanstack/react-query";
import { AnchorIcon, ArrowDown, ArrowDown01Icon, PlayCircleIcon } from "lucide-react";
import { useState } from "react";
import ContextCard from "./context-card";

export function SectionCards({ search }) {

  const { data: contextList } = useQuery({ ...api.context.listQuery });
  const filteredList = contextList?.filter((context) =>
    context.description.toLowerCase().includes(search.toLowerCase())
  );

  const cardClick = (id) => {
    api.context.switch(id)
  };
  return (
    <div className="*:data-[slot=card]:shadow-xs @xl/main:grid-cols-2 @5xl/main:grid-cols-4 grid grid-cols-1 gap-4 px-4 *:data-[slot=card]:bg-gradient-to-t *:data-[slot=card]:from-primary/5 *:data-[slot=card]:to-card dark:*:data-[slot=card]:bg-card lg:px-6">
      {filteredList?.map((context) => (
        <ContextCard key={context.id} context={context} > </ContextCard>
      ))}
    </div>
  )
}
