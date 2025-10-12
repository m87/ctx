import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function colorHash(id: string) {
  let hash = 0;
  for (let i = 0; i < id.length; i++) {
    hash = id.charCodeAt(i) + ((hash << 5) - hash);
  }

  const r = (hash >> 0) & 0xFF;
  const g = (hash >> 8) & 0xFF;
  const b = (hash >> 16) & 0xFF;
  return `rgba(${r}, ${g}, ${b}, 1)`;
}

export function durationAsH(duration: number) {
  return Math.floor(duration / 60000000000 / 60)
}

export function durationAsM(duration: number) {
  return Math.floor(duration / 60000000000 % 60)
}

export function durationAsS(duration: number) {
  return Math.floor(duration / 1000000000 % 60)
}

export function durationAsHM(duration: number) {
  return `${durationAsH(duration)}h ${durationAsM(duration)}m`
}
