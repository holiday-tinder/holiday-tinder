async function request(method, path, body) {
  const res = await fetch(path, {
    method,
    headers: body ? { 'Content-Type': 'application/json' } : {},
    body: body ? JSON.stringify(body) : undefined,
  });
  const data = await res.json().catch(() => ({}));
  if (!res.ok) throw new Error(data.error || `Request failed (${res.status})`);
  return data;
}

const s = (code) => `/api/sessions/${encodeURIComponent(code)}`;

export const api = {
  cities: () => request('GET', '/api/cities'),
  createSession: (name, city, categories) => request('POST', '/api/sessions', { name, city, categories }),
  join: (code, name) => request('POST', `${s(code)}/join`, { name }),
  deck: (code, pid) => request('GET', `${s(code)}/deck?participant=${encodeURIComponent(pid)}`),
  swipe: (code, pid, venueId, liked) =>
    request('POST', `${s(code)}/swipes`, { participantId: pid, venueId, liked }),
  matches: (code) => request('GET', `${s(code)}/matches`),
};
