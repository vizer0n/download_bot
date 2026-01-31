package main

import (
	"errors"
	"io"
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

		if _, err := io.Copy(video_file, resp_body); err != nil {
			return "", err
		}

		return video_path, nil
	default:
		return "", errors.New("Неподдерживаемый сервис")
	}
}
