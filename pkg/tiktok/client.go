package tiktok

import (
	"download_bot/pkg/interfaces"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
)

type TiktokService struct {
	Client *http.Client
	Name   string
}

func NewTiktokClient() *TiktokService {
	jar, _ := cookiejar.New(nil)

	client := &http.Client{
		Jar:     jar,
		Timeout: 30 * time.Second,
	}

	return &TiktokService{
		Client: client,
		Name:   "tiktok",
	}
}

func (t *TiktokService) GetName() string {
	return t.Name
}

func (t *TiktokService) NewVideo(url string) interfaces.VideoObject {
	return NewTiktokVideo(t, url)
}

func (t *TiktokService) Match(url string) bool {
	return strings.Contains(url, "tiktok.com")
}

func (t *TiktokService) newRequest(method, url string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9,en-US;q=0.8")

	return req, nil
}

func (t *TiktokService) GetHTML(url string) (io.ReadCloser, error) {
	req, _ := t.newRequest("GET", url)

	resp, err := t.Client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (t *TiktokService) FetchFile(url string) (io.ReadCloser, error) {
	req, _ := t.newRequest("GET", url)
	req.Header.Set("Referer", "https://www.tiktok.com/")

	resp, err := t.Client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}
