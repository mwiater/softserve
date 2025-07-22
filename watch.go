// watch.go
package softserve

import (
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func watchForChanges(root string, hub *reloadHub) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("watcher error: %v", err)
	}
	defer watcher.Close()

	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return watcher.Add(path)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("walk error: %v", err)
	}

	for {
		select {
		case event := <-watcher.Events:
			log.Printf("📁 Changed: %s", event.Name)
			hub.broadcastReload()
		case err := <-watcher.Errors:
			log.Printf("watch error: %v", err)
		}
	}
}
