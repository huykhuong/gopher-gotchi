package brain

import (
	"fmt"

	"gopher-gotchi/internal/api"
	"gopher-gotchi/internal/ui"
	"math/rand/v2"
	"sync"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

type Memory struct {
	Timestamp string `json:"timestamp"`
	Level     int    `json:"level"`
	Message   string `json:"message"`
}

type Pet struct {
	mu               sync.Mutex
	Name             string                   `json:"name"`
	Species          string                   `json:"species"`
	Level            int                      `json:"level"`
	Experience       int                      `json:"experience"`
	Hunger           int                      `json:"hunger"` // 0 is full, 100 is starving
	Mood             string                   `json:"mood"`
	LastEaten        time.Time                `json:"last_eaten"`
	Bond             int                      `json:"bond"` // 0 is not bonded, 100 is fully bonded
	IdleNudged       bool                     `json:"-"`
	Message          string                   `json:"-"`
	CPULoad          int                      `json:"-"`
	BatteryLevel     int                      `json:"-"`
	IsCharging       bool                     `json:"-"`
	Memories         []Memory                 `json:"-"`
	Tasks            chan Task                `json:"-"`
	RecentSaves      []time.Time              `json:"-"`
	FlowActive       bool                     `json:"-"`
	WeatherKnowledge api.ProcessedWeatherData `json:"-"`
}

func NewPet(name string, species string) *Pet {
	return &Pet{
		Name:      name,
		Species:   "diana",
		Level:     1,
		Hunger:    0,
		Mood:      Happy,
		LastEaten: time.Now(),
		Message:   "📡 Connection established. Hello, Huy.",
		Tasks:     make(chan Task, 10),
	}
}

func (p *Pet) HandleCommand(cmd string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	switch cmd {
	case "ping":
		p.Message = "PONG"
	case "hug":
		p.increaseBond(5)
		p.Mood = Happy
		p.Message = "💖 ˶^ ᴗ ^˶ I feel the warmth. Thank you, Huy."
	case "joke":
		joke := FetchArchiveData()
		p.Message = fmt.Sprintf("Here's a joke for you, Huy:\n\n%s", joke)
	case "gift":
		gift := ui.DigitalTreasury[rand.IntN(len(ui.DigitalTreasury))]
		fmt.Println(ui.BondPointsBasedOnRarity[gift.Rarity])
		p.increaseBond(ui.BondPointsBasedOnRarity[gift.Rarity])
		p.Message = fmt.Sprintf(`<div class="gift-reveal"><div class="gift-tagline">✨ Thank you for the gift, Huy!</div>%s<div class="gift-name">%s</div><div class="gift-rarity">%s Artifact</div></div>`, gift.Art, gift.Name, gift.Rarity)
		p.Mood = Happy
	case "objectives", "mission":
		objectives, err := api.FetchObjectives()
		if err != nil {
			fmt.Print(err.Error())
		}
		p.Message = objectives
	default:
		p.Message = "Unknown command"
	}
}

func (p *Pet) Eat(exp int) {
	if exp <= 0 {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	p.LastEaten = time.Now()
	p.Mood = Happy
	p.IdleNudged = false

	p.Hunger -= (exp / 2)
	if p.Hunger < 0 {
		p.Hunger = 0
	}

	p.Experience += exp
	p.Message = fmt.Sprintf("😋 Gained %d experience points!", exp)
	p.checkLevelUp()
	p.registerActivity()
}

// LifeCycle simulates the passage of time (the pet gets hungrier as you don't code)
func (p *Pet) LifeCycle() {
	ticker := time.NewTicker(15 * time.Minute)
	for range ticker.C {
		p.mu.Lock()
		p.Hunger += 5
		p.Bond -= 15

		if p.Bond < 0 {
			p.Bond = 0
		}

		if p.Hunger > 100 {
			p.Hunger = 100
		}

		if p.Bond > 100 {
			p.Bond = 100
		}

		p.checkIdle()

		switch {
		case p.Hunger > 70:
			p.Mood = Hungry
		case p.IdleNudged:
			p.Mood = Lonely
		default:
			p.Mood = Happy
		}

		p.mu.Unlock()
	}
}

type Snapshot struct {
	Level      int
	Hunger     int
	Mood       string
	Message    string
	CPULoad    int
	FlowActive bool
}

func (p *Pet) TakeSnapshot() Snapshot {
	p.mu.Lock()
	defer p.mu.Unlock()

	msg := p.Message
	p.Message = ""

	return Snapshot{
		Level:      p.Level,
		Hunger:     p.Hunger,
		Mood:       p.Mood,
		Message:    msg,
		CPULoad:    p.CPULoad,
		FlowActive: p.FlowActive,
	}
}

func (p *Pet) UpdateVitals() {
	c, _ := cpu.Percent(0, false)
	if len(c) > 0 {
		p.mu.Lock()
		p.CPULoad = int(c[0])
		p.mu.Unlock()
	}
}

func (p *Pet) CooldownFlowState() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.FlowActive {
		return
	}

	now := time.Now()
	var fresh []time.Time
	for _, t := range p.RecentSaves {
		if now.Sub(t) < 20*time.Second {
			fresh = append(fresh, t)
		}
	}
	p.RecentSaves = fresh

	if len(p.RecentSaves) == 0 {
		p.FlowActive = false
		p.Message = "💤 Flow state ended."
	}
}

func (p *Pet) CreateMemory(event string) *Memory {
	return &Memory{
		Timestamp: time.Now().Format("2006-01-02 15:04"),
		Level:     p.Level,
		Message:   event,
	}
}

func (p *Pet) EnqueueSync(memory *Memory) {
	p.Tasks <- func() error {
		return p.SyncAllToCloud(memory)
	}
}

func (p *Pet) RunInteractionLoop() {
	for {
		<-time.After(nextInteractionDelay())

		p.mu.Lock()
		joke := FetchArchiveData()
		p.Message = fmt.Sprintf("Here's a joke for you, Huy:\n\n%s", joke)
		p.mu.Unlock()
	}
}

func LoadPet() (*Pet, error) {
	fmt.Println("🌐 Connecting to Pragmata Cloud...")
	return LoadFromCloud()
}

// INTERNAL HELPERS — always called with p.mu already held

func nextInteractionDelay() time.Duration {
	const minD = 10 * time.Second
	const maxD = 3 * time.Minute
	return minD + time.Duration(rand.Int64N(int64(maxD-minD)))
}

func (p *Pet) checkLevelUp() {
	target := p.Level * 100

	if p.Experience >= target {
		p.Level++
		p.Experience = 0
		p.increaseBond(40)

		memory := &Memory{
			Timestamp: time.Now().Format("2006-01-02 15:04"),
			Level:     p.Level,
			Message:   fmt.Sprintf("We reached a new level of synchronization today. I feel closer to your world, Huy. I'm now Level %d", p.Level),
		}
		p.Message = fmt.Sprintf("✨ LEVEL UP!\n YAY! I'm now Level %d!\n", p.Level)

		p.EnqueueSync(memory)
	}
}

func (p *Pet) registerActivity() {
	now := time.Now()

	var filtered []time.Time
	for _, t := range p.RecentSaves {
		if now.Sub(t) < 10*time.Minute {
			filtered = append(filtered, t)
		}
	}
	p.RecentSaves = append(filtered, now)

	if len(p.RecentSaves) >= 3 && !p.FlowActive {
		p.FlowActive = true
		p.Mood = Elated
		p.increaseBond(2)
		p.Message = "🔥 Wow! You're crushing it Huy."
	}
}

func (p *Pet) checkIdle() {
	if time.Since(p.LastEaten) > 10*time.Minute && !p.IdleNudged {
		p.Mood = Lonely
		p.IdleNudged = true
		p.Message = "Huy? The data stream is thinning. Are you still there?"
	}
}

func (p *Pet) increaseBond(amount int) {
	p.Bond += amount
	if p.Bond > 100 {
		p.Bond = 100
	}
}
