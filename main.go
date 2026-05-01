package main

import (
	"flag"
	"fmt"
	"gopher-gotchi/internal/api"
	"gopher-gotchi/internal/brain"
	"gopher-gotchi/internal/tray"
	"gopher-gotchi/internal/ui"
	"gopher-gotchi/internal/watcher"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func main() {
	speciesFlag := flag.String("species", "diana", "The species of the companion")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	if handleCLICommands() {
		return
	}

	myPet := loadOrCreatePet(*speciesFlag)

	w := startWatcher(myPet)

	api.StartServer(myPet)

	go runLoop(myPet)

	tray.Init(func() {
		myPet.Save()
		w.Close()
		os.Exit(0)
	})
}

func handleCLICommands() bool {
	if len(os.Args) > 2 && os.Args[1] == "tell" {
		command := os.Args[2]
		resp, err := http.Get("http://localhost:9090/tell?cmd=" + command)
		if err != nil {
			fmt.Println("❌ Could not reach Diana. Maybe she's sleeping?")
			return true
		}
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("📡 Diana responds: %s\n", string(body))
		return true
	}
	return false
}

func loadOrCreatePet(species string) *brain.Pet {
	myPet, err := brain.LoadPet()
	if err != nil {
		myPet = brain.NewPet("Diana", species)
	}
	return myPet
}

func startWatcher(myPet *brain.Pet) *watcher.Watcher {
	home, _ := os.UserHomeDir()
	devPath := filepath.Join(home, "Development")
	w := watcher.NewWatcher()
	w.Start(devPath, myPet)
	return w
}

func runLoop(myPet *brain.Pet) {
	go myPet.LifeCycle()

	uiTicker := time.NewTicker(2 * time.Second)
	quoteTicker := time.NewTicker(2 * time.Minute)
	defer uiTicker.Stop()
	defer quoteTicker.Stop()

	for {
		select {
		case <-uiTicker.C:
			myPet.UpdateVitals()
			face := myPet.GetFace()
			if (face == ui.Themes[myPet.Species].Happy || face == ui.Themes[myPet.Species].Neutral) && rand.Intn(5) == 0 {
				myPet.Log(myPet.GetRandomQuote())
				face = myPet.GetBlinkFace()
			}
			ui.DrawPet(face, myPet.Level, myPet.Hunger, myPet.Mood, myPet.Messages, myPet.CPULoad)
			tray.Update(myPet.Level, myPet.Hunger, myPet.Mood)
		case <-quoteTicker.C:
			myPet.Log(myPet.GetRandomQuote())
		}
	}
}
