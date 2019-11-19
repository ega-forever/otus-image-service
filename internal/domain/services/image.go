package services

import (
	"context"
	"github.com/ega-forever/otus-image-service/internal/domain/interfaces"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ImageService struct {
	imageStorage interfaces.ImageStorage
}

func NewImageService(imageStorage interfaces.ImageStorage) *ImageService {
	return &ImageService{imageStorage: imageStorage}
}

func (es *ImageService) CacheToStorage(ctx context.Context, url string, width int, height int) ([]byte, map[string][]string, error) {

	cachedFile, cachedHeaders, _ := es.imageStorage.FindCachedImageData(url, width, height)

	if cachedFile != nil {
		return cachedFile, cachedHeaders, nil
	}

	filename, headers, err := es.grabAndCacheImage(url, height, width)

	if err != nil {
		return nil, nil, err
	}

	err = es.imageStorage.SaveImageByURL(ctx, url, width, height, filename, headers)

	if err != nil {
		return nil, nil, err
	}

	cachedFile, cachedHeaders, _ = es.imageStorage.FindCachedImageData(url, width, height)

	return cachedFile, cachedHeaders, nil
}

func (es *ImageService) grabAndCacheImage(url string, width int, height int) (string, map[string][]string, error) {

	response, err := http.Get("https://" + url)

	if err != nil {
		response, err = http.Get("http://" + url)

		if err != nil {
			return "", nil, err
		}
	}

	defer response.Body.Close()
	fileParts := strings.Split(url, ".")
	filename := strconv.FormatInt(time.Now().UnixNano(), 10) + "." + fileParts[len(fileParts)-1] //todo add extention
	err = es.imageStorage.SaveImageData(response.Body, filename, width, height)

	if err != nil {
		return "", nil, err
	}

	return filename, response.Header, nil
}
