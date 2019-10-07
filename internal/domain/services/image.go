package services

import (
	"context"
	"github.com/ega-forever/otus-image-service/internal/domain/interfaces"
	"github.com/ega-forever/otus-image-service/internal/domain/models"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

type ImageService struct {
	imageStorage interfaces.ImageStorage
	storeDir     string
}

func NewImageService(imageStorage interfaces.ImageStorage, storeDir string) *ImageService {
	return &ImageService{imageStorage: imageStorage, storeDir: storeDir}
}

func (es *ImageService) SaveToStorage(ctx context.Context, url string, timestamp int64) (*models.Image, error) {

	name, err := es.grabImage(url)

	if err != nil {
		return nil, err
	}

	savedEvent, err := es.imageStorage.SaveImage(ctx, &models.Image{Url: url})

	if err != nil {
		return nil, err
	}

	return savedEvent, nil
}

func (es *ImageService) grabImage(url string) (string, error) {

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
		file, err := os.Create(path.Join(es.storeDir, filename))

		if err != nil {
			ch <- message{err: err}
		}
		defer file.Close()

		_, err = io.Copy(file, response.Body)

		if err != nil {
			ch <- message{err: err}
		}

		ch <- message{name: filename}
	}()

	m := <-ch
	return m.name, m.err
}
