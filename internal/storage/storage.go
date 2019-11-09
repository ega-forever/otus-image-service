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

func New(count int, storeDir string) *Storage {
	ctx := context.Background()
	lru := NewLRU(count)

	_ = os.RemoveAll(storeDir)
	_ = os.MkdirAll(storeDir, os.ModePerm)

	return &Storage{ctx: ctx, lru: lru, storeDir: storeDir}
}

func (storage *Storage) SaveImageByURL(ctx context.Context, url string, width int, height int, filename string, headers map[string][]string) error {

	ch := make(chan error)

	mutex := sync.Mutex{}

	go func() {
		mutex.Lock()
		_, removedItem := storage.lru.Put(url+"."+strconv.Itoa(width)+"."+strconv.Itoa(height), lruItem{filename: filename, headers: headers})

		if removedItem != nil {
			err := os.Remove(path.Join(storage.storeDir, removedItem.(lruItem).filename))
			ch <- err
		} else {
			ch <- nil
		}

		mutex.Unlock()
	}()

	return <-ch
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

	file, err := os.Create(path.Join(storage.storeDir, filename))
	defer file.Close()

	if err != nil {
		return err
	}

	_, err = io.Copy(file, inStream)

	if err != nil {
		return err
	}

	src, err := imaging.Open(path.Join(storage.storeDir, filename))

	if err != nil {
		return err
	}

	src = imaging.Resize(src, width, height, imaging.Lanczos)
	err = imaging.Save(src, path.Join(storage.storeDir, filename))

	return err
}
