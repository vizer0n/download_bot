package main

import (
	"io"
	"net/http"
	"net/http/cookiejar"
	"time"
)

type TiktokHttp struct {
	Client *http.Client
}

func NewTiktokClient() *TiktokHttp {
	jar, _ := cookiejar.New(nil)

	client := &http.Client{
		Jar:     jar,
		Timeout: 30 * time.Second,
	}

	return &TiktokHttp{
		Client: client,
	}
}

func (t *TiktokHttp) NewRequest(method, url string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9,en-US;q=0.8")

	return req, nil
}

func (t *TiktokHttp) TiktokParse(url string, mp4 bool) (io.ReadCloser, error) {
	req, _ := t.NewRequest("GET", url)

	if mp4 {
		req.Header.Set("Referer", "https://www.tiktok.com/")
	}

	resp, err := t.Client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}
