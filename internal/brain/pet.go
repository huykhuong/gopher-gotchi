package brain

import (
	"encoding/json"
	"fmt"
	"gopher-gotchi/internal/ui"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/shirou/gopsutil/cpu"
)

type Pet struct {
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

// Random dialogue for Diana
var dianaQuotes = []string{
	"System check complete. Everything is nominal, Huy.",
	"I'm scanning the data streams. You're doing great.",
	"Your logic patterns are fascinating.",
	"Are we pushing the boundaries today?",
	"I'm glad to be synced with this unit.",
	"You're making progress, Huy. Keep it up!",
	"The system is stable. Ready for more tasks.",
	"I'm here to support you. Anything specific?",
	"You're on the right track. Keep coding!",
	"I'm monitoring the system. Everything is under control.",
	"You're making great strides. Keep pushing!",
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
	default:
		return "Unknown command"
	}
}

func (p *Pet) GetRandomQuote() string {
	return dianaQuotes[rand.Intn(len(dianaQuotes))]
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
	if exp <= 0 {
		return
	}

	p.LastEaten = time.Now()
	p.IdleNudged = false

	p.Hunger -= (exp / 10)
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
		p.Hunger += 10

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
	theme := ui.Themes[p.Species]

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

func GetConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".gopher-gotchi.json")
}

func (p *Pet) Save() error {
	data, err := json.MarshalIndent(p, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(GetConfigPath(), data, 0644)
}

func LoadPet() (*Pet, error) {
	data, err := os.ReadFile(GetConfigPath())
	if err != nil {
		return nil, err
	}
	var p Pet
	err = json.Unmarshal(data, &p)
	return &p, err
}