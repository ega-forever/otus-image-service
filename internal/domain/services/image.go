package services

import (
	"context"
	"github.com/ega-forever/otus-image-service/internal/domain/interfaces"
	"net/http"
	"strconv"
	"time"
)

type ImageService struct {
	imageStorage interfaces.ImageStorage
}

func NewImageService(imageStorage interfaces.ImageStorage) *ImageService {
	return &ImageService{imageStorage: imageStorage}
}

func (es *ImageService) CacheToStorage(ctx context.Context, url string) ([]byte, error) {

	cachedFile, _ := es.imageStorage.FindCachedImageData(url)

	if cachedFile != nil {
		return cachedFile, nil
	}

	filename, err := es.grabAndCacheImage(url)

	if err != nil {
		return nil, err
	}

	err = es.imageStorage.SaveImageByURL(ctx, url, filename)

	if err != nil {
		return nil, err
	}

	cachedFile, _ = es.imageStorage.FindCachedImageData(url)

	return cachedFile, nil
}

func (es *ImageService) grabAndCacheImage(url string) (string, error) {

	type message struct {
		err  error
		name string
	}

	ch := make(chan message)

	go func() {
		response, err := http.Get(url)
		defer response.Body.Close()

		if err != nil {
			ch <- message{err: err}
			return
		}

		filename := strconv.FormatInt(time.Now().UnixNano(), 10)
		err = es.imageStorage.SaveImageData(response.Body, filename)

		if err != nil {
			ch <- message{err: err}
		}
		ch <- message{name: filename}
	}()

	m := <-ch
	return m.name, m.err
}

func (es *ImageService) LRU() {

}
