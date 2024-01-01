package main

import (
	"context"
	"device-service/internal/config"
	h "device-service/internal/http/handler"
	"device-service/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.MustLoad()
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	servc := service.New()

	handler := h.Handler{
		Service: servc,
		Logger:  log,
	}
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.URLFormat)

	router.Get("/devices/{serial_num}", handler.New(h.Get))
	router.Post("/devices", handler.New(h.Post))
	router.Patch("/devices", handler.New(h.Patch))
	router.Delete("/devices/{serial_num}", handler.New(h.Delete))

	log.Info("starting server", slog.String("address", cfg.Address))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server")
			close(done)
		}
	}()

	log.Info("server started")

	<-done
	log.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", err)
		return
	}

	log.Info("server stopped")
}
