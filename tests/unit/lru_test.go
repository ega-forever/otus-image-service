package lru_test

import (
	"github.com/ega-forever/otus-image-service/internal/storage"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestLRU(t *testing.T) {

	lru := storage.NewLRU(5)

	for i := 0; i < 5; i++ {
		lru.Put(i, i)
	}

	val := lru.Get(1)
	assert.Equal(t, val, 1)

	lru.Put(10, 1)

	val = lru.Get(0)
	assert.Equal(t, val, nil)

	lru.Put(11, 1)

	val = lru.Get(1)
	assert.Equal(t, val, 1)

	val = lru.Get(2)
	assert.Equal(t, val, nil)

}
