package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"innowhyte/service"
	"os"
	"path/filepath"
	"time"
)

func main() {
	watchPath := flag.String("watchPath", "watch", "used to watch on the provided path")
	storagePath := flag.String("storagePath", "storage/test.json", "used to store JSON files")
	concurrency := flag.Int("conc", 1, "max no of file processors")
	startClient := flag.Bool("auto-test", false, "to start file ops client ")
	flag.Parse()

	fmt.Println(*watchPath, *storagePath)
	fw, err := service.NewFileWriter(*storagePath, *concurrency)
	if err != nil {
		panic(err)
	}

	watcher, err := service.NewFileWatcher(*watchPath, *concurrency, false, fw)
	if err != nil {
		panic(err)
	}

	if *startClient {
		go startFileOps(filepath.Join(*watchPath, "test.json"), *concurrency)
		go startFileOps(filepath.Join(*watchPath, "test1.json"), *concurrency)
		go startFileOps(filepath.Join(*watchPath, "test2.json"), *concurrency)
	}

	watcher.Watch()
}

// test client
func startFileOps(fp string, concurrency int) {
	time.Sleep(1 * time.Second)
	counter := 1
	go func() {

		c := 0
		for {
			if c%20 == 0 {
				time.Sleep(1 * time.Second)
			}
			os.Create(filepath.Join("watch", fmt.Sprint("counter", c, ".json")))
			c++
		}
	}()
	for {
		go func() {
			file, err := os.OpenFile(fp, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer file.Close()
			data := make(map[string]int64)
			if err = json.NewDecoder(file).Decode(&data); err != nil && err.Error() != "EOF" {
				fmt.Println(err)
				return
			}
			data[fp] = int64(counter) * 2
			json.NewEncoder(file).Encode(data)
		}()

		counter++
	}

}
