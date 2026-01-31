package main

import (
	"errors"
	"net/http"
	"strings"
)

func check_domain(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	domain := resp.Request.URL.Host
	return domain, nil
}

func check_correct_service(domain string) (string, error) {
	var err error = nil
	split_domain := strings.Split(domain, ".")
	service := split_domain[len(split_domain)-2]
	switch service {
	case "tiktok":
		return service, err
	case "youtube":
		return service, err
	case "instagram":
		return service, err
	default:
		err := errors.New("Unknown service")
		return "", err
	}
}
