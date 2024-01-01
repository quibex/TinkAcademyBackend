package handler

import (
	"device-service/internal/device"
	"device-service/internal/service"
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	"net/http"
)

func RespMsg(msg string) Response {
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

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=Service
type Service interface {
	service.Repository
}

type Handler struct {
	Service Service
	Logger  *slog.Logger
}

func ValidateRequest(w http.ResponseWriter, r *http.Request) (*Request, error) {
	var req Request

	err := render.DecodeJSON(r.Body, &req)
	if errors.Is(err, io.EOF) {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, RespMsg("empty request"))
		return nil, errors.New("request body is empty")
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, RespMsg("failed to decode request body: "+err.Error()))
		return nil, errors.New("failed to decode request body: " + err.Error())
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, RespMsg("invalid request: "+err.Error()))
		return nil, errors.New("invalid request: " + err.Error())
	}

	return &req, nil
}
