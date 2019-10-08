package storage

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"path"
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

func (storage *Storage) SaveImageByURL(ctx context.Context, url string, filename string, headers map[string][]string) error {

	ch := make(chan error)

	mutex := sync.Mutex{}

	go func() {
		mutex.Lock()
		_, removedItem := storage.lru.Put(url, lruItem{filename: filename, headers: headers})

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

func (storage *Storage) FindCachedImageData(url string) ([]byte, map[string][]string, error) {

	item := storage.lru.Get(url)

	if item == nil {
		return nil, nil, nil
	}

	file, err := ioutil.ReadFile(path.Join(storage.storeDir, item.(lruItem).filename))

	if err != nil {
		return nil, nil, err
	}

	return file, item.(lruItem).headers, nil
}

func (storage *Storage) SaveImageData(inStream io.ReadCloser, filename string) error {

	file, err := os.Create(path.Join(storage.storeDir, filename))
	defer file.Close()

	if err != nil {
		return err
	}

	_, err = io.Copy(file, inStream)
	return err
}
