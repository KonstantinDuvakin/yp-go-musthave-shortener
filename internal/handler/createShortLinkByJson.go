package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/KonstantinDuvakin/yp-go-musthave-shortener/internal/logger"
	"github.com/KonstantinDuvakin/yp-go-musthave-shortener/internal/model"
	"github.com/KonstantinDuvakin/yp-go-musthave-shortener/internal/repository"
	"go.uber.org/zap"
)

func CreateShortLinkByJson(s *repository.LinksStorage, baseUrl string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var req model.GetLinkShortenReq

		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&req); err != nil && err != io.EOF {
			logger.Log.Error("Error decoding request", zap.Error(err))
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("invalid request body"))
			return
		}

		if req.Url == "" {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("url is required"))
			return
		}

		linkId := s.GetOrCreateLink([]byte(req.Url))

		res := model.GetLinkShortenRes{
			Result: baseUrl + linkId,
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusCreated)

		if err := json.NewEncoder(rw).Encode(res); err != nil {
			logger.Log.Error("Error encoding response", zap.Error(err))
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
