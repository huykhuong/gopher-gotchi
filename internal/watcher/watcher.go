package watcher

import (
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

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

func (w *Watcher) Start(rootPath string, onSave func(lines int)) {
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
					// Pick up any newly created subdirectories.
					_ = w.registerDirs(rootPath)
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

func (w *Watcher) registerDirs(rootPath string) error {
	return filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return nil
		}

		if info.Name()[0] == '.' || info.Name() == "node_modules" {
			return filepath.SkipDir
		}

		// Only watch dirs under a folder listed in folderToRegister.
		for p := path; ; {
			base := filepath.Base(p)
			for _, name := range folderToRegister {
				if base == name {
					return w.fsWatcher.Add(path)
				}
			}
			parent := filepath.Dir(p)
			if parent == p {
				break
			}
			p = parent
		}
		return nil
	})
}

func (w *Watcher) Close() {
	w.fsWatcher.Close()
}