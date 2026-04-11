package watcher

import (
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

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

func (w *Watcher) Start(rootPath string, onSave func(lines int)) {
	// 1. Recursively add all subdirectories to the watcher
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			// Skip hidden folders like .git or .node_modules to save CPU
			if info.Name()[0] == '.' || info.Name() == "node_modules" {
				return filepath.SkipDir
			}
			return w.fsWatcher.Add(path)
		}
		return nil
	})

	if err != nil {
		log.Fatal("Search error:", err)
	}

	go func() {
		for {
			select {
			case event, ok := <-w.fsWatcher.Events:
				if !ok {
					return
				}

				if event.Op&fsnotify.Write == fsnotify.Write {
					onSave(20)
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

func (w *Watcher) Close() {
	w.fsWatcher.Close()
}