package watcher

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"gopher-gotchi/internal/brain"

	"github.com/fsnotify/fsnotify"
)

var sizeCache = make(map[string]int64)

var folderToRegister = []string{
	"bionic",
	"gopher-gotchi",
}

type Watcher struct {
	fsWatcher *fsnotify.Watcher
}

func NewWatcher() *Watcher {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	return &Watcher{fsWatcher: w}
}

func (w *Watcher) Start(rootPath string, eventChan chan<- brain.DataEvent) {
	if err := w.registerDirs(rootPath); err != nil {
		log.Fatal("Search error:", err)
	}

	go func() {
		for {
			select {
			case event, ok := <-w.fsWatcher.Events:
				if !ok {
					return
				}

				if event.Op&fsnotify.Create == fsnotify.Create {
					if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
						_ = w.registerDirs(event.Name)
					}
				}

				if event.Op&fsnotify.Write == fsnotify.Write {
					eventChan <- brain.DataEvent{
						Type: 		brain.FileSaved,
						Payload: 	event.Name,
					}
				}
			case err, ok := <-w.fsWatcher.Errors:
				if !ok {
					return
				}
				log.Println("Watcher error:", err)
			}
		}
	}()
}

func (w *Watcher) registerDirs(rootPath string) error {
	return filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		name := d.Name()
		if d.IsDir() && (name == "" || name[0] == '.' || name == "node_modules") {
			return filepath.SkipDir
		}

		if !underRegisteredFolder(path) {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return nil
		}

		if !d.IsDir() {
			sizeCache[path] = info.Size()
			return nil
		}

		return w.fsWatcher.Add(path)
	})
}

func underRegisteredFolder(path string) bool {
	for p := path; ; {
		base := filepath.Base(p)
		for _, name := range folderToRegister {
			if base == name {
				return true
			}
		}
		parent := filepath.Dir(p)
		if parent == p {
			return false
		}
		p = parent
	}
}

func (w *Watcher) Close() {
	w.fsWatcher.Close()
}

func UpdateXP(path string, myPet *brain.Pet) {
	info, err := os.Stat(path)
	if err != nil {
		return
	}

	currentSize := info.Size()
	lastSize, seen := sizeCache[path]
	sizeCache[path] = currentSize

	if !seen {
		return
	}

	if diff := currentSize - lastSize; diff > 0 {
		if xp := int(diff / 10); xp > 0 {
			myPet.Eat(xp)
		}
	}
}