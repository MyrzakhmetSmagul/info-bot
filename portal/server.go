package portal

import (
	file_manager "tg-bot/file-manager"
	"tg-bot/repository"
)

type Portal struct {
	repository  repository.Repository
	fileManager file_manager.FileManager
	basePath    string
}

func NewPortal(repository repository.Repository, fileManager file_manager.FileManager, basePath string) Portal {
	return Portal{
		repository:  repository,
		fileManager: fileManager,
		basePath:    basePath,
	}
}
