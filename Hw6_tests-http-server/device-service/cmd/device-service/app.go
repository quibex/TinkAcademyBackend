package main

import (
	"context"
	"device-service/internal/config"
	h "device-service/internal/http-server/handler"
	"device-service/internal/http-server/handler/device/create"
	"device-service/internal/http-server/handler/device/delete"
	"device-service/internal/http-server/handler/device/receive"
	"device-service/internal/http-server/handler/device/update"
	"device-service/internal/lib/logger/slogpretty"
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
	log := setupPrettySlog()
	servc := service.New()

	handler := h.Handler{
		Service: servc,
		Logger:  log,
	}
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.URLFormat)

	router.Get("/devices/{serial_num}", receive.New(handler))
	router.Post("/devices", create.New(handler))
	router.Patch("/devices", update.New(handler))
	router.Delete("/devices/{serial_num}", delete.New(handler))

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

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
