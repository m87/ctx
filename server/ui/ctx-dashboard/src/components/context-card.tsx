import { ArrowDown, ChevronDown, ChevronUp, Edit, PlayCircleIcon } from "lucide-react";
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "./ui/card";
import { useEffect, useState } from "react";
import { api, ZonedDateTime, Interval, Context } from "@/api/api";
import IntervalComponent from "./interval-component";
import { Badge } from "./ui/badge";
import { colorHash } from "@/lib/utils";
import { compareAsc } from "date-fns";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useParams } from "react-router-dom";
import { ContextMenu, ContextMenuContent, ContextMenuItem, ContextMenuTrigger } from "./ui/context-menu";
import { Button } from "./ui/button";
import { Dialog, DialogClose, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from "./ui/dialog";
import { Label } from "./ui/label";
import { Input } from "./ui/input";
import { EditContextDialog } from "./edit-context-dialog";

export interface ContextCardProps {
  context: Context
  expandCard: boolean
}


export function ContextCard({ context, expandCard }: ContextCardProps) {
  const [hovered, setHovered] = useState(false);
  const [expanded, setExpand] = useState(false);
  const querClient = useQueryClient();
  const switchMutation = useMutation(api.context.switchMutation(querClient))
  const renameMutation = useMutation(api.context.renameMutation(querClient))
  const updateIntervalMutation = useMutation(api.context.updateIntervalMutation(querClient))
  const { day } = useParams();

  useEffect(() => {
    setExpand(expandCard);
  }, [expandCard])

  return (
    <Card key={context.id} className="flex w-full h-full"
      onMouseEnter={() => setHovered(true)}
      onMouseLeave={() => setHovered(false)}
    >
      <div className="h-full w-2 rounded-l-xl" style={{ backgroundColor: colorHash(context.id) }}></div>
      <div className="@container/card w-full">
        <CardHeader className="relative">
          <CardTitle className="@[250px]/card:text-3xl text-2xl font-semibold tabular-nums flex justify-between w-full items-center">
            <div className="flex w-full items-center ">
              {hovered && <div className="cursor-pointer"><PlayCircleIcon size={30} onClick={() => switchMutation.mutate({ id: context.id, day })} /></div>}
              <div className="flex flex-col items-start">
                <div className="flex flex-grow min-w-0">

                  <div>{context.description} </div>

                  <div className="ml-5 flex-shrink-0 whitespace-nowrap">({Math.floor(context.duration / 60000000000 / 60)} h {Math.floor(context.duration / 60000000000 % 60)} min)</div>
                </div>
                <div className="flex">
                  {context.labels?.map((label: string) => (
                    <Badge variant={"secondary"}>{label}</Badge>
                  ))}
                </div>
              </div>
            </div>
            <div>
              <div onClick={() => setExpand(!expanded)} className="cursor-pointer">
                <ChevronDown
                  className={`transition-transform duration-200 ${expanded ? "rotate-180" : ""}`}
                />
              </div>
            </div>
          </CardTitle>
        </CardHeader>
        {expanded && <CardContent className="flex flex-col gap-2">
          <div className="flex flex-col justify-center">
            {Object.values(context.intervals ?? []).sort((a, b) => compareAsc(a.start.time, b.start.time)).map((interval: Interval) => (
              <IntervalComponent key={interval.id} interval={interval} onChange={(id: string, start: ZonedDateTime, end: ZonedDateTime) => updateIntervalMutation.mutate({ contextId: context.id, intervalId: id, start, end, day })} />
            ))}
          </div>
          <div className="w-full flex justify-end items-end">
            <EditContextDialog context={context} onChange={(ctx) => renameMutation.mutate({ctxId: ctx.id, name: ctx.description})}>
              <Button variant="outline" size="sm"><Edit />Edit context</Button>
            </EditContextDialog>
          </div>
        </CardContent>
        }
      </div>
    </Card>
  )
}

export default ContextCard;
