package create_test

import (
	"bytes"
	"device-service/internal/device"
	h "device-service/internal/http-server/handler"
	"device-service/internal/http-server/handler/device/create"
	"device-service/internal/http-server/handler/mocks"
	"device-service/internal/lib/logger/slogdiscard"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestHandlerCreateValidDevice(t *testing.T) {
	d := device.Device{
		SerialNum: "645",
		Model:     "a",
		IP:        "1.6.3.76",
	}

	t.Log(d)

	serviceMock := mocks.NewService(t)
	serviceMock.On("CreateDevice", mock.Anything).Return(nil)

	handler := h.Handler{
		Service: serviceMock,
		Logger:  slogdiscard.NewDiscardLogger(),
	}

	router := chi.NewRouter()

	router.Post("/devices", create.New(handler))

	input, err := json.Marshal(d)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest(http.MethodPost, "/devices", bytes.NewReader(input))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)
}

func TestHandlerCreateAlreadyExistingDevice(t *testing.T) {
	d := device.Device{
		SerialNum: "645",
		Model:     "a",
		IP:        "1.6.3.76",
	}

	serviceMock := mocks.NewService(t)
	serviceMock.On("CreateDevice", d).Return(fmt.Errorf("device already exists: %s", d.SerialNum))

	handler := h.Handler{
		Service: serviceMock,
		Logger:  slogdiscard.NewDiscardLogger(),
	}

	router := chi.NewRouter()

	router.Post("/devices", create.New(handler))

	input, err := json.Marshal(d)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest(http.MethodPost, "/devices", bytes.NewReader(input))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusConflict, rr.Code)

	body := rr.Body.String()

	var resp h.Response

	require.NoError(t, json.Unmarshal([]byte(body), &resp))

	require.Equal(t, fmt.Sprintf("device already exists: %s", d.SerialNum), resp.Message)
}

func FuzzHandlerCreateValidRequest(f *testing.F) {
	f.Fuzz(func(t *testing.T, serialNum uint16, model string, ip1, ip2, ip3, ip4 uint8) {
		d := device.Device{
			SerialNum: strconv.Itoa(int(serialNum)),
			Model:     model + "a",
			IP:        strconv.Itoa(int(ip1)) + "." + strconv.Itoa(int(ip2)) + "." + strconv.Itoa(int(ip3)) + "." + strconv.Itoa(int(ip4)),
		}

		serviceMock := mocks.NewService(t)
		serviceMock.On("CreateDevice", mock.Anything).Return(nil)

		handler := h.Handler{
			Service: serviceMock,
			Logger:  slogdiscard.NewDiscardLogger(),
		}

		router := chi.NewRouter()

		router.Post("/devices", create.New(handler))

		input, err := json.Marshal(d)
		if err != nil {
			t.Error(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/devices", bytes.NewReader(input))
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		require.Equal(t, http.StatusCreated, rr.Code)
	})
}
