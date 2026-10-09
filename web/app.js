import { api } from './api.js';
import { esc, renderCard } from './card.js';
import { makeSwipeable, fling } from './swipe.js';
import { toast, showMatch, renderMatches, share } from './ui.js';

const $ = (id) => document.getElementById(id);
const STORE = 'nightswipe';
const state = { code: null, pid: null, deck: [], seen: new Set(), poll: null, last: null };

const showError = (msg) => { const e = $('home-error'); e.textContent = msg; e.hidden = !msg; };
const saved = () => { try { return JSON.parse(localStorage.getItem(STORE)) || {}; } catch { return {}; } };

async function enterSession(code, pid) {
  const { venues } = await api.deck(code, pid);
  Object.assign(state, { code, pid, deck: venues, seen: new Set(), last: null });
  localStorage.setItem(STORE, JSON.stringify({ code, pid }));
  history.replaceState(null, '', `/?code=${code}`);
  $('home').hidden = true;
  $('swipe').hidden = false;
  $('session-code').textContent = code;
  renderDeck();
  await refresh(true);
  clearInterval(state.poll);
  state.poll = setInterval(() => refresh(false), 5000);
}

function leave() {
  if (!confirm('Leave this session? You will need to rejoin as a new swiper.')) return;
  clearInterval(state.poll);
  localStorage.removeItem(STORE);
  history.replaceState(null, '', '/');
  state.code = state.pid = null;
  $('swipe').hidden = true;
  $('home').hidden = false;
}

function renderDeck() {
  const deck = $('deck');
  deck.querySelectorAll('.card:not(.gone), .empty').forEach((n) => n.remove());
  if (!state.deck.length) {
    deck.insertAdjacentHTML('beforeend', `<div class="empty"><div class="empty-emoji">🎉</div>
      <h2>You've seen every spot</h2><p class="muted">See what your crew liked.</p>
      <button class="btn primary" id="empty-matches">See picks</button></div>`);
    $('empty-matches').onclick = openMatches;
    return;
  }
  const top = state.deck.slice(0, 3);
  for (let depth = top.length - 1; depth >= 0; depth--) {
    const v = top[depth];
    const card = renderCard(v);
    card.style.setProperty('--depth', depth);
    if (depth === 0) {
      card.classList.add('top');
      makeSwipeable(card, (liked) => decide(v, liked));
    }
    deck.appendChild(card);
  }
}

async function decide(venue, liked) {
  state.deck.shift();
  renderDeck();
  try {
    const { match } = await api.swipe(state.code, state.pid, venue.id, liked);
    if (match) { state.seen.add(match.id); showMatch(match); }
    if (liked) refresh(false);
  } catch (e) { toast(e.message); }
}

async function refresh(initial) {
  if (!state.code) return;
  try {
    const m = await api.matches(state.code);
    state.last = m;
    const n = m.participants.length;
    $('session-info').textContent = `📍 ${m.city} · ${n} ${n === 1 ? 'person' : 'people'} swiping`
      + (n < 2 ? ' · share the code!' : '');
    const matched = m.results.filter((r) => r.matched);
    const badge = $('match-badge');
    badge.textContent = matched.length;
    badge.hidden = !matched.length;
    for (const r of matched) {
      if (state.seen.has(r.venue.id)) continue;
      state.seen.add(r.venue.id);
      if (!initial) showMatch(r.venue);
    }
    if (!$('matches-sheet').hidden) renderMatches(m);
  } catch { /* transient; next poll retries */ }
}

function openMatches() {
  $('matches-sheet').hidden = false;
  renderMatches(state.last);
  refresh(false);
}

function swipeTop(liked) {
  const top = $('deck').querySelector('.card.top:not(.gone)');
  if (top) fling(top, liked);
}

$('create-form').addEventListener('submit', async (e) => {
  e.preventDefault();
  showError('');
  const categories = [...document.querySelectorAll('#categories input:checked')].map((i) => i.value);
  if (!categories.length) return showError('Pick at least one type of place.');
  try {
    const r = await api.createSession($('create-name').value.trim(), $('city').value, categories);
    await enterSession(r.code, r.participantId);
  } catch (err) { showError(err.message); }
});

$('join-form').addEventListener('submit', async (e) => {
  e.preventDefault();
  showError('');
  try {
    const r = await api.join($('join-code').value.trim().toUpperCase(), $('join-name').value.trim());
    await enterSession(r.code, r.participantId);
  } catch (err) { showError(err.message); }
});

$('btn-like').onclick = () => swipeTop(true);
$('btn-nope').onclick = () => swipeTop(false);
$('leave').onclick = leave;
$('share').onclick = () => share(state.code);
$('open-matches').onclick = openMatches;
$('close-matches').onclick = () => { $('matches-sheet').hidden = true; };
$('matches-sheet').onclick = (e) => { if (e.target.id === 'matches-sheet') e.target.hidden = true; };
$('match-close').onclick = () => { $('match-overlay').hidden = true; };
document.addEventListener('keydown', (e) => {
  if ($('swipe').hidden) return;
  if (e.key === 'ArrowRight') swipeTop(true);
  if (e.key === 'ArrowLeft') swipeTop(false);
});

async function init() {
  const link = (new URLSearchParams(location.search).get('code') || '').toUpperCase();
  if (link) $('join-code').value = link;
  try {
    const { cities } = await api.cities();
    $('city').innerHTML = cities.map((c) => `<option>${esc(c)}</option>`).join('');
  } catch { showError('Could not load cities. Please refresh.'); }
  const s = saved();
  if (s.code && s.pid && (!link || link === s.code)) {
    try { await enterSession(s.code, s.pid); return; } catch { localStorage.removeItem(STORE); }
  }
  if (link) $('join-name').focus();
}

init();
