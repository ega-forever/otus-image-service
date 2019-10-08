package interfaces

import (
	"context"
	"io"
)

// todo add more actions
type ImageStorage interface {
	SaveImageByURL(ctx context.Context, url string, filename string) error
	FindCachedImageData(url string) ([]byte, error)
	SaveImageData(inStream io.ReadCloser, filename string) error
}
