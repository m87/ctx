import { Route, Routes } from "react-router-dom";
import { Contexts } from "@/components/contexts";
import TodaySummary from "@/components/today-summary";
import Recent from "@/components/recent";
import AppLayout from "./app-layout";


export function App() {

  return (
    <Routes>
      <Route path="/" element={<AppLayout />}>
        <Route index element={<TodaySummary />} />
        <Route path="/day/:day" element={<TodaySummary />} />
        <Route path="/recent" element={<Recent />} />
        <Route path="/contexts" element={<Contexts />} />
        <Route path="/today" element={<TodaySummary />} />
      </Route>
    </Routes>
  )
}

export default App;
