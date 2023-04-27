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
