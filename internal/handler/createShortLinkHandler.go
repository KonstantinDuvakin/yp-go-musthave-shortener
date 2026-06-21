package handler

import (
	"io"
	"net/http"

	"github.com/KonstantinDuvakin/yp-go-musthave-shortener/internal/repository"
)

func CreateShortLinkHandler(s *repository.LinksStorage, host string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)

		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(err.Error()))
			return
		}

		if len(body) == 0 {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("body is required"))
			return
		}

		linkID := s.GetOrCreateLink(body)

		rw.Header().Add("Content-Type", "text/plain")
		rw.WriteHeader(http.StatusCreated)
		rw.Write([]byte(host + "/" + linkID))
	}
}
