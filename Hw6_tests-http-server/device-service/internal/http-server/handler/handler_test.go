package handler_test

import (
	"bytes"
	"device-service/internal/device"
	h "device-service/internal/http-server/handler"
	"device-service/internal/http-server/handler/device/create"
	"device-service/internal/http-server/handler/device/receive"
	"device-service/internal/http-server/handler/mocks"
	"device-service/internal/lib/logger/slogdiscard"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestSuite struct {
	suite.Suite
	h http.Handler
}

func (s *TestSuite) SetupTest() {
	serviceMock := mocks.NewService(s.T())

	handler := h.Handler{
		Service: serviceMock,
		Logger:  slogdiscard.NewDiscardLogger(),
	}

	router := chi.NewRouter()
	router.Use(middleware.URLFormat)

	router.Get("/devices/{serial_num}", receive.New(handler))
	router.Post("/devices", create.New(handler))

	s.h = router
}

func (s *TestSuite) TestHandlerValidateRequest() {
	cases := []struct {
		name        string
		d           device.Device
		respMessage string
	}{
		{
			name: "With Invalid IP",
			d: device.Device{
				SerialNum: "6242",
				Model:     "15 pro",
				IP:        "34.75.75;",
			},
			respMessage: "invalid request: Key: 'Request.Device.IP' Error:Field validation for 'IP' failed on the 'ip4_addr' tag",
		},
		{
			name: "With Invalid serialNum",
			d: device.Device{
				SerialNum: "drive",
				Model:     "15 pro",
				IP:        "34.75.75.75",
			},
			respMessage: "invalid request: Key: 'Request.Device.SerialNum' Error:Field validation for 'SerialNum' failed on the 'numeric' tag",
		},
		{
			name: "Without Model field",
			d: device.Device{
				SerialNum: "85356",
				IP:        "34.75.75.75",
			},
			respMessage: "invalid request: Key: 'Request.Device.Model' Error:Field validation for 'Model' failed on the 'required' tag",
		},
		{
			name: "Without SerialNum field",
			d: device.Device{
				Model: "15 pro",
				IP:    "34.75.75.75",
			},
			respMessage: "invalid request: Key: 'Request.Device.SerialNum' Error:Field validation for 'SerialNum' failed on the 'required' tag",
		},
		{
			name: "Without IP field",
			d: device.Device{
				SerialNum: "85356",
				Model:     "15 pro",
			},
			respMessage: "invalid request: Key: 'Request.Device.IP' Error:Field validation for 'IP' failed on the 'required' tag",
		},
		{
			name:        "Empty request",
			respMessage: "empty request",
		},
	}

	for _, tc := range cases {
		tc := tc

		s.T().Run(tc.name, func(t *testing.T) {
			t.Parallel()

			handler := s.h

			input, err := json.Marshal(tc.d)
			if (tc.d == device.Device{}) {
				input = []byte{}
			}
			if err != nil {
				t.Error(err)
			}

			req, err := http.NewRequest(http.MethodPost, "/devices", bytes.NewReader(input))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, http.StatusBadRequest, rr.Code)

			body := rr.Body.String()

			var resp h.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respMessage, resp.Message)

		})
	}
}

func (s *TestSuite) TestHandlerValidateSerialNum() {
	cases := []struct {
		name       string
		serialNum  string
		statusCode int
	}{
		{
			name:       "Invalid1 SerialNum",
			statusCode: http.StatusBadRequest,
			serialNum:  "goland",
		},
		{
			name:       "Invalid1 SerialNum",
			statusCode: http.StatusBadRequest,
			serialNum:  "356&*",
		},
		{
			name:       "Invalid1 SerialNum",
			statusCode: http.StatusBadRequest,
			serialNum:  "{]4463",
		},
	}

	for _, tc := range cases {
		tc := tc

		s.T().Run(tc.name, func(t *testing.T) {
			t.Parallel()

			handler := s.h

			req, err := http.NewRequest(http.MethodGet, "/devices/"+tc.serialNum, bytes.NewReader([]byte{}))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, tc.statusCode, rr.Code)

		})
	}
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
