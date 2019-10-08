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

func New(count int, storeDir string) *Storage {
	ctx := context.Background()
	lru := NewLRU(count)

	_ = os.RemoveAll(storeDir)
	_ = os.MkdirAll(storeDir, os.ModePerm)

	return &Storage{ctx: ctx, lru: lru, storeDir: storeDir}
}

func (storage *Storage) SaveImageByURL(ctx context.Context, url string, filename string) error {

	ch := make(chan error)

	mutex := sync.Mutex{}

	go func() {
		mutex.Lock()
		_, removedFilename := storage.lru.Put(url, filename)

		if removedFilename != "" {
			err := os.Remove(path.Join(storage.storeDir, removedFilename))
			ch <- err
		} else {
			ch <- nil
		}

		mutex.Unlock()
	}()

	return <-ch
}

func (storage *Storage) FindCachedImageData(url string) ([]byte, error) {

	filename := storage.lru.Get(url)
	file, err := ioutil.ReadFile(path.Join(storage.storeDir, filename))

	if err != nil {
		return nil, err
	}

	return file, nil
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
