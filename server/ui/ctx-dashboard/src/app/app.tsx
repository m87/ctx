
import { Button } from '@/components/ui/button';
import Timeline from '@/components/ui/timeline';
import { Link } from 'react-router-dom';

export function App() {
  return (
    <div>
      <Button variant="primary">
        <Link to="/dashboard">Dashboard</Link>
      </Button>
      <Button variant="secondary">
        <Link to="/settings">Settings</Link>
      </Button>
      <Timeline
        data={[
          {
            date: "2024-05-18",
            blocks: [
              { label: "Spotkanie", start: 9, end: 11, color: "#3b82f6" },
              { label: "Obiad", start: 13, end: 14.5, color: "#10b981" }
            ]
          },
          {
            date: "2024-05-19",
            blocks: [
              { label: "Praca", start: 8, end: 16, color: "#f97316" }
            ]
          }
        ]}
/>

    </div>
  );
}

export default App;
