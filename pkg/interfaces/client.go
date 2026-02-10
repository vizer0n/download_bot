package interfaces

import "io"

type VideoService interface {
	GetName() string
	NewVideo(string) VideoObject
	Match(string) bool
	GetHTML(string) (io.ReadCloser, error)
	FetchFile(string) (io.ReadCloser, error)
}
