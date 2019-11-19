package storage

import (
	"context"
	"github.com/disintegration/imaging"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"sync"
)

type Storage struct {
	ctx      context.Context
	lru      *LRU
	storeDir string
}

type lruItem struct {
	headers  map[string][]string
	filename string
}

func New(count int, storeDir string) (*Storage, error) {
	ctx := context.Background()
	lru := NewLRU(count)

	err := os.RemoveAll(storeDir)
	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(storeDir, os.ModePerm)

	if err != nil {
		return nil, err
	}

	return &Storage{ctx: ctx, lru: lru, storeDir: storeDir}, nil
}

func (storage *Storage) SaveImageByURL(ctx context.Context, url string, width int, height int, filename string, headers map[string][]string) error {

	mutex := sync.Mutex{}

	mutex.Lock()
	defer mutex.Unlock()

	_, removedItem := storage.lru.Put(url+"."+strconv.Itoa(width)+"."+strconv.Itoa(height), lruItem{filename: filename, headers: headers})

	if removedItem != nil {
		return os.Remove(path.Join(storage.storeDir, removedItem.(lruItem).filename))
	}

	return nil
}

func (storage *Storage) FindCachedImageData(url string, width int, height int) ([]byte, map[string][]string, error) {

	item := storage.lru.Get(url + "." + strconv.Itoa(width) + "." + strconv.Itoa(height))

	if item == nil {
		return nil, nil, nil
	}

	file, err := ioutil.ReadFile(path.Join(storage.storeDir, item.(lruItem).filename))

	if err != nil {
		return nil, nil, err
	}

	return file, item.(lruItem).headers, nil
}

func (storage *Storage) SaveImageData(inStream io.ReadCloser, filename string, width int, height int) error {

	src, err := imaging.Decode(inStream)

	if err != nil {
		return err
	}

	src = imaging.Resize(src, width, height, imaging.Lanczos)
	return imaging.Save(src, path.Join(storage.storeDir, filename))
}
