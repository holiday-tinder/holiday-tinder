const THRESHOLD = 100;

function setPosition(card, dx, dy) {
  card.style.transform = `translate(${dx}px, ${dy}px) rotate(${dx / 15}deg)`;
  card.style.setProperty('--like', Math.max(0, Math.min(1, dx / THRESHOLD)));
  card.style.setProperty('--nope', Math.max(0, Math.min(1, -dx / THRESHOLD)));
}

// fling animates the card off screen and reports the decision once.
export function fling(card, liked) {
  if (card.classList.contains('gone')) return;
  card.classList.add('gone');
  card.style.transition = '';
  const x = (liked ? 1 : -1) * window.innerWidth * 1.3;
  card.style.transform = `translate(${x}px, 40px) rotate(${liked ? 30 : -30}deg)`;
  card.style.setProperty(liked ? '--like' : '--nope', 1);
  setTimeout(() => card.remove(), 500);
  card._decide(liked);
}

export function makeSwipeable(card, onDecision) {
  card._decide = onDecision;
  let startX = 0, startY = 0, dx = 0, dy = 0, dragging = false;

  card.addEventListener('pointerdown', (e) => {
    dragging = true;
    startX = e.clientX;
    startY = e.clientY;
    dx = dy = 0;
    card.setPointerCapture(e.pointerId);
    card.style.transition = 'none';
  });
  card.addEventListener('pointermove', (e) => {
    if (!dragging) return;
    dx = e.clientX - startX;
    dy = e.clientY - startY;
    setPosition(card, dx, dy);
  });
  const end = () => {
    if (!dragging) return;
    dragging = false;
    card.style.transition = '';
    if (Math.abs(dx) > THRESHOLD) fling(card, dx > 0);
    else setPosition(card, 0, 0);
  };
  card.addEventListener('pointerup', end);
  card.addEventListener('pointercancel', end);
}
