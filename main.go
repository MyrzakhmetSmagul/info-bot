package main

import (
	"encoding/json"
	"github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	tgbot_api "tg-bot/client/telegram/tgbot-api"
	event_consumer "tg-bot/consumer/event-consumer"
	"tg-bot/events/telegram"
	file_manager "tg-bot/file-manager"
	"tg-bot/portal"
	db_mysql "tg-bot/repository/db-mysql"
)

type config struct {
	Token        string        `json:"token"`
	BasePath     string        `json:"base_path"`
	MysqlConfig  mysql.Config  `json:"mysql_config"`
	RedisOptions redis.Options `json:"redis_options"`
	Port         string        `json:"port"`
}

func main() {
	cnf := getConfig()
	db := db_mysql.New(cnf.MysqlConfig)
	fileManager := file_manager.New(cnf.BasePath)
	if len(os.Args) != 1 && (os.Args[1] == "-site" || os.Args[1] == "--site") {
		site := portal.NewPortal(db, fileManager, cnf.BasePath)
		go site.Run(cnf.Port)
	}

	eventsProcessor := telegram.New(tgbot_api.New(cnf.Token), db, fileManager)

	log.Println("...service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, 100)

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
