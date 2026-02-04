package tiktok

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"time"

	"download_bot/pkg/interfaces"

	html "golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type TiktokVideo struct {
	width           float64
	height          float64
	duration        float64
	service         TiktokService
	base_url        string
	html            *html.Node
	video_link      string
	video_name      string
	video_path      string
	thumbnail_link  string
	thumbnail_name  string
	thumbnail_path  string
	sound_link      string
	sound_name      string
	sound_performer string
	sound_title     string
	sound_path      string
	sound_duration  float64
}

func NewTiktokVideo(service *TiktokService, url string) interfaces.VideoObject {
	return &TiktokVideo{service: *service, base_url: url}
}

func (tv *TiktokVideo) GetServiceName() string {
	return tv.service.Name
}

func (tv *TiktokVideo) GetHTML() error {
	req, _ := tv.service.newRequest("GET", tv.base_url)

	startTime := time.Now().Second()
	resp, err := tv.service.Client.Do(req)
	if err != nil {
		return err
	}
	endTime := time.Now().Second()
	log.Print("HTML файл загружен за ", endTime-startTime, " секунд")

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

					file, err := os.Create("test.json")
					if err != nil {
						return err
					}
					_, err = file.Write(data)
					if err != nil {
						return err
					}
					if err := json.Unmarshal(data, &data_map); err != nil {
						return err
					}

					if root, ok := data_map["__DEFAULT_SCOPE__"].(map[string]any); ok {
						if webapp, ok := root["webapp.video-detail"].(map[string]any); ok {
							if itemInfo, ok := webapp["itemInfo"].(map[string]any); ok {
								if itemStruct, ok := itemInfo["itemStruct"].(map[string]any); ok {
									id := itemStruct["id"].(string)
									author := itemStruct["author"].(map[string]any)["uniqueId"].(string)
									name := author + "__" + id
									tv.video_name = name + ".mp4"
									tv.thumbnail_name = name + ".jpg"
									tv.sound_link = itemStruct["music"].(map[string]any)["playUrl"].(string)
									tv.sound_name = name + ".mp3"
									tv.sound_performer = itemStruct["music"].(map[string]any)["authorName"].(string)
									tv.sound_title = itemStruct["music"].(map[string]any)["title"].(string)
									tv.sound_duration = itemStruct["music"].(map[string]any)["duration"].(float64)
									if video, ok := itemStruct["video"].(map[string]any); ok {
										tv.duration = video["duration"].(float64)
										tv.width = video["width"].(float64)
										tv.height = video["height"].(float64)
										tv.thumbnail_link = video["originCover"].(string)
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

func (tv *TiktokVideo) GetVideoLink() string {
	return tv.video_link
}

func (tv *TiktokVideo) DownloadAll() error {
	os.MkdirAll("downloads/tiktok/videos", os.ModePerm)
	os.MkdirAll("downloads/tiktok/thumbnails", os.ModePerm)
	os.MkdirAll("downloads/tiktok/sounds", os.ModePerm)

	tv.video_path = "./downloads/tiktok/videos/" + tv.video_name
	tv.thumbnail_path = "downloads/tiktok/thumbnails/" + tv.thumbnail_name
	tv.sound_path = "downloads/tiktok/sounds/" + tv.sound_name

	err := tv.downloadSome(tv.video_link, tv.video_path)
	if err != nil {
		return err
	}

	err = tv.downloadSome(tv.thumbnail_link, tv.thumbnail_path)
	if err != nil {
		return err
	}

	err = tv.downloadSome(tv.sound_link, tv.sound_path)
	if err != nil {
		return err
	}

	return nil
}

func (tv *TiktokVideo) downloadSome(link string, path string) error {
	resp_body, err := tv.service.FetchFile(link)
	if err != nil {
		return err
	}
	defer resp_body.Close()

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := io.Copy(file, resp_body); err != nil {
		return err
	}
	return nil
}

func (tv *TiktokVideo) Delete() error {
	err := os.Remove(tv.video_path)
	if err != nil {
		return err
	}

	err = os.Remove(tv.thumbnail_path)
	if err != nil {
		return err
	}

	err = os.Remove(tv.sound_path)
	if err != nil {
		return err
	}
	return nil
}

func (tv *TiktokVideo) GetDuration() float64 {
	return tv.duration
}

func (tv *TiktokVideo) GetSoundDuration() float64 {
	return tv.sound_duration
}

func (tv *TiktokVideo) GetSoundPerformer() string {
	return tv.sound_performer
}

func (tv *TiktokVideo) GetSoundTitle() string {
	return tv.sound_title
}

func (tv *TiktokVideo) GetSoundName() string {
	return tv.sound_name
}

func (tv *TiktokVideo) GetWidth() float64 {
	return tv.width
}

func (tv *TiktokVideo) GetHeight() float64 {
	return tv.height
}

func (tv *TiktokVideo) GetVideoName() string {
	return tv.video_name
}

func (tv *TiktokVideo) GetVideoPath() string {
	return tv.video_path
}

func (tv *TiktokVideo) GetThumbnailPath() string {
	return tv.thumbnail_path
}

func (tv *TiktokVideo) GetSoundPath() string {
	return tv.sound_path
}
