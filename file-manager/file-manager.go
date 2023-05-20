package file_manager

import (
	"fmt"
	"log"
	"os"
	"path"
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
		return nil, fmt.Errorf("get file was failed: %w", err)
	}

	return file, nil
}

func (f *FileManager) SaveFile(filename string, content []byte) error {
	file, err := os.Create(path.Join(f.basePath, filename))
	if err != nil {
		return fmt.Errorf("create file was failed: %w", err)
	}

	defer file.Close()

	_, err = file.Write(content)
	if err != nil {
		return fmt.Errorf("create file was failed: %w", err)
	}

	return nil
}
