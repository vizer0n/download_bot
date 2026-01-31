package main

import (
	"encoding/json"
	"errors"
	"io"
	"os"

	html "golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type VideoObject interface {
	GetHTML() error
	GetVideoInfo() error
	GetVideoPath() string
	Download() error
}

type TiktokVideo struct {
	service    TiktokService
	base_url   string
	html       *html.Node
	video_link string
	video_name string
	video_path string
}

func NewTiktokVideo(service *TiktokService, url string) VideoObject {
	return &TiktokVideo{service: *service, base_url: url}
}

func (tv *TiktokVideo) GetHTML() error {
	req, _ := tv.service.newRequest("GET", tv.base_url)

	resp, err := tv.service.Client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	tv.html, err = html.Parse(resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func (tv *TiktokVideo) GetVideoInfo() error {
	for n := range tv.html.Descendants() {
		if n.Type == html.ElementNode && n.DataAtom == atom.Script {
			for _, a := range n.Attr {
				if a.Key == "id" && a.Val == "__UNIVERSAL_DATA_FOR_REHYDRATION__" {
					var data_map map[string]any
					data := []byte(n.FirstChild.Data)
					if err := json.Unmarshal(data, &data_map); err != nil {
						return err
					}

					if root, ok := data_map["__DEFAULT_SCOPE__"].(map[string]any); ok {
						if webapp, ok := root["webapp.video-detail"].(map[string]any); ok {
							if itemInfo, ok := webapp["itemInfo"].(map[string]any); ok {
								if itemStruct, ok := itemInfo["itemStruct"].(map[string]any); ok {
									id := itemStruct["id"].(string)
									author := itemStruct["author"].(map[string]any)["uniqueId"].(string)
									tv.video_name = author + "__" + id + ".mp4"
									if video, ok := itemStruct["video"].(map[string]any); ok {
										if playAddr, ok := video["playAddr"].(string); ok {
											tv.video_link = playAddr
											return nil
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return errors.New("playAddr не был извлечён")
}

func (tv *TiktokVideo) Download() error {
	os.MkdirAll("downloads/tiktok", os.ModePerm)
	video_path := "./downloads/tiktok/" + tv.video_name
	resp_body, err := tv.service.FetchVideo(tv.video_link)
	if err != nil {
		return err
	}

	video_file, err := os.Create(video_path)
	if err != nil {
		return err
	}
	defer video_file.Close()

	if _, err := io.Copy(video_file, resp_body); err != nil {
		return err
	}

	tv.video_path = video_path
	return nil
}

func (tv *TiktokVideo) GetVideoPath() string {
	return tv.video_path
}
