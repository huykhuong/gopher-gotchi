package main

import (
	"gopher-gotchi/internal/brain"
	"gopher-gotchi/internal/tray"
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

	myPet, err := brain.LoadPet()
	if err != nil {
		myPet = brain.NewPet("Gopher")
	}

	home, _ := os.UserHomeDir()
	devPath := filepath.Join(home, "Development")
	w := watcher.NewWatcher()

	w.Start(devPath, func(lines int) {
		myPet.Eat(lines)
		myPet.Save()
	})

	// Launch logic & UI in a Goroutine
	// Because the Tray needs the main thread
	go func() {
		go myPet.LifeCycle()

		for {
			face := myPet.GetFace()
			if (face == ui.FaceHappy || face == ui.FaceNeutral) && rand.Intn(5) == 0 {
				face = ui.FaceBlink
			}

			ui.DrawPet(face, myPet.Level, myPet.Hunger, myPet.Mood, myPet.Messages)

			tray.Update(myPet.Level, myPet.Hunger, myPet.Mood)

			time.Sleep(1 * time.Second)
		}
	}()

	// Start the tray (This is a blocking operation)
	tray.Init(func()  {
		myPet.Save()
		w.Close()
		os.Exit(0)
	})
}