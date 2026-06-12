package main

import (
	"net/http"

	"github.com/KonstantinDuvakin/yp-go-musthave-shortener/internal/handler"
	"github.com/KonstantinDuvakin/yp-go-musthave-shortener/internal/repository"
	"github.com/go-chi/chi/v5"
)

func main() {
	storage := repository.NewLinkStorage()

	r := chi.NewRouter()

	r.MethodNotAllowed(func(rw http.ResponseWriter, r *http.Request) {
		http.Error(rw, "Method Not Allowed", http.StatusMethodNotAllowed)
	})

	r.Route("/", func(r chi.Router) {
		r.Post("/", handler.CreateShortLinkHandler(storage))
		r.Get("/{id}", handler.GetShortLinkHandler(storage))
	})

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
