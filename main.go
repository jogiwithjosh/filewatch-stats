package main

import (
	"flag"
	"fmt"
	"innowhyte/service"
)

func main() {
	watchPath := flag.String("watchPath", "", "used to watch on the provided path")
	storagePath := flag.String("storagePath", "", "used to store JSON files")
	flag.Parse()

	fmt.Println(*watchPath, *storagePath)
	fw, err := service.NewFileWriter(*storagePath)
	if err != nil {
		panic(err)
	}

	watcher, err := service.NewFileWatcher(*watchPath, false, fw)
	if err != nil {
		panic(err)
	}
	watcher.Watch()
}
