package create

import (
	"device-service/internal/http-server/handler"
	"github.com/go-chi/render"
	"io"
	"log/slog"
	"net/http"
)

func New(h handler.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		const op = "handler.New"

		log := h.Logger.With(
			slog.String("op", op),
		)

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Error("failed to close body", err)
			}
		}(r.Body)

		req, err := handler.ValidateRequest(w, r)
		if err != nil {
			log.Error("invalid request", err)
			return
		}

		err = h.Service.CreateDevice(req.Device)
		if err != nil {
			log.Error("CreateDevice", err)

			w.WriteHeader(http.StatusConflict)
			render.JSON(w, r, handler.RespMsg(err.Error()))
			return
		}
		w.WriteHeader(http.StatusCreated)
		log.Info("create device")
	}
}
