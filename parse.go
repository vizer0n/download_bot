package main

import (
	"encoding/json"
	"errors"
	html "golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"log"
)

func get_download_link(url string, service string, tiktok *TiktokHttp) (string, string, error) {
	switch service {
	case "tiktok":
		download_url, video_name, err := parse_tiktok(url, tiktok)
		if err != nil {
			log.Print(err)
			return "", "", err
		}
		return download_url, video_name, err
	default:
		return "", "", errors.New(service + " пока что не поддерживется")
	}

}

func get_video_name_tiktok(itemStruct map[string]any) string {
	id := itemStruct["id"].(string)
	author := itemStruct["author"].(map[string]any)["uniqueId"].(string)
	return author + "__" + id + ".mp4"
}

func parse_tiktok(url string, tiktok *TiktokHttp) (string, string, error) {
	respBody, err := tiktok.TiktokParse(url, false)

	doc, err := html.Parse(respBody)
	if err != nil {
		log.Print(err)
		return "", "", err
	}

	for n := range doc.Descendants() {
		if n.Type == html.ElementNode && n.DataAtom == atom.Script {
			for _, a := range n.Attr {
				if a.Key == "id" && a.Val == "__UNIVERSAL_DATA_FOR_REHYDRATION__" {
					var data_map map[string]any
					data := []byte(n.FirstChild.Data)
					if err := json.Unmarshal(data, &data_map); err != nil {
						log.Print(err)
						return "", "", err
					}
					var playAddr string

					if root, ok := data_map["__DEFAULT_SCOPE__"].(map[string]any); ok {
						if webapp, ok := root["webapp.video-detail"].(map[string]any); ok {
							if itemInfo, ok := webapp["itemInfo"].(map[string]any); ok {
								if itemStruct, ok := itemInfo["itemStruct"].(map[string]any); ok {
									if video, ok := itemStruct["video"].(map[string]any); ok {
										playAddr = video["playAddr"].(string)
										video_name := get_video_name_tiktok(itemStruct)
										return playAddr, video_name, nil
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return "Not Found", "", nil
}
