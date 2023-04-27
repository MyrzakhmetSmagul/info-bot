package file_manager

import (
	"log"
	"os"
	"path"
	"tg-bot/lib/e"
)

type FileManager struct {
	basePath string
}

func New(basePath string) FileManager {
	return FileManager{basePath: basePath}
}

func (f *FileManager) GetFile(filename string) ([]byte, error) {
	file, err := os.ReadFile(path.Join(f.basePath, filename))
	log.Println(path.Join(f.basePath, filename))
	if err != nil {
		return nil, e.Wrap("can't get file", err)
	}

	return file, nil
}

func (f *FileManager) SaveFile(filename string, content []byte) error {
	file, err := os.Create(path.Join(f.basePath, filename))
	if err != nil {
		return e.Wrap("can't create file", err)
	}

	defer file.Close()

	_, err = file.Write(content)
	if err != nil {
		return e.Wrap("can't create file", err)
	}

	return nil
}
