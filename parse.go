package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	html "golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func get_download_link(resp http.Response, service string) (string, error) {
	switch service {
	case "tiktok":
		download_url, err := parse_tiktok(resp)
		if err != nil {
			log.Print(err)
			return "", err
		}
		return download_url, err
	default:
		return "", errors.New(service + " пока что не поддерживется")
	}

}

func parse_tiktok(resp http.Response) (string, error) {
	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Print(err)
		return "", err
	}

	for n := range doc.Descendants() {
		if n.Type == html.ElementNode && n.DataAtom == atom.Script {
			for _, a := range n.Attr {
				if a.Key == "id" && a.Val == "__UNIVERSAL_DATA_FOR_REHYDRATION__" {
					var data_map map[string]any
					data := []byte(n.FirstChild.Data)
					if err := json.Unmarshal(data, &data_map); err != nil {
						log.Print(err)
						return "", err
					}
					var playAddr string

					if root, ok := data_map["__DEFAULT_SCOPE__"].(map[string]any); ok {
						if webapp, ok := root["webapp.video-detail"].(map[string]any); ok {
							if itemInfo, ok := webapp["itemInfo"].(map[string]any); ok {
								if itemStruct, ok := itemInfo["itemStruct"].(map[string]any); ok {
									if video, ok := itemStruct["video"].(map[string]any); ok {
										playAddr = video["playAddr"].(string)
										return playAddr, nil
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return "Not Found", nil
}
