package db_mysql

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"log"
	"tg-bot/repository"
)

type storage struct {
	db *sql.DB
}

func (s *storage) AddTransition(transition *repository.Transition) error {
	query := `INSERT INTO transition (name, from_state, to_state, message_id) VALUES (?, ?, ?, ?)`
	// Insert transition
	_, err := s.db.Exec(query, transition.Name, transition.FromState.ID, transition.ToState.ID, transition.Message)
	if err != nil {
		return err
	}

	return nil

}

func (s *storage) GetTransitionsFromState(stateID int) ([]*repository.Transition, error) {
	//TODO implement me
	panic("implement me")
}

func (s *storage) GetTransitionsToState(stateID int) ([]*repository.Transition, error) {
	//TODO implement me
	panic("implement me")
}

func (s *storage) GetAllTransitions() ([]*repository.Transition, error) {
	//TODO implement me
	panic("implement me")
}

func New(cnf mysql.Config) repository.Repository {
	db, err := sql.Open("mysql", cnf.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	return &storage{
		db: db,
	}
}
