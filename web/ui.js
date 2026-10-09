import { api } from './api.js';
import { esc, bubbles, price } from './card.js';

const $ = (id) => document.getElementById(id);

export function toast(msg) {
  const t = $('toast');
  t.textContent = msg;
  t.hidden = false;
  clearTimeout(toast.timer);
  toast.timer = setTimeout(() => { t.hidden = true; }, 2400);
}

export function showMatch(v) {
  $('match-emoji').textContent = v.emoji;
  $('match-text').textContent =
    `Everyone liked ${v.name}! ${v.cuisine} in ${v.neighborhood} — looks like tonight's plan.`;
  $('match-overlay').hidden = false;
}

export function renderMatches(m) {
  if (!m) return;
  $('crew').textContent = `Swiping: ${m.participants.join(', ')}`;
  const list = $('matches-list');
  if (!m.results.length) {
    list.innerHTML = '<li class="muted empty-list">No likes yet — keep swiping!</li>';
    return;
  }
  const total = m.participants.length;
  list.innerHTML = m.results.map((r) => `
    <li class="match-item ${r.matched ? 'matched' : ''}">
      <div class="mi-emoji cat-${r.venue.category}">${esc(r.venue.emoji)}</div>
      <div class="mi-body">
        <strong>${esc(r.venue.name)}</strong>
        <div class="rating small">${bubbles(r.venue.rating)}
          <span class="muted">${r.venue.reviewCount.toLocaleString()} reviews</span></div>
        <div class="muted">${esc(r.venue.cuisine)} · ${price(r.venue.priceLevel)} · ${esc(r.venue.neighborhood)}</div>
        <div class="likers">♥ ${esc(r.likers.join(', '))}</div>
      </div>
      <div class="mi-count">${r.matched ? '✓ Match' : `${r.likes}/${total}`}</div>
    </li>`).join('');
}

export async function share(code) {
  const url = `${location.origin}/?code=${code}`;
  const text = `Help pick tonight's spot! Join my Night Swipe session with code ${code}`;
  try {
    if (navigator.share) {
      await navigator.share({ title: 'Night Swipe', text, url });
    } else {
      await navigator.clipboard.writeText(url);
      toast('Invite link copied');
    }
  } catch {
    toast(`Share code: ${code}`);
  }
}

export { api };
