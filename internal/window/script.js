// Sprite sheet: 8 cols x 9 rows, source 192x208px per frame, displayed at 120x130px
const FRAME_W = 120;
const FRAME_H = 130;

const spriteEl  = document.getElementById('sprite');
const styleEl   = document.getElementById('sprite-style');

function setAnimation(mood) {
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

  // Track mood for voice profile (e.g. "Happy 😊" → "happy")
  currentMoodForVoice = ((s.mood && s.mood.Name) || 'Happy').split(' ')[0].toLowerCase();

  var bondTier = s.bond !== undefined ? Math.min(5, Math.floor(s.bond / 20)) : 0;

  // Bond tier class on sprite (semantic, for future CSS hooks)
  for (var ti = 0; ti <= 5; ti++) spriteEl.classList.remove('bond-tier-' + ti);
  spriteEl.classList.add('bond-tier-' + bondTier);

  // Bond color filters — cool blues → warm purples → rosy pinks → golden
  var tierFilters = [
    'saturate(0.92)',
    'saturate(1.06) hue-rotate(-8deg)',
    'saturate(1.16) hue-rotate(-18deg)',
    'saturate(1.28) hue-rotate(-30deg)',
    'saturate(1.42) hue-rotate(-45deg) brightness(1.05)',
    'saturate(1.60) hue-rotate(-60deg) brightness(1.10)',
  ];
  var flowGlow = s.flowActive
    ? 'drop-shadow(0 0 6px rgba(255,200,100,0.85)) drop-shadow(0 0 14px rgba(255,140,60,0.55))'
    : '';
  spriteEl.style.filter = [tierFilters[bondTier], flowGlow].filter(Boolean).join(' ');

  // Aura — flow state or high bond tier
  var auraEl = document.getElementById('aura');
  for (var ai = 0; ai <= 5; ai++) auraEl.classList.remove('bond-aura-' + ai);
  auraEl.classList.add('bond-aura-' + bondTier);
  auraEl.classList.toggle('active', !!s.flowActive || bondTier >= 4);

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
  speakMessage(text);

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
      alpha: Math.random() * 0.9 + 0.7
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

var _starTwinkle = Array.from({ length: 5 }, function() {
  return { phase: Math.random() * Math.PI * 2, speed: Math.random() * 0.04 + 0.02 };
});
var _starPositions = [
  { x: 148, y: 14 }, { x: 138, y: 44 }, { x: 158, y: 55 },
  { x: 182, y: 54 }, { x: 192, y: 15 }
];

function drawDayNightIndicator() {
  var h = new Date().getHours();
  var isDay = h >= 5 && h < 17;

  if (isDay) {
    var sx = 170, sy = 30;
    wctx.save();
    // Outer soft bloom
    var g1 = wctx.createRadialGradient(sx, sy, 0, sx, sy, 38);
    g1.addColorStop(0,   'rgba(255, 210, 70, 0.22)');
    g1.addColorStop(0.5, 'rgba(255, 170, 40, 0.08)');
    g1.addColorStop(1,   'rgba(255, 130, 20, 0)');
    wctx.fillStyle = g1;
    wctx.beginPath(); wctx.arc(sx, sy, 38, 0, Math.PI * 2); wctx.fill();
    // Mid glow
    var g2 = wctx.createRadialGradient(sx, sy, 0, sx, sy, 22);
    g2.addColorStop(0,   'rgba(255, 230, 110, 0.45)');
    g2.addColorStop(0.7, 'rgba(255, 190, 60, 0.18)');
    g2.addColorStop(1,   'rgba(255, 160, 40, 0)');
    wctx.fillStyle = g2;
    wctx.beginPath(); wctx.arc(sx, sy, 22, 0, Math.PI * 2); wctx.fill();
    // Sun core
    var gc = wctx.createRadialGradient(sx - 2, sy - 2, 0, sx, sy, 10);
    gc.addColorStop(0,   'rgba(255, 252, 200, 1.0)');
    gc.addColorStop(0.5, 'rgba(255, 230, 100, 0.95)');
    gc.addColorStop(1,   'rgba(255, 200, 60, 0.80)');
    wctx.fillStyle = gc;
    wctx.beginPath(); wctx.arc(sx, sy, 10, 0, Math.PI * 2); wctx.fill();
    wctx.restore();

  } else {
    var mx = 170, my = 30;
    wctx.save();
    // Outer halo
    var mg1 = wctx.createRadialGradient(mx, my, 0, mx, my, 34);
    mg1.addColorStop(0,   'rgba(200, 225, 255, 0.18)');
    mg1.addColorStop(0.5, 'rgba(180, 215, 255, 0.07)');
    mg1.addColorStop(1,   'rgba(160, 205, 255, 0)');
    wctx.fillStyle = mg1;
    wctx.beginPath(); wctx.arc(mx, my, 34, 0, Math.PI * 2); wctx.fill();
    // Mid glow
    var mg2 = wctx.createRadialGradient(mx, my, 0, mx, my, 18);
    mg2.addColorStop(0,   'rgba(220, 238, 255, 0.32)');
    mg2.addColorStop(0.7, 'rgba(200, 228, 255, 0.12)');
    mg2.addColorStop(1,   'rgba(180, 218, 255, 0)');
    wctx.fillStyle = mg2;
    wctx.beginPath(); wctx.arc(mx, my, 18, 0, Math.PI * 2); wctx.fill();
    // Crescent — two arcs with nonzero winding (no destination-out, no black bg)
    var mgc = wctx.createRadialGradient(mx - 2, my - 2, 0, mx, my, 11);
    mgc.addColorStop(0,   'rgba(248, 252, 255, 0.95)');
    mgc.addColorStop(0.6, 'rgba(220, 238, 255, 0.85)');
    mgc.addColorStop(1,   'rgba(190, 222, 255, 0.70)');
    wctx.fillStyle = mgc;
    wctx.beginPath();
    wctx.arc(mx, my, 11, 0, Math.PI * 2, false);   // outer disc clockwise
    wctx.arc(mx + 6, my - 2, 8.5, 0, Math.PI * 2, true);  // inner disc counter-clockwise → crescent
    wctx.fill();
    wctx.restore();

    // Twinkling stars nearby
    for (var i = 0; i < _starTwinkle.length; i++) {
      _starTwinkle[i].phase += _starTwinkle[i].speed;
      var a = (Math.sin(_starTwinkle[i].phase) * 0.5 + 0.5) * 0.65 + 0.10;
      var sp = _starPositions[i];
      wctx.save();
      wctx.globalAlpha = a;
      wctx.fillStyle = 'rgba(220, 238, 255, 1)';
      // Cross/sparkle shape
      wctx.fillRect(sp.x - 0.5, sp.y - 2.5, 1, 5);
      wctx.fillRect(sp.x - 2.5, sp.y - 0.5, 5, 1);
      wctx.restore();
    }
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
    wctx.lineWidth = 2;
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

  drawDayNightIndicator();

  requestAnimationFrame(drawWeatherFrame);
}

drawWeatherFrame();
// ===============================

// ===============================
// SPEAK
let voiceEnabled = localStorage.getItem('dianaVoice') === 'on';
let currentMoodForVoice = 'happy';
let _dianaVoice = null;

// Diana's voice — macOS voices: Samantha, Victoria, Karen, Moira, Tessa
// Edge users: change to "Microsoft Ana Online (Natural)"
const DIANA_VOICE_NAME = "Superstar";

// Mood-driven speech profiles
const MOOD_VOICE = {
  happy:    { pitch: 1.25, rate: 1.00 },
  elated:   { pitch: 1.35, rate: 1.08 },
  starving: { pitch: 1.05, rate: 0.92 },
  hungry:   { pitch: 1.05, rate: 0.92 },
  lonely:   { pitch: 1.00, rate: 0.88 },
  idle:     { pitch: 1.10, rate: 0.95 },
};

function cleanTextForSpeech(text) {
  // If it's HTML (gift message), extract just the readable text
  let clean = text;
  if (clean.trimStart().startsWith('<')) {
    clean = clean.replace(/<style\b[^>]*>([\s\S]*?)<\/style>/gi, "");
    const tmp = document.createElement('div');
    tmp.innerHTML = clean;
    clean = tmp.textContent || '';
  }
  // Strip emojis and pictographic symbols so Diana doesn't say "sparkles"
  clean = clean
    .replace(/\p{Extended_Pictographic}/gu, "")
    .trim();

  return clean;
}

function speakMessage(text) {
  if (!voiceEnabled) return;

  const clean = cleanTextForSpeech(text);
  if (!clean) return;

  // Mood profile + tiny randomness so she doesn't sound robotic
  const profile = MOOD_VOICE[currentMoodForVoice] || MOOD_VOICE.happy;
  const pitch = profile.pitch + (Math.random() - 0.5) * 0.10;
  const rate  = profile.rate  + (Math.random() - 0.5) * 0.06;

  window.speechSynthesis.cancel(); // stop overlapping speech

  const utterance = new SpeechSynthesisUtterance(clean);
  utterance.pitch = pitch
  utterance.rate  = rate;
  if (_dianaVoice) utterance.voice = _dianaVoice;

  window.speechSynthesis.speak(utterance);
}

function updateVoiceButton() {
  const btn = document.getElementById('voice-toggle');
  if (!btn) return;
  btn.textContent = voiceEnabled ? '🔊' : '🔇';
  btn.classList.toggle('active', voiceEnabled);
  btn.title = voiceEnabled ? "Diana's voice is on (click to mute)" : "Diana's voice is off (click to enable)";
}

function toggleVoice() {
  voiceEnabled = !voiceEnabled;
  localStorage.setItem('dianaVoice', voiceEnabled ? 'on' : 'off');
  updateVoiceButton();
  if (!voiceEnabled) window.speechSynthesis.cancel();
}

// Wire up the toggle button + initial voice selection
document.getElementById('voice-toggle').addEventListener('click', function(e) {
  e.stopPropagation();
  toggleVoice();
});
updateVoiceButton();
// ===============================