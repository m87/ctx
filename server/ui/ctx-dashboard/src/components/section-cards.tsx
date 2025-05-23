import { api } from "@/api/api";
import {
  Card,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { useQuery } from "@tanstack/react-query";
import { AnchorIcon, ArrowDown, ArrowDown01Icon, PlayCircleIcon } from "lucide-react";
import { useState } from "react";

export function SectionCards({ search }) {
  const [hovered, setHovered] = useState('');

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
        <Card key={context.id} className="@container/card"
          onMouseEnter={() => setHovered(context.id)}
          onMouseLeave={() => setHovered('')}
        >
          <CardHeader className="relative">
            <CardTitle className="@[250px]/card:text-3xl text-2xl font-semibold tabular-nums flex justify-between w-full items-center">
              <div className="flex w-full items-center">
                {hovered == context.id && <div className="cursor-pointer"><PlayCircleIcon size={50} onClick={() => cardClick(context.id)} /></div>}
                <div>    {context.description} </div>
                <div className="flex items-center">
                  <div className="text-xs text-muted-foreground">{context.intervals.length} intervals</div>
                  {context.intervals?.map((interval, index) => (
                    <div>elo</div>
                  ))}
                  </div>
              </div>
              <div><ArrowDown></ArrowDown></div>
            </CardTitle>
          </CardHeader>
        </Card>
      ))}
    </div>
  )
}
