package delete

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
		if err != nil || serialNum == "" {
			log.Error("invalid serial number", err)

			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, handler.RespMsg("invalid serial number"))
			return
		}

		err = h.Service.DeleteDevice(serialNum)
		if err != nil {
			log.Error("DeleteDevice", err)

			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, handler.RespMsg(err.Error()))
			return
		}
		w.WriteHeader(http.StatusAccepted)
		log.Info("delete device")
	}
}
