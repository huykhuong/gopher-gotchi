package brain

import (
	"encoding/json"
	"fmt"
	"gopher-gotchi/internal/ui"
	"os"
	"path/filepath"
	"time"
)

type Pet struct {
	Name 				string		`json:"name"`
	Species			string		`json:"species"`
	Level				int				`json:"level"`
	Experience	int				`json:"experience"`
	Hunger			int 			`json:"hunger"` // 0 is full, 100 is starving 
	Mood				string		`json:"mood"`
	LastEaten		time.Time	`json:"last_eaten"`
	Messages		[]string	`json:"-"` // No need to save the log to JSON
}

func NewPet(name string, species string) *Pet {
	if _, ok := ui.Themes[species]; !ok {
		species = "gopher" // default to gopher if the theme is not found
	}

	return &Pet{
		Name:				name,
		Species:		species,
		Level:			1,
		Hunger:			0,
		Mood:				"Happy",
		LastEaten: 	time.Now(),
		Messages:  []string{"🐣 I'm alive!"},
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

func (p *Pet) Eat(linesChanged int) {
	if linesChanged <= 0 {
		return
	}

	p.Hunger -= (linesChanged / 10)
	if p.Hunger < 0 {
		p.Hunger = 0
	}

	p.Experience += linesChanged
	p.checkLevelUp()

	p.Log(fmt.Sprintf("😋 %s ate %d lines of code! (Exp: %d)\n", p.Name, linesChanged, p.Experience))
}

func (p *Pet) checkLevelUp() {
	target := p.Level * 100
	if p.Experience >= target {
		p.Level++
		p.Experience = 0
		p.Log(fmt.Sprintf("✨ LEVEL UP! %s is now Level %d!\n", p.Name, p.Level))
	}
}

// LifeCycle simulates the passage of time (the pet gets hungrier as you don't code)
func (p *Pet) LifeCycle() {
	ticker := time.NewTicker(15 * time.Minute)
	for range ticker.C {
		p.Hunger += 10
		if p.Hunger > 100 {
			p.Hunger = 100
			p.Mood = "Starving 💀"
		} else if p.Hunger > 70 {
			p.Mood = "Grumpy 💢"
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