package receive_test

import (
	"bytes"
	"device-service/internal/device"
	h "device-service/internal/http-server/handler"
	"device-service/internal/http-server/handler/device/receive"
	"device-service/internal/http-server/handler/mocks"
	"device-service/internal/lib/logger/slogdiscard"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestHandlerReceiveValidRequest(t *testing.T) {
	serviceMock := mocks.NewService(t)
	serviceMock.On("GetDevice", mock.Anything).Return(device.Device{}, nil)

	handler := h.Handler{
		Service: serviceMock,
		Logger:  slogdiscard.NewDiscardLogger(),
	}

	router := chi.NewRouter()
	router.Use(middleware.URLFormat)

	router.Get("/devices/{serial_num}", receive.New(handler))

	var input []byte

	req, err := http.NewRequest(http.MethodGet, "/devices/"+"464532", bytes.NewReader(input))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestHandlerReceiveNotExistingDevice(t *testing.T) {
	serialNum := "464532"

	serviceMock := mocks.NewService(t)
	serviceMock.On("GetDevice", serialNum).Return(device.Device{}, fmt.Errorf("device not found: %s", serialNum))

	handler := h.Handler{
		Service: serviceMock,
		Logger:  slogdiscard.NewDiscardLogger(),
	}

	router := chi.NewRouter()
	router.Use(middleware.URLFormat)

	router.Get("/devices/{serial_num}", receive.New(handler))

	var input []byte

	req, err := http.NewRequest(http.MethodGet, "/devices/"+serialNum, bytes.NewReader(input))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)

	body := rr.Body.String()

	var resp h.Response

	require.NoError(t, json.Unmarshal([]byte(body), &resp))

	require.Equal(t, fmt.Sprintf("device not found: %s", serialNum), resp.Message)
}

func FuzzHandlerCreateValidRequest(f *testing.F) {
	f.Fuzz(func(t *testing.T, serialNum uint16) {
		serviceMock := mocks.NewService(t)
		serviceMock.On("DeleteDevice", mock.Anything).Return(device.Device{}, nil)

		handler := h.Handler{
			Service: serviceMock,
			Logger:  slogdiscard.NewDiscardLogger(),
		}

		router := chi.NewRouter()
		router.Use(middleware.URLFormat)

		router.Get("/devices/{serial_num}", receive.New(handler))

		var input []byte

		req, err := http.NewRequest(http.MethodGet, "/devices/"+strconv.Itoa(int(serialNum)), bytes.NewReader(input))
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
	})
}
