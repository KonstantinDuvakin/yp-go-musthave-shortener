package handler

import (
	"net/http"

	"github.com/KonstantinDuvakin/yp-go-musthave-shortener/internal/repository"
	"github.com/go-chi/chi/v5"
)

func GetShortLinkHandler(s *repository.LinksStorage) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		linkId := chi.URLParam(r, "id")
		link := s.GetLink(linkId)
		if link == "" {
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte("Not Found"))
			return
		}

		rw.Header().Set("Location", link)
		rw.WriteHeader(http.StatusTemporaryRedirect)
	}
}
