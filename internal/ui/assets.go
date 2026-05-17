package ui

type Gift struct {
	Name   string
	Rarity string
	Art    string
}

var BondPointsBasedOnRarity = map[string]int{
	"Common":    1,
	"Uncommon":  2,
	"Rare":      3,
	"Epic":      4,
	"Legendary": 5,
}

var DigitalTreasury = []Gift{
	{
		Name:   "Holographic Rose",
		Rarity: "Rare",
		Art: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 80 72" width="80" height="72">
<defs>
  <radialGradient id="rg" cx="50%" cy="45%" r="50%">
    <stop offset="0%" stop-color="#ff4da6" stop-opacity=".55"/>
    <stop offset="100%" stop-color="#ff4da6" stop-opacity="0"/>
  </radialGradient>
  <radialGradient id="cg" cx="35%" cy="30%" r="65%">
    <stop offset="0%" stop-color="#ff80c0"/>
    <stop offset="100%" stop-color="#b5005b"/>
  </radialGradient>
  <filter id="rf">
    <feGaussianBlur stdDeviation="2.5" result="b"/>
    <feMerge><feMergeNode in="b"/><feMergeNode in="SourceGraphic"/></feMerge>
  </filter>
  <style>
    @keyframes rosePulse{0%,100%{transform:scale(1);opacity:.45}50%{transform:scale(1.15);opacity:.7}}
    @keyframes roseFloat{0%,100%{transform:translateY(0)}50%{transform:translateY(-2px)}}
    .rh{transform-origin:40px 35px;animation:rosePulse 2.2s ease-in-out infinite}
    .rb{transform-origin:40px 35px;animation:roseFloat 3s ease-in-out infinite}
  </style>
</defs>
<ellipse class="rh" cx="40" cy="35" rx="32" ry="26" fill="url(#rg)"/>
<g class="rb">
  <line x1="40" y1="56" x2="40" y2="69" stroke="#4caf50" stroke-width="2.5" stroke-linecap="round"/>
  <path d="M40 64 Q35 60 33 55" stroke="#4caf50" stroke-width="1.5" fill="none" stroke-linecap="round"/>
  <ellipse cx="40" cy="25" rx="7" ry="12" fill="url(#cg)" opacity=".9" transform="rotate(0 40 37)"/>
  <ellipse cx="40" cy="25" rx="7" ry="12" fill="url(#cg)" opacity=".88" transform="rotate(52 40 37)"/>
  <ellipse cx="40" cy="25" rx="7" ry="12" fill="url(#cg)" opacity=".85" transform="rotate(104 40 37)"/>
  <ellipse cx="40" cy="25" rx="7" ry="12" fill="url(#cg)" opacity=".88" transform="rotate(156 40 37)"/>
  <ellipse cx="40" cy="25" rx="7" ry="12" fill="url(#cg)" opacity=".85" transform="rotate(208 40 37)"/>
  <ellipse cx="40" cy="25" rx="7" ry="12" fill="url(#cg)" opacity=".88" transform="rotate(260 40 37)"/>
  <ellipse cx="40" cy="25" rx="7" ry="12" fill="url(#cg)" opacity=".85" transform="rotate(312 40 37)"/>
  <circle cx="40" cy="37" r="9" fill="#e91e63" filter="url(#rf)"/>
  <circle cx="40" cy="37" r="6" fill="#ff4081"/>
  <circle cx="37" cy="34" r="2.5" fill="#fff" opacity=".55"/>
</g>
</svg>`,
	},
	{
		Name:   "Binary Star",
		Rarity: "Uncommon",
		Art: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 80 72" width="80" height="72">
<defs>
  <radialGradient id="sg" cx="50%" cy="50%" r="50%">
    <stop offset="0%" stop-color="#ffe066" stop-opacity=".6"/>
    <stop offset="100%" stop-color="#ff9800" stop-opacity="0"/>
  </radialGradient>
  <radialGradient id="starFill" cx="40%" cy="30%" r="65%">
    <stop offset="0%" stop-color="#fff176"/>
    <stop offset="100%" stop-color="#f57f17"/>
  </radialGradient>
  <filter id="sf">
    <feGaussianBlur stdDeviation="3" result="b"/>
    <feMerge><feMergeNode in="b"/><feMergeNode in="SourceGraphic"/></feMerge>
  </filter>
  <style>
    @keyframes starSpin{from{transform:rotate(0deg)}to{transform:rotate(360deg)}}
    @keyframes starPulse{0%,100%{opacity:.5;transform:scale(1)}50%{opacity:.85;transform:scale(1.2)}}
    .sr{transform-origin:40px 36px;animation:starSpin 6s linear infinite}
    .sh{transform-origin:40px 36px;animation:starPulse 2s ease-in-out infinite}
  </style>
</defs>
<ellipse class="sh" cx="40" cy="36" rx="30" ry="26" fill="url(#sg)"/>
<g class="sr">
  <line x1="40" y1="10" x2="40" y2="14" stroke="#ffe066" stroke-width="1.5" opacity=".7"/>
  <line x1="40" y1="58" x2="40" y2="62" stroke="#ffe066" stroke-width="1.5" opacity=".7"/>
  <line x1="14" y1="36" x2="18" y2="36" stroke="#ffe066" stroke-width="1.5" opacity=".7"/>
  <line x1="62" y1="36" x2="66" y2="36" stroke="#ffe066" stroke-width="1.5" opacity=".7"/>
  <line x1="21" y1="17" x2="24" y2="20" stroke="#ffe066" stroke-width="1.5" opacity=".5"/>
  <line x1="56" y1="52" x2="59" y2="55" stroke="#ffe066" stroke-width="1.5" opacity=".5"/>
  <line x1="56" y1="17" x2="59" y2="20" stroke="#ffe066" stroke-width="1.5" opacity=".5" transform="scale(-1,1) translate(-80,0)"/>
</g>
<polygon points="40,14 43.8,25.5 56,25.5 46.1,32.5 49.9,44 40,37 30.1,44 33.9,32.5 24,25.5 36.2,25.5"
  fill="url(#starFill)" filter="url(#sf)"/>
<polygon points="40,18 43,27 52,27 45,32 47.5,41 40,36 32.5,41 35,32 28,27 37,27"
  fill="#fff9c4" opacity=".4"/>
<circle cx="40" cy="30" r="3" fill="#fff" opacity=".7"/>
</svg>`,
	},
	{
		Name:   "Lucky Fragment",
		Rarity: "Epic",
		Art: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 80 72" width="80" height="72">
<defs>
  <radialGradient id="gg" cx="50%" cy="50%" r="50%">
    <stop offset="0%" stop-color="#00e5ff" stop-opacity=".5"/>
    <stop offset="100%" stop-color="#7c4dff" stop-opacity="0"/>
  </radialGradient>
  <linearGradient id="gemTop" x1="30%" y1="0%" x2="70%" y2="100%">
    <stop offset="0%" stop-color="#e0f7ff"/>
    <stop offset="40%" stop-color="#00bcd4"/>
    <stop offset="100%" stop-color="#4a148c"/>
  </linearGradient>
  <linearGradient id="gemLeft" x1="0%" y1="0%" x2="100%" y2="100%">
    <stop offset="0%" stop-color="#26c6da"/>
    <stop offset="100%" stop-color="#6a1b9a"/>
  </linearGradient>
  <linearGradient id="shimmer" x1="0%" y1="0%" x2="100%" y2="0%">
    <stop offset="0%" stop-color="#fff" stop-opacity="0"/>
    <stop offset="50%" stop-color="#fff" stop-opacity=".45"/>
    <stop offset="100%" stop-color="#fff" stop-opacity="0"/>
  </linearGradient>
  <clipPath id="gemClip">
    <polygon points="40,10 58,24 58,48 40,62 22,48 22,24"/>
  </clipPath>
  <filter id="gf">
    <feGaussianBlur stdDeviation="3" result="b"/>
    <feMerge><feMergeNode in="b"/><feMergeNode in="SourceGraphic"/></feMerge>
  </filter>
  <style>
    @keyframes gemPulse{0%,100%{opacity:.45;transform:scale(1)}50%{opacity:.75;transform:scale(1.18)}}
    @keyframes shimmerMove{0%{transform:translateX(-80px)}100%{transform:translateX(80px)}}
    .glow{transform-origin:40px 36px;animation:gemPulse 2.4s ease-in-out infinite}
    .shim{animation:shimmerMove 2.5s ease-in-out infinite}
  </style>
</defs>
<ellipse class="glow" cx="40" cy="36" rx="30" ry="26" fill="url(#gg)"/>
<polygon points="40,10 58,24 58,48 40,62 22,48 22,24" fill="url(#gemTop)" filter="url(#gf)"/>
<polygon points="40,10 22,24 40,36" fill="url(#gemLeft)" opacity=".9"/>
<polygon points="40,10 58,24 40,36" fill="#e0f7ff" opacity=".35"/>
<polygon points="22,24 22,48 40,36" fill="#006064" opacity=".6"/>
<polygon points="58,24 58,48 40,36" fill="#4a148c" opacity=".5"/>
<polygon points="22,48 40,62 40,36" fill="#26c6da" opacity=".55"/>
<polygon points="58,48 40,62 40,36" fill="#7c4dff" opacity=".5"/>
<rect x="0" y="10" width="30" height="52" fill="url(#shimmer)" clip-path="url(#gemClip)" class="shim"/>
<circle cx="34" cy="22" r="2.5" fill="#fff" opacity=".8"/>
<circle cx="29" cy="28" r="1.2" fill="#fff" opacity=".5"/>
</svg>`,
	},
	{
		Name:   "Pixel Heart",
		Rarity: "Rare",
		Art: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 80 72" width="80" height="72">
<defs>
  <radialGradient id="hg2" cx="50%" cy="50%" r="50%">
    <stop offset="0%" stop-color="#ff4d6d" stop-opacity=".55"/>
    <stop offset="100%" stop-color="#ff4d6d" stop-opacity="0"/>
  </radialGradient>
  <linearGradient id="heartFill" x1="30%" y1="0%" x2="70%" y2="100%">
    <stop offset="0%" stop-color="#ff8fa3"/>
    <stop offset="100%" stop-color="#c9184a"/>
  </linearGradient>
  <filter id="hf2">
    <feGaussianBlur stdDeviation="2.5" result="b"/>
    <feMerge><feMergeNode in="b"/><feMergeNode in="SourceGraphic"/></feMerge>
  </filter>
  <style>
    @keyframes heartBeat{0%,100%{transform:scale(1)}14%{transform:scale(1.15)}28%{transform:scale(1)}42%{transform:scale(1.1)}70%{transform:scale(1)}}
    @keyframes heartGlow{0%,100%{opacity:.4;transform:scale(1)}50%{opacity:.7;transform:scale(1.2)}}
    .hb{transform-origin:40px 38px;animation:heartBeat 1.4s ease-in-out infinite}
    .hg2{transform-origin:40px 38px;animation:heartGlow 1.4s ease-in-out infinite}
  </style>
</defs>
<ellipse class="hg2" cx="40" cy="38" rx="30" ry="25" fill="url(#hg2)"/>
<g class="hb" filter="url(#hf2)">
  <path d="M40 58 C40 58 18 44 18 30 C18 22 24 16 32 16 C36 16 40 19 40 19 C40 19 44 16 48 16 C56 16 62 22 62 30 C62 44 40 58 40 58Z" fill="url(#heartFill)"/>
  <path d="M28 22 C24 22 22 26 22 30 C22 36 28 42 34 48" stroke="#ff8fa3" stroke-width="2" fill="none" opacity=".5" stroke-linecap="round"/>
  <circle cx="34" cy="24" r="3" fill="#fff" opacity=".45"/>
</g>
</svg>`,
	},
	{
		Name:   "Ancient Coin",
		Rarity: "Uncommon",
		Art: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 80 72" width="80" height="72">
<defs>
  <radialGradient id="coinGlow" cx="50%" cy="50%" r="50%">
    <stop offset="0%" stop-color="#ffd54f" stop-opacity=".6"/>
    <stop offset="100%" stop-color="#ff8f00" stop-opacity="0"/>
  </radialGradient>
  <linearGradient id="coinFace" x1="25%" y1="15%" x2="75%" y2="85%">
    <stop offset="0%" stop-color="#fff8e1"/>
    <stop offset="40%" stop-color="#ffd54f"/>
    <stop offset="100%" stop-color="#e65100"/>
  </linearGradient>
  <filter id="coinF">
    <feGaussianBlur stdDeviation="2" result="b"/>
    <feMerge><feMergeNode in="b"/><feMergeNode in="SourceGraphic"/></feMerge>
  </filter>
  <style>
    @keyframes coinSpin{0%,100%{transform:scaleX(1)}25%{transform:scaleX(0.15)}50%{transform:scaleX(1)}75%{transform:scaleX(0.15)}}
    @keyframes coinPulse{0%,100%{opacity:.45;transform:scale(1)}50%{opacity:.75;transform:scale(1.15)}}
    .cs{transform-origin:40px 36px;animation:coinSpin 2.4s ease-in-out infinite}
    .cp{transform-origin:40px 36px;animation:coinPulse 2.4s ease-in-out infinite}
  </style>
</defs>
<ellipse class="cp" cx="40" cy="36" rx="30" ry="24" fill="url(#coinGlow)"/>
<g class="cs" filter="url(#coinF)">
  <ellipse cx="40" cy="36" rx="22" ry="22" fill="#bf6c00"/>
  <ellipse cx="40" cy="36" rx="20" ry="20" fill="url(#coinFace)"/>
  <circle cx="40" cy="36" r="12" fill="none" stroke="#bf6c00" stroke-width="1.5" opacity=".5"/>
  <text x="40" y="41" text-anchor="middle" font-size="14" font-weight="bold" fill="#7f3e00" opacity=".7" font-family="serif">✦</text>
  <circle cx="34" cy="30" r="3" fill="#fff" opacity=".35"/>
</g>
</svg>`,
	},
	{
		Name:   "Lightning Shard",
		Rarity: "Epic",
		Art: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 80 72" width="80" height="72">
<defs>
  <radialGradient id="lglow" cx="50%" cy="50%" r="50%">
    <stop offset="0%" stop-color="#e8ff00" stop-opacity=".55"/>
    <stop offset="100%" stop-color="#00b0ff" stop-opacity="0"/>
  </radialGradient>
  <linearGradient id="boltFill" x1="30%" y1="0%" x2="70%" y2="100%">
    <stop offset="0%" stop-color="#ffffff"/>
    <stop offset="30%" stop-color="#ffff00"/>
    <stop offset="100%" stop-color="#00b0ff"/>
  </linearGradient>
  <filter id="lf">
    <feGaussianBlur stdDeviation="3" result="b"/>
    <feMerge><feMergeNode in="b"/><feMergeNode in="SourceGraphic"/></feMerge>
  </filter>
  <style>
    @keyframes boltFlicker{0%,100%{opacity:1;transform:scale(1)}20%{opacity:.7;transform:scale(.97)}40%{opacity:1;transform:scale(1.02)}60%{opacity:.8;transform:scale(.98)}80%{opacity:1}}
    @keyframes boltGlow{0%,100%{opacity:.4;transform:scale(1)}50%{opacity:.8;transform:scale(1.25)}}
    .bf{transform-origin:40px 36px;animation:boltFlicker 1.6s ease-in-out infinite}
    .bg2{transform-origin:40px 36px;animation:boltGlow 1.6s ease-in-out infinite}
  </style>
</defs>
<ellipse class="bg2" cx="40" cy="36" rx="28" ry="24" fill="url(#lglow)"/>
<g class="bf" filter="url(#lf)">
  <polygon points="46,10 30,38 42,38 34,64 56,30 44,30" fill="url(#boltFill)"/>
  <polygon points="46,10 30,38 42,38 34,64 56,30 44,30" fill="none" stroke="#fff" stroke-width="0.8" opacity=".5"/>
  <circle cx="38" cy="20" r="2" fill="#fff" opacity=".7"/>
</g>
</svg>`,
	},
	{
		Name:   "Moon Prism",
		Rarity: "Rare",
		Art: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 80 72" width="80" height="72">
<defs>
  <radialGradient id="moonGlow" cx="50%" cy="50%" r="50%">
    <stop offset="0%" stop-color="#ce93d8" stop-opacity=".55"/>
    <stop offset="100%" stop-color="#4a148c" stop-opacity="0"/>
  </radialGradient>
  <linearGradient id="moonFill" x1="20%" y1="10%" x2="80%" y2="90%">
    <stop offset="0%" stop-color="#f3e5f5"/>
    <stop offset="60%" stop-color="#ce93d8"/>
    <stop offset="100%" stop-color="#6a1b9a"/>
  </linearGradient>
  <filter id="mf">
    <feGaussianBlur stdDeviation="2.5" result="b"/>
    <feMerge><feMergeNode in="b"/><feMergeNode in="SourceGraphic"/></feMerge>
  </filter>
  <style>
    @keyframes moonFloat{0%,100%{transform:translateY(0) rotate(-5deg)}50%{transform:translateY(-3px) rotate(5deg)}}
    @keyframes moonGlowP{0%,100%{opacity:.4;transform:scale(1)}50%{opacity:.7;transform:scale(1.18)}}
    @keyframes starTwinkle{0%,100%{opacity:.2}50%{opacity:1}}
    .mf2{transform-origin:40px 36px;animation:moonFloat 3.5s ease-in-out infinite}
    .mg2{transform-origin:40px 36px;animation:moonGlowP 3.5s ease-in-out infinite}
    .s1{animation:starTwinkle 1.2s ease-in-out infinite}
    .s2{animation:starTwinkle 1.8s ease-in-out .4s infinite}
    .s3{animation:starTwinkle 2.1s ease-in-out .9s infinite}
  </style>
</defs>
<ellipse class="mg2" cx="40" cy="36" rx="30" ry="26" fill="url(#moonGlow)"/>
<g class="mf2" filter="url(#mf)">
  <path d="M52 20 A20 20 0 1 0 52 52 A14 14 0 1 1 52 20Z" fill="url(#moonFill)"/>
  <circle cx="34" cy="26" r="2" fill="#fff" opacity=".5"/>
</g>
<circle class="s1" cx="62" cy="18" r="1.5" fill="#e1bee7"/>
<circle class="s2" cx="20" cy="22" r="1.2" fill="#e1bee7"/>
<circle class="s3" cx="58" cy="54" r="1" fill="#e1bee7"/>
<circle class="s1" cx="15" cy="46" r="1.3" fill="#e1bee7"/>
<circle class="s2" cx="66" cy="40" r="1" fill="#e1bee7"/>
</svg>`,
	},
	{
		Name:   "Cyber Orb",
		Rarity: "Legendary",
		Art: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 80 72" width="80" height="72">
<defs>
  <radialGradient id="orbGlow" cx="50%" cy="50%" r="50%">
    <stop offset="0%" stop-color="#00e5ff" stop-opacity=".65"/>
    <stop offset="100%" stop-color="#0d47a1" stop-opacity="0"/>
  </radialGradient>
  <radialGradient id="orbFill" cx="35%" cy="30%" r="65%">
    <stop offset="0%" stop-color="#e0f7fa"/>
    <stop offset="40%" stop-color="#00bcd4"/>
    <stop offset="100%" stop-color="#0d47a1"/>
  </radialGradient>
  <filter id="of">
    <feGaussianBlur stdDeviation="3" result="b"/>
    <feMerge><feMergeNode in="b"/><feMergeNode in="SourceGraphic"/></feMerge>
  </filter>
  <style>
    @keyframes orbPulse{0%,100%{transform:scale(1);opacity:.5}50%{transform:scale(1.22);opacity:.8}}
    @keyframes orbRing{from{transform:rotate(0deg)}to{transform:rotate(360deg)}}
    @keyframes orbRingRev{from{transform:rotate(0deg)}to{transform:rotate(-360deg)}}
    .og{transform-origin:40px 36px;animation:orbPulse 2s ease-in-out infinite}
    .or1{transform-origin:40px 36px;animation:orbRing 4s linear infinite}
    .or2{transform-origin:40px 36px;animation:orbRingRev 3s linear infinite}
  </style>
</defs>
<ellipse class="og" cx="40" cy="36" rx="30" ry="26" fill="url(#orbGlow)"/>
<circle cx="40" cy="36" r="18" fill="url(#orbFill)" filter="url(#of)"/>
<g class="or1">
  <ellipse cx="40" cy="36" rx="24" ry="8" fill="none" stroke="#00e5ff" stroke-width="1.2" opacity=".5"/>
</g>
<g class="or2">
  <ellipse cx="40" cy="36" rx="8" ry="24" fill="none" stroke="#40c4ff" stroke-width="1" opacity=".4"/>
</g>
<circle cx="40" cy="36" r="18" fill="none" stroke="#80deea" stroke-width="0.8" opacity=".3"/>
<circle cx="33" cy="28" r="4" fill="#fff" opacity=".4"/>
<circle cx="33" cy="28" r="2" fill="#fff" opacity=".6"/>
</svg>`,
	},
}
