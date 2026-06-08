package main

import (
	"net/http"

	"github.com/KonstantinDuvakin/yp-go-musthave-shortener/internal/handler"
	"github.com/KonstantinDuvakin/yp-go-musthave-shortener/internal/repository"
)

func main() {
	storage := &repository.LinksStorage{
		Links: make(map[string]string),
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{id}", handler.GetShortLinkHandler(storage))
	mux.HandleFunc("POST /", handler.CreateShortLinkHandler(storage))

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
