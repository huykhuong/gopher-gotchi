package main

import (
	"flag"
	"fmt"
	"gopher-gotchi/internal/api"
	"gopher-gotchi/internal/brain"
	"gopher-gotchi/internal/watcher"
	"gopher-gotchi/internal/window"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

func main() {
	var wg sync.WaitGroup
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt, syscall.SIGTERM)

	speciesFlag := flag.String("species", "diana", "The species of the companion")
	devFlag := flag.Bool("dev", false, "Load UI from filesystem for live editing")
	flag.Parse()

	myPet := loadOrCreatePet(*speciesFlag)

	eventsChan := make(chan brain.DataEvent, 10)

	// Create the window first — before the file watcher opens thousands of
	// file descriptors, which would exhaust macOS's limit and crash Metal init.
	var win *window.Window
	cleanup := func() {
		if win != nil {
			win.Terminate()
		}

		close(eventsChan)

		wg.Add(1)
		myPet.Tasks <- func() error {
			defer wg.Done()
			return myPet.SyncAllToCloud(nil)
		}
		wg.Wait()
		close(myPet.Tasks)

		fmt.Println("💾 Finalizing cloud sync... Goodbye, Huy.")
		time.Sleep(2 * time.Second)

		os.Exit(0)
	}

	win = window.New(cleanup, myPet.HandleCommand, *devFlag)

	startWatcher(eventsChan)
	api.StartServer(window.Spritesheet)
	brain.StartWorkerPool(myPet.Tasks)

	go runEventLoop(eventsChan, myPet)
	go runLoop(myPet, win)

	go func() {
		<-stopSignal
		cleanup()
	}()

	// Run blocks the main thread until the window is closed.
	win.Run()
}

func loadOrCreatePet(species string) *brain.Pet {
	myPet, err := brain.LoadPet()
	if err != nil {
		myPet = brain.NewPet("Diana", species)
	}
	return myPet
}

func startWatcher(eventChan chan<- brain.DataEvent) *watcher.Watcher {
	home, _ := os.UserHomeDir()
	devPath := filepath.Join(home, "Development")
	w := watcher.NewWatcher()
	w.Start(devPath, eventChan)
	return w
}

func runEventLoop(eventsChan chan brain.DataEvent, myPet *brain.Pet) {
	for event := range eventsChan {
		switch event.Type {
		case brain.FileSaved:
			path := event.Payload.(string)
			watcher.UpdateXP(path, myPet)

			myPet.EnqueueSync(nil)
		}
	}
}

func runLoop(myPet *brain.Pet, win *window.Window) {
	go myPet.LifeCycle()
	go myPet.RunInteractionLoop()

	if weather, err := api.GetWeatherData(); err == nil {
		myPet.WeatherKnowledge = weather
	}

	uiTicker := time.NewTicker(2 * time.Second)
	saveStreakCooldownTicker := time.NewTicker(20 * time.Second)
	weatherTicker := time.NewTicker(15 * time.Minute)
	defer uiTicker.Stop()
	defer saveStreakCooldownTicker.Stop()
	defer weatherTicker.Stop()

	for {
		select {
		case <-uiTicker.C:
			myPet.UpdateVitals()
			snap := myPet.TakeSnapshot()

			win.Update(window.PetState{
				Level:            snap.Level,
				Hunger:           snap.Hunger,
				Mood:             brain.GetAnimationForMood(snap.Mood),
				Message:          snap.Message,
				CPULoad:          snap.CPULoad,
				FlowActive:       snap.FlowActive,
				WeatherKnowledge: myPet.WeatherKnowledge,
				Bond:             myPet.Bond,
			})

		case <-saveStreakCooldownTicker.C:
			myPet.CooldownFlowState()
		case <-weatherTicker.C:
			if weather, err := api.GetWeatherData(); err == nil {
				myPet.WeatherKnowledge = weather
			}
		}
	}
}
