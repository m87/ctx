import { Route, Routes } from "react-router-dom";
import TodaySummary from "@/components/today-summary";
import AppLayout from "./app-layout";


export function App() {

  return (
    <Routes>
      <Route path="/" element={<AppLayout />}>
        <Route index element={<TodaySummary />} />
        <Route path="/day/:day" element={<TodaySummary />} />
      </Route>
    </Routes>
  )
}

export default App;
