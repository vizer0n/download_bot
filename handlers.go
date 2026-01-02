package main

import (
	"log"

	tele "gopkg.in/telebot.v4"
)

func register_handlers(bot *tele.Bot) {
	bot.Handle("/start", func(c tele.Context) error {
		return c.Send("Привет! Отправь ссылку на видео tiktok, youtube, instagram и я скачаю его")
	})

	bot.Handle(tele.OnText, func(c tele.Context) error {
		// Логирование сообщений
		user := c.Chat().Username
		msg_text := c.Message().Text
		log.Print("User " + user + " send text: " + msg_text)

		// Проверка доменного имени
		domain, resp, err := check_domain(msg_text)
		if err != nil {
			log.Print(err)
			return c.Send("Введите корректную ссылку")
		}

		// Проверка сервиса на корректность (Youtube, Tiktok, Instagram)
		message, service, err := check_correct_service(domain)
		if err != nil {
			log.Print(err)
			return c.Send("Неизвестный сервис. Отправьте ссылку на Tiktok, Youtube, Instagram")
		}

		// Получение ссылки для скачивания видео
		download_url, err := get_download_link(*resp, service)
		if err != nil {
			log.Print(err)
			return c.Send(err.Error())
		}

		c.Send(download_url)

		return c.Send(message)

	})
}
