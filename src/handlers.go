package main

import (
	"log"

	tele "gopkg.in/telebot.v4"
)

func register_handlers(bot *tele.Bot, router *Router) {
	bot.Handle("/start", func(c tele.Context) error {
		return c.Send("Привет! Отправь ссылку на видео tiktok, youtube, instagram и я скачаю его")
	})

	bot.Handle(tele.OnText, func(c tele.Context) error {
		// Логирование сообщений
		user := c.Chat().Username
		msg_text := c.Message().Text
		log.Print("User " + user + " send text: " + msg_text)

		// Определение сервиса
		service, err := router.Resolve(msg_text)
		if err != nil {
			log.Print(err)
			return c.Send(err.Error())
		}

		// Инициализация видео объекта
		video := service.NewVideo(msg_text)

		// Получение HTML страницы в виде *html.Node
		err = video.GetHTML()
		if err != nil {
			log.Print(err)
			return c.Send("Произошла ошибка, попробуйте позже")
		}

		// Получение ссылки для скачивания видео и его название
		err = video.GetVideoInfo()
		if err != nil {
			log.Print(err)
			return c.Send("Произошла ошибка, попробуйте позже")
		}

		err = video.Download()
		if err != nil {
			log.Print()
		}

		c.Notify(tele.UploadingVideo)
		return c.Send(&tele.Video{File: tele.FromDisk(video.GetVideoPath())})

	})
}
