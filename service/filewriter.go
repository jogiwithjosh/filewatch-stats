package service

import (
	"os"
)

type FileWriter interface {
	Write(string, int) error
}

type JSONWriter struct {
	storageFile *os.File
}

func NewFileWriter(storagePath string) (FileWriter, error) {
	_, err := os.Stat(storagePath)
	var file *os.File
	if err != nil {
		file, err = os.Create(storagePath)
		if err != nil {
			return nil, err
		}
	}

	return &JSONWriter{
		storageFile: file,
	}, nil
}

func (jw *JSONWriter) Write(filePath string, noOfBytes int) error {
	//dataRow := fmt.Sprintf(`{"filePath": %s, }`)
	return nil
}
