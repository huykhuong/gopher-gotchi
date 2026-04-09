package main

import (
	"fmt"
	"gopher-gotchi/internal/brain"
	"gopher-gotchi/internal/ui"
	"gopher-gotchi/internal/watcher"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

func main() {
	// Seed the random number generator so the blinks aren't predictable
	rand.Seed(time.Now().UnixNano())

	myPet := brain.NewPet("Gopher")
	go myPet.LifeCycle()

	home, _ := os.UserHomeDir()
	devPath := filepath.Join(home, "Development")

	w := watcher.NewWatcher()
	defer w.Close()

	w.Start(devPath, func(lines int) {
		myPet.Eat(lines)
	})

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C{
		fmt.Print("\033[H\033[2J") // Clear screen
		
		face := myPet.GetFace()

		// Random blink logic
		if (face == ui.FaceHappy || face == ui.FaceNeutral) && rand.Intn(5) == 0 {
			face = ui.FaceBlink
		}

		ui.DrawPet(face, myPet.Level, myPet.Hunger, myPet.Mood)

		time.Sleep(500 * time.Millisecond)
		
	}
}