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

func (es *ImageService) CacheToStorage(ctx context.Context, url string) ([]byte, map[string][]string, error) {

	cachedFile, cachedHeaders, _ := es.imageStorage.FindCachedImageData(url)

	if cachedFile != nil {
		return cachedFile, cachedHeaders, nil
	}

	filename, headers, err := es.grabAndCacheImage(url)

	if err != nil {
		return nil, nil, err
	}

	err = es.imageStorage.SaveImageByURL(ctx, url, filename, headers)

	if err != nil {
		return nil, nil, err
	}

	cachedFile, cachedHeaders, _ = es.imageStorage.FindCachedImageData(url)

	return cachedFile, cachedHeaders, nil
}

func (es *ImageService) grabAndCacheImage(url string) (string, map[string][]string, error) {

	type message struct {
		err     error
		name    string
		headers map[string][]string
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

		ch <- message{name: filename, headers: response.Header}
	}()

	m := <-ch
	return m.name, m.headers, m.err
}
