package delete_test

import (
	"bytes"
	"device-service/internal/device"
	h "device-service/internal/http-server/handler"
	"device-service/internal/http-server/handler/device/delete"
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

func TestHandlerDeleteValidRequest(t *testing.T) {
	serviceMock := mocks.NewService(t)
	serviceMock.On("DeleteDevice", mock.Anything).Return(nil)

	handler := h.Handler{
		Service: serviceMock,
		Logger:  slogdiscard.NewDiscardLogger(),
	}

	router := chi.NewRouter()
	router.Use(middleware.URLFormat)

	router.Delete("/devices/{serial_num}", delete.New(handler))

	var input []byte

	req, err := http.NewRequest(http.MethodDelete, "/devices/"+"3634", bytes.NewReader(input))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusAccepted, rr.Code)
}

func TestHandlerDeleteNotExistingDevice(t *testing.T) {
	serialNum := "464532"

	serviceMock := mocks.NewService(t)
	serviceMock.On("DeleteDevice", serialNum).Return(fmt.Errorf("device not found: %s", serialNum))

	handler := h.Handler{
		Service: serviceMock,
		Logger:  slogdiscard.NewDiscardLogger(),
	}

	router := chi.NewRouter()
	router.Use(middleware.URLFormat)

	router.Delete("/devices/{serial_num}", delete.New(handler))

	var input []byte

	req, err := http.NewRequest(http.MethodDelete, "/devices/"+serialNum, bytes.NewReader(input))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)

	body := rr.Body.String()

	var resp h.Response

	require.NoError(t, json.Unmarshal([]byte(body), &resp))

	require.Equal(t, fmt.Sprintf("device not found: %s", serialNum), resp.Message)
}

func FuzzHandlerDeleteValidRequest(f *testing.F) {
	f.Fuzz(func(t *testing.T, serialNum uint16) {
		serviceMock := mocks.NewService(t)
		serviceMock.On("DeleteDevice", mock.Anything).Return(device.Device{}, nil)

		handler := h.Handler{
			Service: serviceMock,
			Logger:  slogdiscard.NewDiscardLogger(),
		}

		router := chi.NewRouter()
		router.Use(middleware.URLFormat)

		router.Delete("/devices/{serial_num}", delete.New(handler))

		var input []byte

		req, err := http.NewRequest(http.MethodDelete, "/devices/"+strconv.Itoa(int(serialNum)), bytes.NewReader(input))
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		require.Equal(t, http.StatusAccepted, rr.Code)
	})
}
