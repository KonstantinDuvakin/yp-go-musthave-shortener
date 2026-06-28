package repository

import (
	"bufio"
	"encoding/json"
	"os"
	"strconv"
	"sync"

	"github.com/KonstantinDuvakin/yp-go-musthave-shortener/internal/helpers"
	"github.com/KonstantinDuvakin/yp-go-musthave-shortener/internal/middlewares/logger"
	"github.com/KonstantinDuvakin/yp-go-musthave-shortener/internal/model"
	"go.uber.org/zap"
)

type LinksStorage struct {
	mu            sync.RWMutex
	Links         map[string]string
	RevertedLinks map[string]string
	filepath      string
}

func NewLinkStorage(filepath string) *LinksStorage {
	return &LinksStorage{
		Links:         make(map[string]string),
		RevertedLinks: make(map[string]string),
		filepath:      filepath,
	}
}

func (ls *LinksStorage) saveLink(key string, value string) {
	if ls.Links == nil {
		ls.Links = make(map[string]string)
	}

	if ls.RevertedLinks == nil {
		ls.RevertedLinks = make(map[string]string)
	}
	ls.Links[key] = value
	ls.RevertedLinks[value] = key
}

func (ls *LinksStorage) GetLink(key string) string {
	ls.mu.RLock()
	defer ls.mu.RUnlock()

	return ls.Links[key]
}

func (ls *LinksStorage) GetOrCreateLink(link []byte) string {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	key, ok := ls.RevertedLinks[string(link)]
	if ok {
		return key
	}

	key = helpers.GenerateShortLink(8)
	ls.saveLink(key, string(link))
	if err := ls.saveLinkToFile(key, string(link)); err != nil {
		logger.Log.Warn("Error saving link to file:", zap.Error(err))
	}
	return key
}

func (ls *LinksStorage) saveLinkToFile(key, value string) error {
	file, err := os.OpenFile(ls.filepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	defer file.Close()

	link := model.Link{Id: strconv.Itoa(len(ls.Links)), ShortUrl: key, OriginalUrl: value}
	enc := json.NewEncoder(file)
	if err = enc.Encode(link); err != nil {
		return err
	}

	return nil
}

func (ls *LinksStorage) RestoreLinkFromFile() error {
	file, err := os.Open(ls.filepath)
	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		data := scanner.Bytes()

		link := model.Link{}

		if err = json.Unmarshal(data, &link); err != nil {
			return err
		}

		ls.saveLink(link.ShortUrl, link.OriginalUrl)
	}

	return nil
}
