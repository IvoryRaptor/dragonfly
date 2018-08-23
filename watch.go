package dragonfly

import (
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"log"
	"strings"
)

type FileWatchService struct {
	file       string
	done       chan bool
	FileChange func(data []byte) error
}

func (l *FileWatchService) Config(kernel IKernel, config map[interface{}]interface{}) error {
	l.file = config["file"].(string)
	return nil
}

func (l *FileWatchService) LoadFile() error {
	data, err := ioutil.ReadFile(l.file)
	if err != nil {
		return err
	}
	return l.FileChange(data)
}

func (l *FileWatchService) Start() error {
	err := l.LoadFile()
	if err != nil {
		log.Fatal(err)
	}
	l.done = make(chan bool)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		defer watcher.Close()
		for {
			select {
			case event := <-watcher.Events:
				//log.Println("event:", event.Name)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
				}
				if event.Op.String() == "CREATE" && strings.LastIndex(l.file, event.Name) > 0 {
					log.Println("File Change: ", event.Name)
					err = l.LoadFile()
					if err != nil {
						log.Fatal(err)
					}
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			case <-l.done:
				return
			}
		}
	}()
	index := strings.LastIndex(l.file, "/")
	watcher.Add(l.file[:index])
	return nil
}

func (l *FileWatchService) Stop() {
	l.done <- true
}
