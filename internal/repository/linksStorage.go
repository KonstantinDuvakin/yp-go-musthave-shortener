package repository

import (
	"sync"

	"github.com/KonstantinDuvakin/yp-go-musthave-shortener/internal/helpers"
)

type LinksStorage struct {
	mu    sync.RWMutex
	Links map[string]string
}

func NewLinkStorage() *LinksStorage {
	return &LinksStorage{
		Links: make(map[string]string),
	}
}

func (ls *LinksStorage) saveLink(key string, value string) {
	if ls.Links == nil {
		ls.Links = make(map[string]string)
	}
	ls.Links[key] = value
}

func (ls *LinksStorage) GetLink(key string) string {
	ls.mu.RLock()
	defer ls.mu.RUnlock()

	return ls.Links[key]
}

func (ls *LinksStorage) GetOrCreateLink(link []byte) string {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	for key, value := range ls.Links {
		if value == string(link) {
			return key
		}
	}

	key := helpers.GenerateShortLink(8)
	ls.saveLink(key, string(link))
	return key
}
