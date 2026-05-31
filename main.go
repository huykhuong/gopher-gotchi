package main

import (
	"context"
	"flag"
	"fmt"
	"gopher-gotchi/internal/api"
	"gopher-gotchi/internal/brain"
	"gopher-gotchi/internal/watcher"
	"gopher-gotchi/internal/window"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

func main() {
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	devFlag := flag.Bool("dev", false, "Load UI from filesystem for live editing")
	flag.Parse()

	myPet := loadOrCreatePet()

	eventsChan := make(chan brain.DataEvent, 10)

	// Create the window first — before the file watcher opens thousands of
	// file descriptors, which would exhaust macOS's limit and crash Metal init.
	var win *window.Window

	workerDone := brain.StartWorkerPool(myPet.Tasks)

	var cleanupOnce sync.Once
	cleanupDone := make(chan struct{})

	cleanup := func() {
		cleanupOnce.Do(func() {
			defer close(cleanupDone)

			// Stop new events first, then drain in-flight sync tasks.
			close(eventsChan)
			close(myPet.Tasks)
			<-workerDone

			// After all tasks are drained, do the final sync.
			if err := myPet.SyncAllToCloud(nil); err != nil {
				log.Println("⚠️ Final sync error:", err)
			} else {
				log.Println("💾 Cloud sync complete. Goodbye, Huy.")
			}

			if win != nil {
				win.Terminate()
			}
		})
	}

	win = window.New(cleanup, myPet.HandleCommand, *devFlag)

	startWatcher(eventsChan)
	brain.StartClipboardWatcher(ctx, eventsChan)
	api.StartServer(window.Spritesheet)

	go runEventLoop(eventsChan, myPet)
	go runLoop(myPet, win)

	go func() {
		<-stopSignal
		cleanup()
	}()

	// Run blocks the main thread until win.Terminate() is called (by cleanup).
	win.Run()

	<-cleanupDone
}

func loadOrCreatePet() *brain.Pet {
	myPet, err := brain.LoadPet()
	if err != nil {
		myPet = brain.NewPet("Diana")
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
		case brain.ClipboardErrorDetected:
			myPet.Message = fmt.Sprintf("Don't sweat the %s, Huy. Take a deep breath, we can trace this block together!", event.Payload.(string))
			myPet.Mood = brain.Concerned
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
