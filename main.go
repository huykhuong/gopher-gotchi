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
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/gen2brain/beeep"
)

func main() {
	var wg sync.WaitGroup
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt, syscall.SIGTERM)

	speciesFlag := flag.String("species", "diana", "The species of the companion")
	flag.Parse()

	if handleCLICommands() {
		return
	}

	myPet := loadOrCreatePet(*speciesFlag)

	eventsChan := make(chan brain.DataEvent, 10)
	tasks := make(chan brain.Task, 10)

	w := startWatcher(eventsChan)
	api.StartServer(myPet)
	brain.StartWorkerPool(tasks)

	go runEventLoop(eventsChan, tasks, myPet)
	go runLoop(myPet)

	<-stopSignal

	tray.Init(func() {
		close(eventsChan)

		wg.Add(1)
		tasks <- func() error {
			defer wg.Done()
			return brain.SaveToCloud(myPet)
		}
		wg.Wait()
		close(tasks)

		fmt.Println("💾 Finalizing cloud sync... Goodbye, Huy.")
		time.Sleep(2 * time.Second)

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

func startWatcher(eventChan chan<- brain.DataEvent) *watcher.Watcher {
	home, _ := os.UserHomeDir()
	devPath := filepath.Join(home, "Development")
	w := watcher.NewWatcher()
	w.Start(devPath, eventChan)
	return w
}

func runEventLoop(eventsChan chan brain.DataEvent, tasks chan brain.Task, myPet *brain.Pet) {
	for event := range eventsChan {
		switch event.Type {
		case brain.FileSaved:
			path := event.Payload.(string)
			watcher.UpdateXP(path, myPet)

			tasks <- func() error {
				return brain.SaveToCloud(myPet)
			}
		}
	}
}

func runLoop(myPet *brain.Pet) {
	go myPet.LifeCycle()

	uiTicker := time.NewTicker(2 * time.Second)
	defer uiTicker.Stop()

	for range uiTicker.C {
		myPet.UpdateVitals()
		face := myPet.GetFace()
		if (face == ui.Themes[myPet.Species].Happy || face == ui.Themes[myPet.Species].Neutral) && rand.Intn(5) == 0 {
			face = myPet.GetBlinkFace()
		}
		ui.DrawPet(face, myPet.Level, myPet.Hunger, myPet.Mood, myPet.Messages, myPet.CPULoad)
		tray.Update(myPet.Level, myPet.Hunger, myPet.Mood)

		hour := time.Now().Hour()
		if hour >= 23 {
			// Desktop notification
			beeep.Alert("Go easy on yourself 🌙", "Huy, it's late. Don't forget to rest.", "")
		}
	}
}
