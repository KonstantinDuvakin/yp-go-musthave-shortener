package handler

import (
	"net/http"

	"github.com/KonstantinDuvakin/yp-go-musthave-shortener/internal/repository"
)

func GetShortLinkHandler(storage *repository.LinksStorage) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		linkId := r.PathValue("id")
		link := storage.GetLink(linkId)

		if link == "" {
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte("Not Found"))
			return
		}

		rw.Header().Set("Location", link)
		rw.WriteHeader(http.StatusTemporaryRedirect)
	}
}
