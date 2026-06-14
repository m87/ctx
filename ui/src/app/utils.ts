export function colorHash(id: string) {
  let hash = 0;
  for (let i = 0; i < id.length; i++) {
    hash = id.charCodeAt(i) + ((hash << 5) - hash);
  }

  const r = (hash >> 0) & 0xff;
  const g = (hash >> 8) & 0xff;
  const b = (hash >> 16) & 0xff;
  return `rgba(${r}, ${g}, ${b}, 1)`;
}

export function durationAsH(duration: number) {
  return Math.floor(duration / 60000000000 / 60);
}

export function durationAsM(duration: number) {
  return Math.floor((duration / 60000000000) % 60);
}

export function durationAsS(duration: number) {
  return Math.floor((duration / 1000000000) % 60);
}

export function durationAsHM(duration: number) {
  const h = durationAsH(duration);
  const m = durationAsM(duration);
  const s = durationAsS(duration);

  if (h > 0) {
    return m > 0 ? `${h}h ${m}m` : `${h}h`;
  }
  if (m > 0) {
    return `${m}m`;
  }
  if (s > 0) {
    return `${s}s`;
  }
  return '0m';
}
