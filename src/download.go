package main

import (
	"errors"
	"io"
	"log"
	"os"
)

func download_video(url string, video_name string, service string, tiktok *TiktokHttp) (string, error) {
	switch service {
	case "tiktok":
		os.MkdirAll("downloads/tiktok", os.ModePerm)
		video_path := "./downloads/tiktok/" + video_name
		resp_body, err := tiktok.TiktokParse(url, true)
		if err != nil {
			return "", err
		}

		video_file, err := os.Create(video_path)
		if err != nil {
			return "", err
		}
		defer video_file.Close()

		written, err := io.Copy(video_file, resp_body)
		if err != nil {
			return "", err
		}
		log.Print(written)

		return video_path, nil
	default:
		return "", errors.New("Неподдерживаемый сервис")
	}
}
