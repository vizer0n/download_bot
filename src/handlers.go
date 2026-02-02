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
		log.Print("Создан новый объект ", video.GetServiceName(), " video")

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

		log.Print("Начата загрузка видео...")
		err = video.Download()
		if err != nil {
			log.Print(err)
			return c.Send("Не удалось скачать видео, отправьте ссылку ещё раз")
		}
		log.Print("Видео загружено, идёт отправка на сервер...")

		upload_video := &tele.Video{
			File:      tele.FromDisk(video.GetVideoPath()),
			Duration:  int(video.GetDuration()),
			Width:     int(video.GetWidth()),
			Height:    int(video.GetHeight()),
			Streaming: true,
			FileName:  video.GetVideoName(),
		}
		c.Notify(tele.UploadingVideo)
		err = c.Send(upload_video)
		if err != nil {
			log.Print(err)
			return c.Send("Приоизошла ошибка при отправке видео, попробуйте ещё раз")
		}
		log.Print("Видео отправлено! Идёт удаление файла")

		return video.Delete()
	})
}
