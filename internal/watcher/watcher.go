package watcher

import (
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
)

func WatchConfig(path string, onChange func()) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	err = watcher.Add(path)
	if err != nil {
		return err
	}

	lastChange := time.Now()

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
				if time.Since(lastChange) > 1*time.Second {
					lastChange = time.Now()
					log.Println("🔁 Detected config change")
					onChange()
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			log.Println("watcher error:", err)
		}
	}
}
