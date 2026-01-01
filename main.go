package main

import (
	"log"
	"time"

	"os"

	tele "gopkg.in/telebot.v4"

	"github.com/joho/godotenv"
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
	bot.Handle(tele.OnText, func(c tele.Context) error {
		user := c.Chat().Username
		msg_text := c.Message().Text
		log_message := "User " + user + " send: " + msg_text
		log.Print(log_message)
		return c.Send(msg_text)
	})

	bot.Start()
}
