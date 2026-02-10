package interfaces

type VideoObject interface {
	GetServiceName() string
	GetHTML() error
	GetVideoInfo() error
	GetVideoPath() string
	GetThumbnailPath() string
	GetSoundPath() string
	GetSoundPerformer() string
	GetSoundTitle() string
	GetSoundName() string
	GetDuration() float64
	GetSoundDuration() float64
	GetWidth() float64
	GetHeight() float64
	GetVideoName() string
	GetVideoLink() string
	DownloadAll() error
	Delete() error
}
