package main

import (
	"log"
	"os"
	"time"

	hd "download_bot/pkg/handlers"
	rt "download_bot/pkg/router"
	tt "download_bot/pkg/tiktok"

	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v4"
)

func getToken() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	TOKEN := os.Getenv("TOKEN")
	return TOKEN
}

func main() {
	tiktok := tt.NewTiktokClient()

	router := rt.NewRouter(tiktok)

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

	hd.RegisterHandlers(bot, router)

	bot.Start()
}
