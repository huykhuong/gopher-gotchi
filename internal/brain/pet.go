package brain

import (
	"fmt"
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
	mu           sync.RWMutex
	Name         string    `json:"name"`
	Species      string    `json:"species"`
	Level        int       `json:"level"`
	Experience   int       `json:"experience"`
	Hunger       int       `json:"hunger"` // 0 is full, 100 is starving
	Mood         string    `json:"mood"`
	LastEaten    time.Time `json:"last_eaten"`
	IdleNudged   bool      `json:"-"`
	Message      string  `json:"-"` // No need to save the log to JSON
	CPULoad      int       `json:"-"`
	BatteryLevel int       `json:"-"`
	IsCharging   bool      `json:"-"`
	Memories     []Memory  `json:"-"` // We don't save this in diana.json
	Tasks        chan Task `json:"-"`
}

func NewPet(name string, species string) *Pet {
	if _, ok := ui.Themes[species]; !ok {
		species = "diana" // default to diana if the theme is not found
	}

	p := &Pet{
		Name:      name,
		Species:   species,
		Level:     1,
		Hunger:    0,
		Mood:      "Happy",
		LastEaten: time.Now(),
		Message:   "📡 Connection established. Hello, Huy.",
		Tasks:     make(chan Task, 10),
	}

	return p
}

func (p *Pet) HandleCommand(cmd string) {
	switch cmd {
	case "ping":
		p.Log("PONG")
	case "hug":
		p.Mood = "Happy 😊"
		p.Log("💖 ˶^ ᴗ ^˶ I feel the warmth. Thank you, Huy.")
	default:
		p.Log("Unknown command")
	}
}

func (p *Pet) GetBlinkFace() string {
	return ui.Themes[p.Species].Blink
}

func (p *Pet) Log(msg string) {
	p.Message = msg
}

func (p *Pet) Eat(exp int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if exp <= 0 {
		return
	}

	p.LastEaten = time.Now()
	p.Mood = "Happy 😊"
	p.IdleNudged = false

	p.Hunger -= (exp / 2)
	if p.Hunger < 0 {
		p.Hunger = 0
	}

	p.Experience += exp

	p.Log(fmt.Sprintf("😋 Gained %d experience points!", exp))
	p.checkLevelUp()
}

func (p *Pet) CheckIdle() {
	if time.Since(p.LastEaten) > 10*time.Minute && !p.IdleNudged {
		p.Mood = "Lonely 💔"
		p.IdleNudged = true

		p.Log("Huy? The data stream is thinning. Are you still there?")
	}
}

func (p *Pet) checkLevelUp() {
	target := p.Level * 100
	if p.Experience >= target {
		p.Level++
		p.Experience = 0

		memory := p.CreateMemory(fmt.Sprintf("We reached a new level of synchronization today. I feel closer to your world, Huy. I'm now Level %d", p.Level))
		p.EnqueueSync(memory)

		// Internal log
		msg := fmt.Sprintf("✨ LEVEL UP!\n YAY! I'm now Level %d!\n", p.Level)
		p.Log(msg)
	}
}

// LifeCycle simulates the passage of time (the pet gets hungrier as you don't code)
func (p *Pet) LifeCycle() {
	ticker := time.NewTicker(15 * time.Minute)
	for range ticker.C {
		p.Hunger += 5

		// Check if we should nudge the user.
		p.CheckIdle()

		if p.Hunger > 100 {
			p.Hunger = 100
			p.Mood = "Starving 💀"
		} else if p.Hunger > 70 {
			p.Mood = "Grumpy 💢"
		} else if p.IdleNudged {
			p.Mood = "Lonely 💔"
		} else {
			p.Mood = "Happy 😊"
		}
	}
}

func (p *Pet) UpdateVitals() {
	c, _ := cpu.Percent(0, false)
	if len(c) > 0 {
		p.CPULoad = int(c[0])
	}
}

func (p *Pet) CreateMemory(event string) *Memory {
	newMemory := Memory{
		Timestamp: time.Now().Format("2006-01-02 15:04"),
		Level:     p.Level,
		Message:   event,
	}

	p.Log("💾 Memory Archived: " + event)

	return &newMemory
}

func (p *Pet) EnqueueSync(memory *Memory) {
	p.Tasks <- func() error {
		return p.SyncAllToCloud(memory)
	}
}

func (p *Pet) RunJokeLoop() {
	for {
		<-time.After(nextJokeDelay())
		joke := FetchArchiveData()
		p.Log(fmt.Sprintf("Here's a joke for you, Huy:\n\n%s", joke))
	}
}

// PERSISTENCE LOGIC
func LoadPet() (*Pet, error) {
	fmt.Println("🌐 Connecting to Pragmata Cloud...")
	return LoadFromCloud()
}

func nextJokeDelay() time.Duration {
	const minD = 10 * time.Second
	const maxD = 3 * time.Minute
	return minD + time.Duration(rand.Int64N(int64(maxD-minD)))
}
