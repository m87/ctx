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
  const hStr = h > 0 ? `${h}h` : '';
  const mStr = m > 0 ? `${m}m` : '';
  return `${hStr} ${mStr}`;
}
