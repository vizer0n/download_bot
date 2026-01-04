package main

import (
	"log"

	tele "gopkg.in/telebot.v4"
)

func register_handlers(bot *tele.Bot, tiktok *TiktokHttp) {
	bot.Handle("/start", func(c tele.Context) error {
		return c.Send("Привет! Отправь ссылку на видео tiktok, youtube, instagram и я скачаю его")
	})

	bot.Handle(tele.OnText, func(c tele.Context) error {
		// Логирование сообщений
		user := c.Chat().Username
		msg_text := c.Message().Text
		log.Print("User " + user + " send text: " + msg_text)

		// Проверка доменного имени
		domain, err := check_domain(msg_text)
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
		log.Print(message)

		// Получение ссылки для скачивания видео
		download_url, video_name, err := get_download_link(msg_text, service, tiktok)
		if err != nil {
			return c.Send(err.Error())
		}
		log.Print("get_download_link() сработана успешно")

		video_path, err := download_video(download_url, video_name, service, tiktok)
		if err != nil {
			return c.Send(err.Error())
		}
		return c.Send(&tele.Video{File: tele.FromDisk(video_path)})

	})
}
