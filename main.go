package main

import (
	"encoding/json"
	"github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	file_manager "tg-bot/file-manager"
	"tg-bot/portal"
	db_mysql "tg-bot/repository/db-mysql"
)

type config struct {
	Token        string        `json:"token"`
	BasePath     string        `json:"base_path"`
	MysqlConfig  mysql.Config  `json:"mysql_config"`
	RedisOptions redis.Options `json:"redis_options"`
}

func main() {
	cnf := getConfig()
	site := portal.NewPortal(db_mysql.New(cnf.MysqlConfig), file_manager.New(cnf.BasePath))
	site.Run()
	//
	//eventsProcessor := telegram.New(
	//	tg.New(cnf.Token),

	//	redis.NewClient(&cnf.RedisOptions),
	//)
	//
	//log.Println("...service started")
	//
	//consumer := consumer.New(eventsProcessor, eventsProcessor, 100)
	//
	//if err := consumer.Start(); err != nil {
	//	log.Fatal(err)
	//}
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
