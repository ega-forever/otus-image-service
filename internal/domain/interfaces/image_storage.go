package interfaces

import (
	"context"
	"io"
)

// todo add more actions
type ImageStorage interface {
	SaveImageByURL(ctx context.Context, url string, filename string, headers map[string][]string) error
	FindCachedImageData(url string) ([]byte, map[string][]string, error)
	SaveImageData(inStream io.ReadCloser, filename string) error
}
