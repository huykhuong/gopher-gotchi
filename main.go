package main

import (
	"flag"
	"fmt"
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
	speciesFlag := flag.String("species", "diana", "The species of the companion")
	flag.Parse()

	// Seed the random number generator so the blinks aren't predictable
	rand.Seed(time.Now().UnixNano())

	myPet, err := brain.LoadPet()
	if err != nil {
		myPet = brain.NewPet("Diana", *speciesFlag)
	}

	home, _ := os.UserHomeDir()
	devPath := filepath.Join(home, "Development")

	w := watcher.NewWatcher()

	w.Start(devPath, myPet)

	uiTicker := time.NewTicker(2 * time.Second)
	quoteTicker := time.NewTicker(2 * time.Minute)
	defer uiTicker.Stop()
	defer quoteTicker.Stop()

	// Launch logic & UI in a Goroutine
	// Because the Tray needs the main thread
	go func() {
		go myPet.LifeCycle()

		fmt.Println("test")

		for {
			select {
			case <- uiTicker.C:
				myPet.UpdateVitals()
				face := myPet.GetFace()
				if (face == ui.Themes[myPet.Species].Happy || face == ui.Themes[myPet.Species].Neutral) && rand.Intn(5) == 0 {
					myPet.Log(myPet.GetRandomQuote())
					face = myPet.GetBlinkFace()
				}

				ui.DrawPet(face, myPet.Level, myPet.Hunger, myPet.Mood, myPet.Messages, myPet.CPULoad)

				tray.Update(myPet.Level, myPet.Hunger, myPet.Mood)
			case <- quoteTicker.C:
				myPet.Log(myPet.GetRandomQuote())
			}
		}
	}()

	// Start the tray (This is a blocking operation)
	tray.Init(func()  {
		myPet.Save()
		w.Close()
		os.Exit(0)
	})
}