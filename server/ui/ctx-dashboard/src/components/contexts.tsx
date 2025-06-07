import {Card, CardDescription, CardHeader, CardTitle} from "@/components/ui/card";
import {PauseCircle} from "lucide-react";
import {SectionCards} from "@/components/section-cards";
import {useQuery} from "@tanstack/react-query";
import {api} from "@/api/api";


export function Contexts() {

    const { data: contextList } = useQuery({ ...api.context.listQuery });
    const { data: currentContext } = useQuery({...api.context.currentQuery, refetchInterval: 1000});
    const pauseClick = () => {
        api.context.free()
    };
    return <>
        <div className="@container/main flex flex-col gap-2 h-full flex-1 min-h-0min-h-0">
            <div className="flex flex-col gap-4 py-4 h-full flex-1 min-h-0">
                {currentContext &&
                    <div
                        className="*:data-[slot=card]:shadow-xs @xl/main:grid-cols-2 @5xl/main:grid-cols-4 grid grid-cols-1 gap-4 px-4 *:data-[slot=card]:bg-gradient-to-t *:data-[slot=card]:from-primary/5 *:data-[slot=card]:to-card dark:*:data-[slot=card]:bg-card lg:px-6">
                        <Card className="@container/card bg-slate-100 animate-pulse">
                            <CardHeader className="relative ">
                                <CardDescription className="flex justify-between">
                                    <div className="text-black">Current context</div>
                                    <div className="text-black">12:23:12</div>
                                </CardDescription>
                                <CardTitle
                                    className="@[250px]/card:text-3xl text-2xl font-semibold tabular-nums flex justify-between w-full">
                                    <div>{currentContext?.description}</div>
                                    <div><PauseCircle className="size-10" onClick={pauseClick}></PauseCircle></div>
                                </CardTitle>
                            </CardHeader>
                        </Card>
                    </div>
                }

                <div className="h-[100px]">
                    <SectionCards contextList={contextList}></SectionCards>
                </div>

            </div>
        </div>
    </>
}
