package main

import (
	"log"
	"os"

	"github.com/tucnak/telebot"
)

func init() {
	file, err := os.OpenFile("karbarban.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Open log file", err)
	}
	log.SetOutput(file)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	cfg, err := ReadConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	err = ConnectToDatabase(cfg.Database.Dialect, cfg.Database.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	bot, err := telebot.NewBot(cfg.Telegram.Token)
	if err != nil {
		log.Fatal(err)
	}

	go telegram(bot)
	crawler(bot)
}
