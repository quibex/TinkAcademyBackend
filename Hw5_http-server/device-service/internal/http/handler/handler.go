package handler

import (
	"device-service/internal/service"
	"device-service/internal/service/device"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	"net/http"
)

const (
	Get = iota
	Post
	Patch
	Delete
)

func respError(msg string) Response {
	return Response{
		Message: msg}
}

type Request struct {
	device.Device
}

type Response struct {
	Message string        `json:"message,omitempty"`
	Device  device.Device `json:"device,omitempty"`
}

type Handler struct {
	Service service.Service
	Logger  *slog.Logger
}

func (h Handler) New(method uint8) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		const op = "handler.New"

		log := h.Logger.With(
			slog.String("op", op),
		)

		serialNum := chi.URLParam(r, "serial_num")
		err := validator.New().Var(serialNum, "omitempty,numeric")
		if err != nil {
			log.Error("invalid serial number", err)

			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, respError("invalid serial number"))
		}

		req := &Request{}
		if serialNum == "" {
			req, err = ValidateRequest(w, r)
			if err != nil {
				log.Error("invalid request", err)
				return
			}
		}

		switch method {
		case Get:
			processGetReq(w, r, h, log, serialNum)
		case Post:
			processPostReq(w, r, h, log, req.Device)
		case Delete:
			processDelReq(w, r, h, log, serialNum)
		case Patch:
			processPatchReq(w, r, h, log, req.Device)
		}

	}
}

func processGetReq(w http.ResponseWriter, r *http.Request, h Handler, log *slog.Logger, serialNum string) {
	d, err := h.Service.GetDevice(serialNum)
	if err != nil {
		log.Error("GetDevice", err)

		w.WriteHeader(http.StatusNotFound)
		render.JSON(w, r, Response{Message: err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, Response{Device: d})

	log.Info("send device's info")
}

func processPostReq(w http.ResponseWriter, r *http.Request, h Handler, log *slog.Logger, reqDevice device.Device) {
	err := h.Service.CreateDevice(reqDevice)
	if err != nil {
		log.Error("CreateDevice", err)

		w.WriteHeader(http.StatusConflict)
		render.JSON(w, r, respError(err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
	log.Info("add device")
}

func processPatchReq(w http.ResponseWriter, r *http.Request, h Handler, log *slog.Logger, reqDevice device.Device) {
	err := h.Service.UpdateDevice(reqDevice)
	if err != nil {
		log.Error("UpdateDevice", err)

		w.WriteHeader(http.StatusConflict)
		render.JSON(w, r, respError(err.Error()))
		return
	}
	w.WriteHeader(http.StatusAccepted)
	log.Info("update device")
}

func processDelReq(w http.ResponseWriter, r *http.Request, h Handler, log *slog.Logger, serialNum string) {
	err := h.Service.DeleteDevice(serialNum)
	if err != nil {
		log.Error("DeleteDevice", err)

		w.WriteHeader(http.StatusConflict)
		render.JSON(w, r, respError(err.Error()))
		return
	}
	w.WriteHeader(http.StatusAccepted)
	log.Info("delete device")
}

func ValidateRequest(w http.ResponseWriter, r *http.Request) (*Request, error) {
	var req Request

	err := render.DecodeJSON(r.Body, &req)
	if errors.Is(err, io.EOF) {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, respError("empty request"))
		return nil, errors.New("request body is empty")
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, respError("failed to decode request"))
		return nil, errors.New("failed to decode request body")
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		validateErr := err.(validator.ValidationErrors)
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, validateErr)
		return nil, errors.New("invalid request")
	}

	return &req, nil
}
