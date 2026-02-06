package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v4"
)

func getToken() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	TOKEN := os.Getenv("TOKEN")
	if TOKEN == "" {
		log.Fatal("token is empty")
	}
	return TOKEN
}

func main() {
	tiktok := NewTiktokClient()

	router := NewRouter(tiktok)

	TOKEN := getToken()

	pref := tele.Settings{
		Token:  TOKEN,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	register_handlers(bot, router)

	bot.Start()
}
