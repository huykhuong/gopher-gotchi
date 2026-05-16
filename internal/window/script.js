// Sprite sheet: 8 cols x 9 rows, source 192x208px per frame, displayed at 120x130px
const FRAME_W = 120;
const FRAME_H = 130;

const spriteEl  = document.getElementById('sprite');
const styleEl   = document.getElementById('sprite-style');
let currentAnim = null;

function setAnimation(mood) {
  if (currentAnim === mood) return;
  currentAnim = mood;
  console.log(mood);
  styleEl.textContent = `@keyframes sp { from { background-position-x: 0px } to { background-position-x: ${-(mood.Frames * FRAME_W)}px } }`;
  spriteEl.style.backgroundPositionY = -(mood.Row * FRAME_H) + 'px';
  spriteEl.style.animation = `sp ${mood.Duration}ms steps(${mood.Frames}) infinite`;
}

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

// State updates
let bubbleTimer = null;

function updateState(s) {
  document.getElementById('level').textContent = 'LVL ' + (s.level || 1);
  document.getElementById('mood').textContent  = s.mood.Name || '';

  // Flow aura — glows when on a save streak
  document.getElementById('aura').classList.toggle('active', !!s.flowActive);
  document.getElementById('sprite').classList.toggle('flow', !!s.flowActive);
  const hunger = s.hunger || 0;

  setAnimation(s.mood);

  const h = Math.min(100, Math.max(0, hunger));
  document.getElementById('hunger-bar').style.width = h + '%';
  document.getElementById('hunger-val').textContent = h + '%';

  const c = Math.min(100, Math.max(0, s.cpuLoad || 0));
  document.getElementById('cpu-bar').style.width = c + '%';
  document.getElementById('cpu-val').textContent = c + '%';

  if (s.weatherKnowledge) {
    applyWeather(s.weatherKnowledge.Condition, s.weatherKnowledge.IsDay);
  }

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

// ===============================
// DRAG LOGIC
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
// ===============================


// ===============================
// WEATHER CANVAS
// ===============================
const weatherCanvas = document.getElementById('weather-canvas');
const wctx = weatherCanvas.getContext('2d');
weatherCanvas.width  = 200;
weatherCanvas.height = 340;

let weatherKey = '';
let drops = [], stars = [], clouds = [], motes = [];
let lightningTimer = 0, lightningAlpha = 0;
let wframe = 0;

function initRain(count) {
  drops = Array.from({ length: count }, function() {
    return {
      x: Math.random() * 220 - 10,
      y: Math.random() * 340,
      len: Math.random() * 9 + 5,
      speed: Math.random() * 4 + 5,
      alpha: Math.random() * 0.3 + 0.12
    };
  });
}

function initStars() {
  stars = Array.from({ length: 26 }, function() {
    return {
      x: Math.random() * 200,
      y: Math.random() * 210,
      r: Math.random() * 0.9 + 0.3,
      phase: Math.random() * Math.PI * 2,
      speed: Math.random() * 0.022 + 0.007
    };
  });
}

function initClouds() {
  clouds = Array.from({ length: 3 }, function(_, i) {
    return {
      x: Math.random() * 250 - 60,
      y: 30 + i * 55 + Math.random() * 25,
      speed: Math.random() * 0.12 + 0.04,
      scale: Math.random() * 0.45 + 0.55,
      alpha: Math.random() * 0.13 + 0.06
    };
  });
}

function initMotes() {
  motes = Array.from({ length: 14 }, function() {
    return {
      x: Math.random() * 200,
      y: Math.random() * 300 + 20,
      r: Math.random() * 1.4 + 0.4,
      phase: Math.random() * Math.PI * 2,
      speed: Math.random() * 0.013 + 0.005,
      vx: (Math.random() - 0.5) * 0.25,
      vy: -(Math.random() * 0.25 + 0.08),
      maxAlpha: Math.random() * 0.35 + 0.08
    };
  });
}

function applyWeather(condition, isDay) {
  var key = condition + '|' + isDay;
  if (key === weatherKey) return;
  weatherKey = key;
  drops = []; stars = []; clouds = []; motes = [];

  if (condition === 'Raining') {
    initRain(45);
  } else if (condition === 'Stormy') {
    initRain(80);
    lightningTimer = Math.random() * 200 + 120;
  } else if (condition === 'Cloudy') {
    initClouds();
  } else if (condition === 'Sunny / Clear') {
    if (isDay) initMotes();
    else        initStars();
  } else {
    initClouds(); // foggy / unspecified
  }
}

function drawCloud(x, y, scale, alpha) {
  var puffs = [[0,0,28],[-22,8,20],[22,8,20],[-10,10,22],[12,10,18]];
  wctx.save();
  wctx.globalAlpha = alpha;
  wctx.fillStyle = '#c8d8ea';
  for (var i = 0; i < puffs.length; i++) {
    var px = puffs[i][0], py = puffs[i][1], r = puffs[i][2];
    wctx.beginPath();
    wctx.arc(x + px * scale, y + py * scale, r * scale, 0, Math.PI * 2);
    wctx.fill();
  }
  wctx.restore();
}

function getTimeOfDay() {
  var h = new Date().getHours();
  if (h >= 5  && h < 10) return 'morning';
  if (h >= 10 && h < 17) return 'day';
  if (h >= 17 && h < 21) return 'evening';
  return 'night';
}

function drawGlow(cx, cy, radius, colorStops) {
  var g = wctx.createRadialGradient(cx, cy, 0, cx, cy, radius);
  for (var i = 0; i < colorStops.length; i++) {
    g.addColorStop(colorStops[i][0], colorStops[i][1]);
  }
  wctx.save();
  wctx.fillStyle = g;
  wctx.beginPath();
  wctx.arc(cx, cy, radius, 0, Math.PI * 2);
  wctx.fill();
  wctx.restore();
}

function drawTimeAmbient(tod) {
  if (tod === 'morning') {
    // Warm golden sunrise spilling from top-right corner
    drawGlow(200, 0, 200, [
      [0,    'rgba(255, 215, 80,  0.65)'],
      [0.35, 'rgba(255, 170, 55,  0.30)'],
      [0.65, 'rgba(255, 120, 40,  0.10)'],
      [1,    'rgba(255,  80, 20,  0)']
    ]);
    // Soft secondary bloom for warmth bleed
    drawGlow(200, 0, 120, [
      [0,   'rgba(255, 240, 160, 0.30)'],
      [1,   'rgba(255, 200,  80, 0)']
    ]);

  } else if (tod === 'day') {
    // Still sunlight from top-right, cooler and brighter mid-day
    drawGlow(200, 0, 190, [
      [0,    'rgba(255, 235, 130, 0.50)'],
      [0.40, 'rgba(255, 205,  80, 0.20)'],
      [0.70, 'rgba(255, 180,  50, 0.07)'],
      [1,    'rgba(255, 160,  30, 0)']
    ]);

  } else if (tod === 'evening') {
    // Rich orange-to-violet bloom from bottom-right
    drawGlow(190, 340, 230, [
      [0,    'rgba(255, 110, 50,  0.55)'],
      [0.38, 'rgba(210,  70, 120, 0.28)'],
      [1,    'rgba(130,  40, 180, 0)']
    ]);

  } else {
    // Night — soft moonlight pooling from top centre, top-half only
    drawGlow(100, 0, 185, [
      [0,    'rgba(210, 230, 255, 0.52)'],
      [0.40, 'rgba(180, 210, 255, 0.22)'],
      [0.70, 'rgba(150, 185, 255, 0.07)'],
      [1,    'rgba(120, 160, 255, 0)']
    ]);
    // Narrow bright core — the moon's direct beam
    drawGlow(100, 0, 90, [
      [0,   'rgba(235, 245, 255, 0.35)'],
      [0.5, 'rgba(210, 230, 255, 0.12)'],
      [1,   'rgba(190, 215, 255, 0)']
    ]);
  }
}

function drawWeatherFrame() {
  wctx.clearRect(0, 0, 200, 340);
  wframe++;

  var condition = weatherKey.split('|')[0];
  var isDay     = weatherKey.split('|')[1] === 'true';

  drawTimeAmbient(getTimeOfDay());

  if (condition === 'Sunny / Clear') {
    if (isDay) {
      // Ambient warm glow top-right
      var grd = wctx.createRadialGradient(162, 52, 8, 162, 52, 72);
      grd.addColorStop(0,   'rgba(255,215,90,0.18)');
      grd.addColorStop(0.5, 'rgba(255,175,55,0.08)');
      grd.addColorStop(1,   'rgba(255,140,30,0)');
      wctx.fillStyle = grd;
      wctx.beginPath();
      wctx.arc(162, 52, 72, 0, Math.PI * 2);
      wctx.fill();

      // Sun core
      wctx.save();
      wctx.globalAlpha = 0.52;
      var sc = wctx.createRadialGradient(162, 52, 0, 162, 52, 14);
      sc.addColorStop(0,   '#fffce0');
      sc.addColorStop(0.6, '#ffd055');
      sc.addColorStop(1,   'rgba(255,175,55,0)');
      wctx.fillStyle = sc;
      wctx.beginPath();
      wctx.arc(162, 52, 14, 0, Math.PI * 2);
      wctx.fill();
      wctx.restore();

      // Rotating rays
      wctx.save();
      wctx.globalAlpha = 0.10;
      wctx.strokeStyle = '#ffd055';
      wctx.lineWidth = 1.5;
      for (var r = 0; r < 8; r++) {
        var angle = (r / 8) * Math.PI * 2 + wframe * 0.003;
        wctx.beginPath();
        wctx.moveTo(162 + Math.cos(angle) * 18, 52 + Math.sin(angle) * 18);
        wctx.lineTo(162 + Math.cos(angle) * 34, 52 + Math.sin(angle) * 34);
        wctx.stroke();
      }
      wctx.restore();

      // Floating light motes
      for (var mi = 0; mi < motes.length; mi++) {
        var m = motes[mi];
        m.phase += m.speed;
        m.x += m.vx;
        m.y += m.vy;
        if (m.y < -4) { m.y = 342; m.x = Math.random() * 200; }
        wctx.save();
        wctx.globalAlpha = (Math.sin(m.phase) * 0.5 + 0.5) * m.maxAlpha;
        wctx.fillStyle = '#ffe070';
        wctx.beginPath();
        wctx.arc(m.x, m.y, m.r, 0, Math.PI * 2);
        wctx.fill();
        wctx.restore();
      }

    } else {
      // Moon
      var mx = 156, my = 48;
      var mg = wctx.createRadialGradient(mx, my, 6, mx, my, 48);
      mg.addColorStop(0, 'rgba(190,210,255,0.10)');
      mg.addColorStop(1, 'rgba(170,195,255,0)');
      wctx.fillStyle = mg;
      wctx.beginPath();
      wctx.arc(mx, my, 48, 0, Math.PI * 2);
      wctx.fill();

      // Moon crescent
      wctx.save();
      wctx.globalAlpha = 0.72;
      wctx.fillStyle = '#dce8ff';
      wctx.beginPath();
      wctx.arc(mx, my, 11, 0, Math.PI * 2);
      wctx.fill();
      wctx.globalCompositeOperation = 'destination-out';
      wctx.fillStyle = 'rgba(0,0,0,0.88)';
      wctx.beginPath();
      wctx.arc(mx + 6, my - 3, 9, 0, Math.PI * 2);
      wctx.fill();
      wctx.restore();

      // Twinkling stars
      for (var si = 0; si < stars.length; si++) {
        var st = stars[si];
        st.phase += st.speed;
        wctx.save();
        wctx.globalAlpha = (Math.sin(st.phase) * 0.35 + 0.65) * 0.65;
        wctx.fillStyle = '#c8dcff';
        wctx.beginPath();
        wctx.arc(st.x, st.y, st.r, 0, Math.PI * 2);
        wctx.fill();
        wctx.restore();
      }
    }

  } else if (condition === 'Cloudy') {
    for (var ci = 0; ci < clouds.length; ci++) {
      var cl = clouds[ci];
      cl.x += cl.speed;
      if (cl.x > 280) cl.x = -110;
      drawCloud(cl.x, cl.y, cl.scale, cl.alpha);
    }

  } else if (condition === 'Raining' || condition === 'Stormy') {
    wctx.lineWidth = 0.8;
    for (var di = 0; di < drops.length; di++) {
      var d = drops[di];
      d.y += d.speed;
      d.x -= d.speed * 0.18;
      if (d.y > 342) { d.y = -8; d.x = Math.random() * 220 - 10; }
      wctx.save();
      wctx.globalAlpha = d.alpha;
      wctx.strokeStyle = 'rgba(175,210,255,0.9)';
      wctx.beginPath();
      wctx.moveTo(d.x, d.y);
      wctx.lineTo(d.x - d.len * 0.18, d.y + d.len);
      wctx.stroke();
      wctx.restore();
    }

    if (condition === 'Stormy') {
      if (lightningAlpha > 0) {
        lightningAlpha = Math.max(0, lightningAlpha - 0.07);
        wctx.fillStyle = 'rgba(215,235,255,' + lightningAlpha + ')';
        wctx.fillRect(0, 0, 200, 340);
      }
      if (--lightningTimer <= 0) {
        lightningAlpha = 0.32;
        lightningTimer = Math.random() * 280 + 140;
      }
    }

  } else if (condition !== '') {
    // Foggy / Unspecified — slow drifting wisps
    if (clouds.length === 0) initClouds();
    for (var fi = 0; fi < clouds.length; fi++) {
      var fc = clouds[fi];
      fc.x += fc.speed * 0.45;
      if (fc.x > 280) fc.x = -110;
      drawCloud(fc.x, fc.y, fc.scale * 1.6, fc.alpha * 0.6);
    }
  }

  requestAnimationFrame(drawWeatherFrame);
}

drawWeatherFrame();

// {
//   "name": "Diana",
//   "species": "diana",
//   "level": 31,
//   "experience": 503,
//   "hunger": 0,
//   "mood": "happy",
//   "last_eaten": "2026-05-14T20:38:02.630367+07:00"
// }