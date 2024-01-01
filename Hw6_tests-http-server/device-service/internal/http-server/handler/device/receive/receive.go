package receive

import (
	"device-service/internal/http-server/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
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

		serialNum := chi.URLParam(r, "serial_num")
		err := validator.New().Var(serialNum, "numeric")
		if (err != nil) || (serialNum == "") {
			log.Error("invalid serial number", err)

			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, handler.RespMsg("invalid serial number"))
			return
		}

		d, err := h.Service.GetDevice(serialNum)
		if err != nil {
			log.Error("GetDevice", err)

			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, handler.Response{Message: err.Error()})
			return
		}
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, handler.Response{Device: d})

		log.Info("send device's info")
	}
}
