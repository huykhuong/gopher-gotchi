// Sprite sheet: 8 cols x 9 rows, source 192x208px per frame, displayed at 120x130px
const FRAME_W = 120;
const FRAME_H = 130;

const ANIMS = {
  idle:   { row: 0, frames: 6, duration: 1100 },
  right:  { row: 1, frames: 8, duration: 1060 },
  left:   { row: 2, frames: 8, duration: 1060 },
  wave:   { row: 3, frames: 4, duration: 1000  },
  jump:   { row: 4, frames: 5, duration: 840  },
  failed: { row: 5, frames: 8, duration: 1220 },
  wait:   { row: 6, frames: 6, duration: 1010 },
  run:    { row: 7, frames: 6, duration: 820  },
  review: { row: 8, frames: 6, duration: 1030 },
};

const spriteEl  = document.getElementById('sprite');
const styleEl   = document.getElementById('sprite-style');
let currentAnim = null;

function setAnimation(name) {
  if (currentAnim === name) return;
  currentAnim = name;
  const a = ANIMS[name];
  if (!a) return;
  styleEl.textContent = `@keyframes sp { from { background-position-x: 0px } to { background-position-x: ${-(a.frames * FRAME_W)}px } }`;
  spriteEl.style.backgroundPositionY = -(a.row * FRAME_H) + 'px';
  spriteEl.style.animation = `sp ${a.duration}ms steps(${a.frames}) infinite`;
}

setAnimation('idle');

// Radial action menu
const actionBtn  = document.getElementById('action-btn');
const radialMenu = document.getElementById('radial-menu');

actionBtn.addEventListener('click', function(e) {
  e.stopPropagation();
  radialMenu.classList.toggle('open');
});

document.addEventListener('click', function(e) {
  if (!e.target.closest('.radial-item') && e.target !== actionBtn) {
    radialMenu.classList.remove('open');
  }
});

document.querySelectorAll('.radial-item').forEach(function(item) {
  item.addEventListener('click', function(e) {
    e.stopPropagation();
    const action = item.dataset.action;
    window.petAction(action);
    radialMenu.classList.remove('open');
  });
});

// Drag — ignore clicks on any interactive element
let dragging = false, sx = 0, sy = 0;
document.addEventListener('mousedown', function(e) {
  if (e.target.closest('.quit-btn, .action-btn, .radial-item, .mission-panel')) return;
  dragging = true; sx = e.screenX; sy = e.screenY;
  e.preventDefault();
});
document.addEventListener('mousemove', function(e) {
  if (!dragging) return;
  const dx = e.screenX - sx, dy = e.screenY - sy;
  sx = e.screenX; sy = e.screenY;
  window.moveWindowBy(dx, dy).catch(function(){});
});
document.addEventListener('mouseup', function() { dragging = false; });

// State updates
let bubbleTimer = null;

function updateState(s) {
  document.getElementById('level').textContent = 'LVL ' + (s.level || 1);
  document.getElementById('mood').textContent  = s.mood || '';

  // Flow aura — glows when on a save streak
  document.getElementById('aura').classList.toggle('active', !!s.flowActive);
  document.getElementById('sprite').classList.toggle('flow', !!s.flowActive);

  const mood   = (s.mood || '').toLowerCase();
  const hunger = s.hunger || 0;
  
  if      (hunger > 70)             setAnimation('failed');
  else if (mood.includes('happy'))  {
    var random = Math.floor(Math.random() * 3);
    if (random === 0) setAnimation('wave');
    else if (random === 1) setAnimation('right');
    else setAnimation('run');
  }
  else if (mood.includes('lonely')) setAnimation('wait');
  else if (mood.includes('grumpy')) setAnimation('failed');
  else if (mood.includes('elated')) setAnimation('jump');
  else                              setAnimation('idle');

  const h = Math.min(100, Math.max(0, hunger));
  document.getElementById('hunger-bar').style.width = h + '%';
  document.getElementById('hunger-val').textContent = h + '%';

  const c = Math.min(100, Math.max(0, s.cpuLoad || 0));
  document.getElementById('cpu-bar').style.width = c + '%';
  document.getElementById('cpu-val').textContent = c + '%';

  if (s.message) {
    // Detect a mission list by the task line pattern: "> [ID] title (status)"
    if (/^>\s*\[[^\]]+\].*\([^)]+\)\s*$/m.test(s.message)) {
      showMission(s.message);
    } else {
      showBubble(s.message);
    }
  }
}

// Mission panel — rendered when an objectives message arrives
const missionPanel = document.getElementById('mission-panel');
const missionList  = document.getElementById('mission-list');
document.getElementById('mission-close').addEventListener('click', function() {
  missionPanel.classList.remove('show');
});

function escapeHTML(s) {
  const d = document.createElement('div');
  d.textContent = s == null ? '' : String(s);
  return d.innerHTML;
}

function statusClass(name) {
  return (name || '').toLowerCase().trim().replace(/\s+/g, '-');
}

function showMission(text) {
  const items = [];
  const lineRegex = /^>\s*\[([^\]]+)\]\s*(.+?)\s*\(([^)]+)\)\s*$/;
  text.split('\n').forEach(function(line) {
    const m = line.match(lineRegex);
    if (m) items.push({ id: m[1], title: m[2], status: m[3] });
  });

  if (items.length === 0) {
    missionList.innerHTML = '<div class="mission-empty">✨ All clear, Huy.<br>No active objectives in this sector.</div>';
  } else {
    missionList.innerHTML = items.map(function(it) {
      return '<div class="mission-item">' +
               '<div class="mission-id">' + escapeHTML(it.id) + '</div>' +
               '<div class="mission-text">' + escapeHTML(it.title) + '</div>' +
               '<div class="mission-status ' + statusClass(it.status) + '">' + escapeHTML(it.status) + '</div>' +
             '</div>';
    }).join('');
  }

  missionPanel.classList.add('show');
}

function showBubble(text) {
  const el = document.getElementById('bubble');
  const isGift = text.trimStart().startsWith('<');
  if (isGift) {
    el.innerHTML = text;
  } else {
    el.innerText = text;
  }
  el.classList.add('show');
  if (bubbleTimer) clearTimeout(bubbleTimer);
  bubbleTimer = setTimeout(function() { el.classList.remove('show'); }, 5000);
}