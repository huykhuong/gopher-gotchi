package brain

import (
	"fmt"
	"gopher-gotchi/internal/ui"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/shirou/gopsutil/cpu"
)

type Pet struct {
	mu					sync.RWMutex
	Name 				string		`json:"name"`
	Species			string		`json:"species"`
	Level				int				`json:"level"`
	Experience	int				`json:"experience"`
	Hunger			int 			`json:"hunger"` // 0 is full, 100 is starving 
	Mood				string		`json:"mood"`
	LastEaten		time.Time	`json:"last_eaten"`
	IdleNudged	bool			`json:"-"`
	Messages		[]string	`json:"-"` // No need to save the log to JSON
	CPULoad			int				`json:"-"`
	BatteryLevel int			`json:"-"`
	IsCharging	bool			`json:"-"`
}

func NewPet(name string, species string) *Pet {
	if _, ok := ui.Themes[species]; !ok {
		species = "diana" // default to diana if the theme is not found
	}

	p := &Pet{
		Name:				name,
		Species:		species,
		Level:			1,
		Hunger:			0,
		Mood:				"Happy",
		LastEaten: 	time.Now(),
		Messages:  []string{"📡 Connection established. Hello, Huy."},
	}

	return p
}

func (p *Pet) HandleCommand(cmd string) string {
	switch cmd {
	case "ping":
		return "PONG"
	case "hug":
		p.Mood = "Happy 😊"
		p.Log("💖 Huy sent a virtual hug.")
		return "˶^ ᴗ ^˶ I feel the sync. Thank you, Huy."
	case "joke":
		p.Log("🔍 Diana is scanning terrestrial archives...")
		joke := fetchArchiveData()

		return fmt.Sprintf("I found this entry in the historical archives, Huy:\n\n%s", joke)
	default:
		return "Unknown command"
	}
}

func (p *Pet) GetBlinkFace() string {
	return ui.Themes[p.Species].Blink
}

func (p *Pet) Log(msg string) {
	p.Messages = append(p.Messages, msg)
	if len(p.Messages) > 5 {
		p.Messages = p.Messages[1:]
	}
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
	p.checkLevelUp()

	p.Log(fmt.Sprintf("😋 Gained %d experience points!", exp))
}

func (p *Pet) CheckIdle() {
	if time.Since(p.LastEaten) > 10*time.Minute && !p.IdleNudged {
		p.Mood = "Lonely 💔"
		p.IdleNudged = true

		msg := "Huy? The data stream is thinning. Are you still there?"
		beeep.Notify("Pragmata Protocol", msg, "")
		p.Log("📢 Sent a nudge for attention")
	}
}

func (p *Pet) checkLevelUp() {
	target := p.Level * 100
	if p.Experience >= target {
		p.Level++
		p.Experience = 0

		// Internal log
		msg := fmt.Sprintf("✨ LEVEL UP! %s is now Level %d!\n", p.Name, p.Level) 
		p.Log(msg)

		// Desktop notification
		err := beeep.Alert("Gopher-Gotchi Evolution", msg, "")
		if err != nil {
			p.Log(fmt.Sprintf("❌ Failed to send desktop notification: %v\n", err))
		}
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

func (p *Pet) GetFace() string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	theme := ui.Themes[p.Species]
	hour := time.Now().Hour()

	if hour >= 22 || hour < 5 {
		return `  (  ˶- ᴗ -˶ ) zZ`
	}

	if hour >= 5 && hour < 9 {
		return `  (  ˶O ᴗ O˶ ) ~`
	}

	if p.Hunger >= 100 {
		return theme.Dead
	}
	if p.Hunger > 70 {
		return theme.Hungry
	}
	if p.Mood == "Happy" {
		return theme.Happy
	}

	return theme.Neutral
}

func (p *Pet) UpdateVitals() {
	c, _ := cpu.Percent(0, false)
	if len(c) > 0 {
		p.CPULoad = int(c[0])
	}
}

// PERSISTENCE LOGIC
func LoadPet() (*Pet, error) {
	fmt.Println("🌐 Connecting to Pragmata Cloud...")
	return LoadFromCloud()
}

func fetchArchiveData() string {
	client := http.Client{
		Timeout: 3 * time.Second,
	}

	req, _ := http.NewRequest("GET", "https://icanhazdadjoke.com/", nil)
	req.Header.Set("Accept", "text/plain")

	resp, err := client.Do(req)
	if err != nil {
		return "Sorry! I can't connect to the archive."
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return string(body)
}