import React, { useEffect, useState } from "react";
import {
  Clock1, Clock2, Clock3, Clock4, Clock5, Clock6,
  Clock7, Clock8, Clock9, Clock10, Clock11, Clock12
} from "lucide-react";

const ICONS = [
  Clock1, Clock2, Clock3, Clock4, Clock5, Clock6,
  Clock7, Clock8, Clock9, Clock10, Clock11, Clock12,
];

export function TickingClock({ size = 16 }: { size?: number }) {
  const [i, setI] = useState(0);

  useEffect(() => {
    const id = setInterval(() => setI((x) => (x + 1) % 12), 1000);
    return () => clearInterval(id);
  }, []);

  const Icon = ICONS[i];
  return <Icon size={size} />;
}