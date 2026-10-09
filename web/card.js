const LABELS = { restaurant: '🍽️ Restaurant', bar: '🍸 Bar', club: '🪩 Club' };

export function esc(s) {
  return String(s).replace(/[&<>"']/g, (c) =>
    ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' })[c]);
}

// TripAdvisor-style rating bubbles.
export function bubbles(rating) {
  let html = '';
  for (let i = 1; i <= 5; i++) {
    const cls = rating >= i ? 'full' : rating >= i - 0.5 ? 'half' : '';
    html += `<i class="bubble ${cls}"></i>`;
  }
  return `<span class="bubbles" aria-label="${rating} of 5 bubbles">${html}</span>`;
}

export const price = (n) => '$'.repeat(Math.max(1, n));

export function renderCard(v) {
  const el = document.createElement('article');
  el.className = `card cat-${v.category}`;
  el.innerHTML = `
    <div class="card-art">
      <span class="card-emoji">${esc(v.emoji)}</span>
      <span class="stamp like-stamp">LIKE</span>
      <span class="stamp nope-stamp">NOPE</span>
      <span class="tag">${LABELS[v.category] || ''}</span>
    </div>
    <div class="card-info">
      <h2>${esc(v.name)}</h2>
      <div class="rating">${bubbles(v.rating)} <b>${v.rating.toFixed(1)}</b>
        <span class="muted">${v.reviewCount.toLocaleString()} reviews</span></div>
      <p class="meta">${esc(v.cuisine)} · ${price(v.priceLevel)} · 📍 ${esc(v.neighborhood)}</p>
      <p class="desc">${esc(v.description)}</p>
      <div class="highlights">${v.highlights.map((h) => `<span>${esc(h)}</span>`).join('')}</div>
    </div>`;
  return el;
}
