package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KonstantinDuvakin/yp-go-musthave-shortener/internal/config"
	"github.com/KonstantinDuvakin/yp-go-musthave-shortener/internal/handler"
	"github.com/KonstantinDuvakin/yp-go-musthave-shortener/internal/repository"
	"github.com/go-chi/chi/v5"
)

func main() {
	storage := repository.NewLinkStorage()
	c := config.NewConfig()

	r := chi.NewRouter()

	r.MethodNotAllowed(func(rw http.ResponseWriter, r *http.Request) {
		http.Error(rw, "Method Not Allowed", http.StatusMethodNotAllowed)
	})

	r.Route("/", func(r chi.Router) {
		r.Post("/", handler.CreateShortLinkHandler(storage, c.BaseUrl))
		r.Get("/{id}", handler.GetShortLinkHandler(storage))
	})

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	server := &http.Server{
		Addr:    c.ServerAddr,
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = server.Shutdown(shutdownCtx)
}
