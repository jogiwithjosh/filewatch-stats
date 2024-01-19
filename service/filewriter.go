package service

import (
	"encoding/json"
	"fmt"
	"os"
)

type Stat struct {
	FilePath string
	ByteSize int64
}

type FileWriter interface {
	Start() chan Stat
}

type JSONWriter struct {
	storagePath string
	writeCh     chan Stat
}

func NewFileWriter(storagePath string, concurrency int) (FileWriter, error) {
	return &JSONWriter{
		storagePath: storagePath,
		writeCh:     make(chan Stat, concurrency),
	}, nil
}

func (jw *JSONWriter) write(filePath string, noOfBytes int64) error {
	file, err := os.OpenFile(jw.storagePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	data := make(map[string]int64)
	if err = json.NewDecoder(file).Decode(&data); err != nil && err.Error() != "EOF" {
		return err
	}
	data[filePath] = noOfBytes
	file.Truncate(0)
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")
	return encoder.Encode(data)
}

func (jw *JSONWriter) Start() chan Stat {
	go jw.receiveStat()
	return jw.writeCh
}

func (jw *JSONWriter) receiveStat() {
	for stat := range jw.writeCh {
		if err := jw.write(stat.FilePath, stat.ByteSize); err != nil {
			fmt.Println(err)
		}
	}
}
