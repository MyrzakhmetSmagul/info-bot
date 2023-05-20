package db_mysql

import "tg-bot/repository"

func (s storage) CreateCommand(command *repository.Command) error {
	//TODO implement me
	panic("implement me")
}

func (s storage) UpdateCommand(command repository.Command) error {
	//TODO implement me
	panic("implement me")
}

func (s storage) GetCommand(commandID int) (repository.Command, error) {
	//TODO implement me
	panic("implement me")
}

func (s storage) GetAllCommands() ([]repository.Command, error) {
	//TODO implement me
	panic("implement me")
}

func (s storage) DeleteCommand(commandID int) error {
	//TODO implement me
	panic("implement me")
}
