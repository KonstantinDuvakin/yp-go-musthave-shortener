package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KonstantinDuvakin/yp-go-musthave-shortener/internal/config"
	"github.com/KonstantinDuvakin/yp-go-musthave-shortener/internal/handler"
	"github.com/KonstantinDuvakin/yp-go-musthave-shortener/internal/logger"
	"github.com/KonstantinDuvakin/yp-go-musthave-shortener/internal/repository"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func main() {
	storage := repository.NewLinkStorage()
	c := config.NewConfig()

	err := logger.Initialize(c.LogLevel)
	if err != nil {
		logger.Log.Warn("Failed to initialize logger", zap.Error(err))
	}

	r := chi.NewRouter()

	r.MethodNotAllowed(func(rw http.ResponseWriter, r *http.Request) {
		http.Error(rw, "Method Not Allowed", http.StatusMethodNotAllowed)
	})

	r.Route("/", func(r chi.Router) {
		r.Post("/", logger.RequestLoggerWrapper(handler.CreateShortLinkHandler(storage, c.BaseUrl)))
		r.Get("/{id}", logger.RequestLoggerWrapper(handler.GetShortLinkHandler(storage)))
		r.Post("/api/shorten", logger.RequestLoggerWrapper(handler.CreateShortLinkByJson(storage, c.BaseUrl)))
	})

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	server := &http.Server{
		Addr:    c.ServerAddr,
		Handler: r,
	}

	go func() {
		logger.Log.Info("Running server", zap.String("address", c.ServerAddr))
		if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Error("server error: ", zap.Error(err))
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = server.Shutdown(shutdownCtx)
}
