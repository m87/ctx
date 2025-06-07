import { useQuery } from "@tanstack/react-query";
import ContextCard from "./context-card";
import { api } from "@/api/api";
import { SectionCards } from "./section-cards";


export function TodaySummary() {
    const {data: summary} = useQuery({...api.summary.todaySummaryQuery});

    return (
        <div className="flex flex-col">
            <div className="flex-1 flex items-center justify-center">
            </div>
            <SectionCards contextList={summary?.contexts}></SectionCards>
        </div>
    );
}

export default TodaySummary;