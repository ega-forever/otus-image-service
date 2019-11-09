package interfaces

import (
	"context"
	"io"
)

// todo add more actions
type ImageStorage interface {
	SaveImageByURL(ctx context.Context, url string, width int, height int, filename string, headers map[string][]string) error
	FindCachedImageData(url string, width int, height int) ([]byte, map[string][]string, error)
	SaveImageData(inStream io.ReadCloser, filename string, width int, height int) error
}
