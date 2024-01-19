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
	directory   string
	concurrency int
	watcherr    *watcher.Watcher
	fw          FileWriter
}

func NewFileWatcher(directory string, concurrency int, recursive bool, fw FileWriter) (FileWatcher, error) {
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
		watcherr:    watcherr,
		directory:   directory,
		concurrency: concurrency,
		fw:          fw,
	}, nil
}

func (w *FileWatcherImpl) Watch() error {
	defer w.watcherr.Close()
	writeCh := w.fw.Start()

	go func() {
		for {
			select {
			case event := <-w.watcherr.Event:
				if !event.FileInfo.IsDir() {
					writeCh <- Stat{FilePath: event.Path, ByteSize: event.Size()}
				}
			case err := <-w.watcherr.Error:
				fmt.Println(err)
			case <-w.watcherr.Closed:
				fmt.Println("closed")
				return
			}
		}
	}()

	return w.watcherr.Start(time.Second * 1)
}
