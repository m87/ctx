import { api } from "@/api/api";
import {
  Card,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { useQuery } from "@tanstack/react-query";

export function SectionCards() {
  const { data: contextList } = useQuery({...api.context.contextListQuery, refetchInterval: 10000});

  return (
    <div className="*:data-[slot=card]:shadow-xs @xl/main:grid-cols-2 @5xl/main:grid-cols-4 grid grid-cols-1 gap-4 px-4 *:data-[slot=card]:bg-gradient-to-t *:data-[slot=card]:from-primary/5 *:data-[slot=card]:to-card dark:*:data-[slot=card]:bg-card lg:px-6">
      {contextList?.map((context) => (
        <Card key={context.id} className="@container/card">
          <CardHeader className="relative">
            <CardTitle className="@[250px]/card:text-3xl text-2xl font-semibold tabular-nums flex justify-between w-full">
              <div>{context.description}</div>
            </CardTitle>
          </CardHeader>
          </Card>
      ))}
    </div>
  )
}
