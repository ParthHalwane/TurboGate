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

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Println("🔄 Config file changed, reloading...")
				time.Sleep(300 * time.Millisecond) // debounce
				onChange()
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			log.Printf("watcher error: %v", err)
		}
	}
}
