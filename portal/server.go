package portal

import (
	file_manager "tg-bot/file-manager"
	"tg-bot/repository"
)

type Portal struct {
	repository  repository.Repository
	fileManager file_manager.FileManager
}

func NewPortal(repository repository.Repository, fileManager file_manager.FileManager) Portal {
	return Portal{
		repository:  repository,
		fileManager: fileManager,
	}
}
