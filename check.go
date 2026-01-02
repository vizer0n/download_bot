package main

import (
	"errors"
	"net/http"
	"strings"
)

func check_domain(url string) (string, *http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", resp, err
	}
	domain := resp.Request.URL.Host
	return domain, resp, nil
}

func check_correct_service(domain string) (string, string, error) {
	var err error = nil
	split_domain := strings.Split(domain, ".")
	service := split_domain[len(split_domain)-2]
	switch service {
	case "tiktok":
		msg := "Загрузка с тик тока в разработке (прямо сейчас)"
		return msg, service, err
	case "youtube":
		msg := "Загрузка с ютуба пока недоступна"
		return msg, service, err
	case "instagram":
		msg := "Загрузка с instagram недоступна"
		return msg, service, err
	default:
		err := errors.New("Unknown service")
		return "", "", err
	}
}
