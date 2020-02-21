package util

import (
	"errors"
	"fmt"
	"github.com/allegro/bigcache"
	"time"
)

type Cache interface {
	Get(key string) ([]byte, error)
	Set(key string, entry []byte) error
	GetAll() (map[string][]byte, error)
	Delete(key string) error
	FindValue(value string) (string, error)
}

func NewBigCache(eviction time.Duration) (*BigCache, error) {
	var err error
	cache, err := bigcache.NewBigCache(bigcache.DefaultConfig(eviction))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to initialize cahce: %v", err))
	}
	return &BigCache{cache}, nil
}

type BigCache struct {
	*bigcache.BigCache
}

func (c *BigCache) GetAll() (map[string][]byte, error) {
	items := make(map[string][]byte, 0)
	iterator := c.Iterator()
	for iterator.SetNext() {
		current, err := iterator.Value()
		if err == nil {
			items[current.Key()] = current.Value()
		}
	}
	return items, nil
}

func (c *BigCache) FindValue(value string) (string, error) {
	for iter := c.Iterator(); iter.SetNext(); {
		info, err := iter.Value()
		if err != nil {
			continue
		}
		if string(info.Value()) == value {
			return info.Key(), nil
		}
	}
	return "", errors.New("cache value not found")
}



