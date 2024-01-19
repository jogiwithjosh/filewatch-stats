package service

import (
	"fmt"
	"time"

	"github.com/radovskyb/watcher"
)

type FileWatcher interface {
	Watch() error
}

type FileWatcherImpl struct {
	directory string
	watcherr  *watcher.Watcher
	fw        FileWriter
}

func NewFileWatcher(directory string, recursive bool, fw FileWriter) (FileWatcher, error) {
	watcherr := watcher.New()

	watcherr.FilterOps(watcher.Create, watcher.Write)
	if recursive {
		if err := watcherr.AddRecursive(directory); err != nil {
			return nil, err
		}
	} else {
		if err := watcherr.Add(directory); err != nil {
			return nil, err
		}
	}

	return &FileWatcherImpl{
		watcherr:  watcherr,
		directory: directory,
		fw:        fw,
	}, nil
}

func (w *FileWatcherImpl) Watch() error {
	go func() {
		for {
			select {
			case event := <-w.watcherr.Event:
				if !event.FileInfo.IsDir() {
					fmt.Println(event.FileInfo.Size())
					fmt.Println(w.fw.Write(event.Path, int(event.FileInfo.Size())))
				}
				fmt.Println(event)

			case err := <-w.watcherr.Error:
				fmt.Println(err)
			case <-w.watcherr.Closed:
				fmt.Println("closed")
			default:
				// do nothing
			}
		}
	}()
	w.watcherr.Start(1 * time.Second)
	//w.watcherr.Wait()
	return nil
}
