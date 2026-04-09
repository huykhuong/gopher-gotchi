package brain

import (
	"fmt"
	"gopher-gotchi/internal/ui"
	"time"
)

type Pet struct {
	Name 				string
	Level				int
	Experience	int
	Hunger			int // 0 is full, 100 is starving
	Mood				string
	LastEaten		time.Time
}

func NewPet(name string) *Pet {
	return &Pet{
		Name:				name,
		Level:			1,
		Hunger:			0,
		Mood:				"Happy",
		LastEaten: 	time.Now(), 
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

	fmt.Printf("😋 %s ate %d lines of code! (Exp: %d)\n", p.Name, linesChanged, p.Experience)
}

func (p *Pet) checkLevelUp() {
	target := p.Level * 100
	if p.Experience >= target {
		p.Level++
		p.Experience = 0
		fmt.Printf("✨ LEVEL UP! %s is now Level %d!\n", p.Name, p.Level)
	}
}

// LifeCycle simulates the passage of time (the pet gets hungrier as you don't code)
func (p *Pet) LifeCycle() {
	ticker := time.NewTicker(15 * time.Minute)
	for range ticker.C {
		p.Hunger += 5
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
	if p.Hunger >= 100 {
		return ui.FaceDead
	}
	if p.Hunger > 70 {
		return ui.FaceHungry
	}
	if p.Mood == "Happy" {
		return ui.FaceHappy
	}

	return ui.FaceNeutral
}