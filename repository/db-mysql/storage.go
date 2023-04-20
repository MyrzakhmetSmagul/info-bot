package db_mysql

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"log"
	"os"
	"path"
	"tg-bot/lib/e"
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

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	storage := &storage{
		db: db,
	}

	if err = storage.createDB(); err != nil {
		log.Fatal(err)
	}

	return storage
}

func (s *storage) createDB() (err error) {
	const errMsg = "can't create db"
	defer func() { err = e.WrapIfErr(errMsg, err) }()

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	script, err := os.ReadFile(path.Clean("./repository/db-mysql/tg-bot_mysql_create.sql"))
	if err != nil {
		return err
	}

	if _, err := s.db.Query(string(script)); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
