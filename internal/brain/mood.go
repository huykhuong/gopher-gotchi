package brain

import (
	"math/rand/v2"
)

const (
	Happy  = "happy"
	Hungry = "hungry"
	Lonely = "lonely"
	Elated = "elated"
	Idle   = "idle"
)

type Mood struct {
	Name     string
	Row      int
	Frames   int
	Duration int
}

var happyVariants = []Mood{
	{Row: 3, Frames: 4, Duration: 1000, Name: "Happy 😊"},
	{Row: 1, Frames: 8, Duration: 1060, Name: "Happy 😊"},
	{Row: 7, Frames: 6, Duration: 820, Name: "Happy 😊"},
}

var moods = map[string]Mood{
	"hungry": {
		Name:     "Starving 💀",
		Row:      5,
		Frames:   8,
		Duration: 1220,
	},
	"lonely": {
		Name:     "Lonely 💔",
		Row:      6,
		Frames:   6,
		Duration: 1010,
	},
	"elated": {
		Name:     "Elated 😍",
		Row:      7,
		Frames:   6,
		Duration: 820,
	},
	"idle": {
		Name:     "Idle 😴",
		Row:      0,
		Frames:   6,
		Duration: 1100,
	},
}

func GetAnimationForMood(mood string) Mood {
	if mood == Happy {
		return happyVariants[rand.IntN(len(happyVariants))]
	}

	return moods[mood]
}
