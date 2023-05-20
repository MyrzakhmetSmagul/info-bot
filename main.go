package main

import (
	"encoding/json"
	"github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"tg-bot/portal"
)

type config struct {
	Token        string        `json:"token"`
	BasePath     string        `json:"base_path"`
	MysqlConfig  mysql.Config  `json:"mysql_config"`
	RedisOptions redis.Options `json:"redis_options"`
}

func main() {
	portal.Run()
	//portal.Start()
	//cnf := getConfig()
	//
	//eventsProcessor := telegram.New(
	//	tg.New(cnf.Token),
	//	dbMysql.New(cnf.MysqlConfig),
	//	file_manager.New(cnf.BasePath),
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
