import { ChevronDown, Clock, Delete, Edit, MessageSquareText, PlayCircleIcon, PlayIcon, Tag, Trash } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "./ui/card";
import { useEffect, useState } from "react";
import { api, Context } from "@/api/api";
import { Badge } from "./ui/badge";
import { colorHash } from "@/lib/utils";
import { compareAsc } from "date-fns";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useParams } from "react-router-dom";
import { IntervalTable } from "./intervals-table";
import { Button } from "./ui/button";
import { RenameContextDialog } from "./dialogs/rename-context-dialog";
import { Separator } from "./ui/separator";
import { DeleteContextDialog } from "./dialogs/delete-context-dialog";

export interface ContextCardProps {
  context: Context;
  expandCard: boolean;
}

export function ContextCard({ context, expandCard }: ContextCardProps) {
  const [hovered, setHovered] = useState(false);
  const [expanded, setExpand] = useState(false);
  const queryClient = useQueryClient();
  const switchMutation = useMutation(api.context.switchMutation(queryClient));
  const { day } = useParams();

  useEffect(() => {
    setExpand(expandCard);
  }, [expandCard]);

  const hours = Math.floor(context.duration / 60000000000 / 60);
  const minutes = Math.floor((context.duration / 60000000000) % 60);

  return (
    <Card
      key={context.id}
      className="flex w-full transition-all hover:shadow-md hover:scale-[1.01] relative"
      style={{ borderLeftColor: "rgba(0,0,0,0)" }}
      onMouseEnter={() => setHovered(true)}
      onMouseLeave={() => setHovered(false)}
    >
      <div
        className="h-full w-2 rounded-l-xl"
        style={{ backgroundColor: colorHash(context.id) }}
      />

      <div className="relative w-full">
        <div
          className={`@container/card w-full transition-[padding] duration-200`}
        >
          <CardHeader className="relative p-3">
            <CardTitle className="@[250px]/card:text-3xl text-2xl font-semibold tabular-nums flex justify-between items-center w-full">
              <div className="flex items-center gap-4 w-full">

                <div className="flex flex-col min-w-0">
                  <div className="flex justify-between w-full items-center">
                    <div className="flex items-end">
                      <span className="truncate font-medium">{context.description}</span>
                    </div>
                    <div className="flex items-center text-sm text-muted-foreground ml-3 whitespace-nowrap gap-1">
                      <Clock size={16}></Clock>
                      <span>
                        {hours > 0 && `${hours} h `}{minutes} min
                      </span>
                    </div>

                                <div className="flex items-center text-sm text-muted-foreground ml-3 whitespace-nowrap gap-1">
                      <Tag size={16}></Tag>
                      <span>
                          5
                      </span>
                    </div>

                                <div className="flex items-center text-sm text-muted-foreground ml-3 whitespace-nowrap gap-1">
                      <MessageSquareText size={16}></MessageSquareText>
                      <span>
                        2
                      </span>
                    </div>
                  </div>

                  {context.labels?.length > 0 && (
                    <div className="flex flex-wrap gap-1 mt-1">
                      {context.labels.map((label: string) => (
                        <Badge key={label} variant="secondary">
                          {label}
                        </Badge>
                      ))}
                    </div>
                  )}
                </div>
               <div className="flex gap-2">


               </div>

              </div>

              <div
                className={`transition-opacity duration-200 ${hovered ? "opacity-100" : "opacity-0 pointer-events-none"
                  }`}
              >
                <Button variant="outline" size="sm" className="flex items-center gap-1"
                  onClick={() => switchMutation.mutate({ id: context.id, day })}
                >
                  <PlayIcon size={16} /> Switch context
                </Button>
              </div>

              <div
                onClick={() => setExpand(!expanded)}
                className="cursor-pointer ml-3 p-1 hover:bg-muted rounded-md transition"
              >
                <ChevronDown
                  className={`transition-transform duration-200 ${expanded ? "rotate-180" : ""}`}
                />
              </div>
            </CardTitle>
          </CardHeader>

          {expanded && (
            <CardContent className="flex flex-col gap-3">
              <div className="flex gap-2"><Badge variant={"secondary"}>asd</Badge>
              <Badge variant={"secondary"}>asd</Badge>
              <Badge variant={"secondary"}>asd</Badge>
              <Badge variant={"secondary"}>asd</Badge>
              <Badge variant={"secondary"}>asd</Badge>
              <Badge variant={"secondary"}>asd</Badge></div>
              
              <IntervalTable
                ctxId={context.id}
                intervals={Object.values(context.intervals ?? []).sort((a, b) =>
                  compareAsc(a.start.time, b.start.time)
                )}
              />
               <div className="flex flex-col">
                  <h4 className="font-medium">Comments</h4>
                  <div>asdlsajl  d lsajd  lksaj s lkdsalk djlksajd lkdsajlksa djdl</div>
              </div>
              <Separator />
              <div className="flex mb-2 gap-1 w-full justify-end">
                <RenameContextDialog context={context}>
                  <Button variant="outline" size="sm" className="flex items-center gap-1">
                    <Edit size={16} /> Rename
                  </Button>
                </RenameContextDialog>
                <DeleteContextDialog context={context}>
                <Button variant="destructive" size="sm" className="flex items-center gap-1">
                  <Trash size={16} /> Delete
                </Button>
                </DeleteContextDialog>
              </div>
             
            </CardContent>
          )}
        </div>
      </div>
    </Card>
  );
}

export default ContextCard;
