package main

import (
	"encoding/json"
	"github.com/go-sql-driver/mysql"
	"log"
	"os"
	tg "tg-bot/client/telegram/tgbot-api"
	consumer "tg-bot/consumer/event-consumer"
	"tg-bot/events/telegram"
	dbMysql "tg-bot/repository/db-mysql"
	file_manager "tg-bot/repository/file-manager"
)

type config struct {
	Token       string       `json:"token"`
	BasePath    string       `json:"base_path"`
	MysqlConfig mysql.Config `json:"mysql_config"`
}

func main() {
	cnf := getConfig()

	eventsProcessor := telegram.New(
		tg.New(cnf.Token),
		dbMysql.New(cnf.MysqlConfig),
		file_manager.New(cnf.BasePath),
	)

	log.Println("...service started")

	consumer := consumer.New(eventsProcessor, eventsProcessor, 100)

	if err := consumer.Start(); err != nil {
		log.Fatal(err)
	}
}

func getConfig() config {
	data, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatal("file ", err)
	}

	var cnf config
	err = json.Unmarshal(data, &cnf)
	if err != nil {
		log.Fatal("json", err)
	}

	return cnf
}
