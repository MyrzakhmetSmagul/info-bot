package db_mysql

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
	"tg-bot/repository"
)

type storage struct {
	db *sql.DB
}

func New(cnf mysql.Config) repository.Repository {
	db, err := sql.Open("mysql", cnf.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	s := storage{
		db: db,
	}

	err = s.getDefault()
	if err != nil {
		log.Fatal(err)
	}

	return &s
}

func (s storage) getDefault() error {
	query := `SELECT id, name FROM state WHERE name='/start'`
	err := s.db.QueryRow(query).Scan(&repository.DefaultState.ID, &repository.DefaultState.Name)
	if err != nil {
		return fmt.Errorf("db_mysql.getDefault: %w", err)
	}
	msg, err := s.GetMessage("/start", repository.Ru)
	if err != nil {
		return fmt.Errorf("db_mysql.getDefault: %w", err)
	}
	repository.StartMessageGroup, err = s.GetMessageGroup(msg.ID, repository.Ru)
	if err != nil {
		return fmt.Errorf("db_mysql.getDefault: %w", err)
	}
	repository.StartMessageGroup.KzMsg, err = s.GetMessageByID(repository.StartMessageGroup.KzMsg.ID)
	if err != nil {
		return fmt.Errorf("db_mysql.getDefault: %w", err)
	}
	repository.StartMessageGroup.EnMsg, err = s.GetMessageByID(repository.StartMessageGroup.RuMsg.ID)
	if err != nil {
		return fmt.Errorf("db_mysql.getDefault: %w", err)
	}
	repository.StartMessageGroup.EnMsg, err = s.GetMessageByID(repository.StartMessageGroup.EnMsg.ID)
	if err != nil {
		return fmt.Errorf("db_mysql.getDefault: %w", err)
	}
	return nil
}
