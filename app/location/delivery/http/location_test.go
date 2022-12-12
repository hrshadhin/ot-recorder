package http_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	lHttp "ot-recorder/app/location/delivery/http"
	"ot-recorder/app/model"
	"ot-recorder/app/model/mocks"
	"ot-recorder/app/response"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var BaseURLV1 = "/api/v1"

func TestPing(t *testing.T) {
	mockUsecase := new(mocks.LocationUsecase)
	mockUsecase.On("Ping", mock.Anything, mock.AnythingOfType("*model.Location")).Return(nil)

	pingReq := lHttp.PingRequest{
		Type:  "location",
		Tst:   time.Now().Unix(),
		Acc:   13,
		Alt:   -42,
		Batt:  40,
		Bs:    1,
		Lat:   23.0000000,
		Lon:   90.0000000,
		M:     1,
		Tid:   "p1",
		Vac:   1,
		Vel:   0,
		Bssid: "c0:00:00:00:00:00",
		Ssid:  "dev-test",
	}

	endPoint := BaseURLV1 + "/ping"

	t.Run("success", func(t *testing.T) {
		tempReq := pingReq
		j, err := json.Marshal(tempReq)
		assert.NoError(t, err)
		c, rec := buildEchoRequest(t, endPoint, echo.POST, strings.NewReader(string(j)), true)

		handler := lHttp.LocationHandler{
			LUseCase: mockUsecase,
		}
		err = handler.Ping(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("headers missing", func(t *testing.T) {
		tempReq := pingReq

		j, err := json.Marshal(tempReq)
		assert.NoError(t, err)
		c, rec := buildEchoRequest(t, endPoint, echo.POST, strings.NewReader(string(j)), false)

		handler := lHttp.LocationHandler{
			LUseCase: mockUsecase,
		}
		err = handler.Ping(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestPingNonPingRequest(t *testing.T) {
	mockUsecase := new(mocks.LocationUsecase)
	mockUsecase.On("Ping", mock.Anything, mock.AnythingOfType("*model.Location")).Return(nil)

	pingReq := lHttp.PingRequest{
		Type:  "location",
		Tst:   time.Now().Unix(),
		Acc:   13,
		Alt:   -42,
		Batt:  40,
		Bs:    1,
		Lat:   23.0000000,
		Lon:   90.0000000,
		M:     1,
		Tid:   "p1",
		Vac:   1,
		Vel:   0,
		Bssid: "c0:00:00:00:00:00",
		Ssid:  "dev-test",
	}

	endPoint := BaseURLV1 + "/ping"

	t.Run("by pass non ping request", func(t *testing.T) {
		tempReq := pingReq
		tempReq.Type = "waypoint"

		j, err := json.Marshal(tempReq)
		assert.NoError(t, err)
		c, rec := buildEchoRequest(t, endPoint, echo.POST, strings.NewReader(string(j)), true)

		handler := lHttp.LocationHandler{
			LUseCase: mockUsecase,
		}
		err = handler.Ping(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestLastLocationSuccess(t *testing.T) {
	endPoint := BaseURLV1 + "/last-location"
	mockLoc := model.LocationDetails{
		Username:         "dev",
		Device:           "dev-test",
		DateTime:         time.Now().Format("2006-01-02 15:04:05"),
		Accuracy:         14,
		Altitude:         -45,
		BatteryLevel:     "80%",
		BatteryStatus:    "Unplugged",
		Latitude:         23.0000000,
		Longitude:        90.0000000,
		Mode:             "Moving",
		VerticalAccuracy: 1,
		Velocity:         0,
		IPAddress:        "127.0.0.1",
	}
	mockUsecase := new(mocks.LocationUsecase)
	mockUsecase.On("LastLocation", mock.Anything, mock.AnythingOfType("string")).Return(&mockLoc, nil)

	handler := lHttp.LocationHandler{
		LUseCase: mockUsecase,
	}

	ctx, res := buildEchoRequest(t, endPoint+"?username=dev", echo.GET, nil, false)
	handle := handler.LastLocation

	assert.NoError(t, handle(ctx))
	assert.Equal(t, http.StatusOK, res.Code)

	var r response.Response
	assert.NoError(t, json.Unmarshal(res.Body.Bytes(), &r))
	assert.Equal(t, "request success", r.Message)

	resultsMap := r.Data.(map[string]interface{})
	assert.Equal(t, mockLoc.Device, resultsMap["device"])

	mockUsecase.AssertExpectations(t)
}

func TestLastLocationNotFound(t *testing.T) {
	endPoint := BaseURLV1 + "/last-location"

	mockUsecase := new(mocks.LocationUsecase)
	mockUsecase.On("LastLocation", mock.Anything, mock.AnythingOfType("string")).Return(nil, response.ErrNotFound)

	handler := lHttp.LocationHandler{
		LUseCase: mockUsecase,
	}

	ctx, res := buildEchoRequest(t, endPoint+"?username=foobar", echo.GET, nil, false)

	handle := handler.LastLocation
	assert.NoError(t, handle(ctx))
	assert.Equal(t, http.StatusNotFound, res.Code)

	mockUsecase.AssertExpectations(t)
}

func buildEchoRequest(
	t *testing.T,
	path,
	method string,
	payload io.Reader,
	auth bool,
) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	var req *http.Request
	var err error

	if method == echo.POST {
		req, err = http.NewRequest(method, path, payload)
	} else {
		req, err = http.NewRequest(method, path, nil)
	}
	assert.NoError(t, err)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("X-Real-IP", "127.0.0.1")

	if auth {
		req.Header.Set("x-limit-u", "dev")
		req.Header.Set("x-limit-d", "phoneAndroid")
	}

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath(path)

	return c, rec
}
