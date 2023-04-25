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

var StartState int

func New(cnf mysql.Config) repository.Repository {
	db, err := sql.Open("mysql", cnf.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	storage := &storage{
		db: db,
	}

	storage.getNeedState()

	return storage
}

func (s *storage) getNeedState() {
	query := `SELECT id, name FROM state WHERE name='/start'`
	err := s.db.QueryRow(query).Scan(&repository.StartState.ID, repository.StartState.Name)
	if err != nil {
		log.Fatal(err)
	}
	query = `SELECT id, name FROM state WHERE name='/language'`
	err = s.db.QueryRow(query).Scan(&repository.LanguageState.ID, repository.LanguageState.Name)
	if err != nil {
		log.Fatal(err)
	}
	query = `SELECT id, name FROM state WHERE name='/question'`
	err = s.db.QueryRow(query).Scan(&repository.QuestionState.ID, repository.QuestionState.Name)
	if err != nil {
		log.Fatal(err)
	}
}
