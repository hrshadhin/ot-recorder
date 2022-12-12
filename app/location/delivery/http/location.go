package http

import (
	"errors"
	"ot-recorder/app/model"
	"ot-recorder/app/response"
	"ot-recorder/app/validation"

	"github.com/labstack/echo/v4"
)

// LocationHandler represent the http handler for Location
type LocationHandler struct {
	LUseCase model.LocationUsecase
}

func NewUserHandler(e *echo.Echo, us model.LocationUsecase) {
	handler := &LocationHandler{
		LUseCase: us,
	}

	v1 := e.Group("/api/v1")
	v1.POST("/ping", handler.Ping)
	v1.GET("/last-location", handler.LastLocation)
}

func (u *LocationHandler) Ping(c echo.Context) error {
	req := c.Request()

	if len(req.Header.Get("x-limit-u")) == 0 ||
		len(req.Header.Get("x-limit-d")) == 0 ||
		len(req.Header.Get("X-Real-IP")) == 0 {
		msg := "[x-limit-u, x-limit-d, X-Real-IP] any one of these headers are missing!" +
			"\n Set username & device id in app settings. Also set X-Real-IP header in Poxy settings"

		return c.JSON(response.RespondError(response.ErrBadRequest, errors.New(msg)))
	}

	var pingReq PingRequest
	err := c.Bind(&pingReq)
	if err != nil {
		return c.JSON(response.RespondError(response.ErrUnprocessableEntity, err))
	}

	if pingReq.Type == "location" {
		if ok, err := validation.Validate(&pingReq); !ok {
			valErrors, valErr := validation.FormatErrors(err)
			if valErr != nil {
				return c.JSON(response.RespondError(response.ErrBadRequest, valErr))
			}

			return c.JSON(response.RespondValidationError(response.ErrBadRequest, valErrors))
		}

		location := mapLocationRequestToModel(&pingReq, &req.Header)

		ctx := c.Request().Context()
		err = u.LUseCase.Ping(ctx, location)
		if err != nil {
			return c.JSON(response.RespondError(err))
		}
	}

	return c.JSON(response.RespondEmpty())
}

func (u *LocationHandler) LastLocation(c echo.Context) error {
	ctx := c.Request().Context()

	username := c.QueryParam("username")
	location, err := u.LUseCase.LastLocation(ctx, username)
	if err != nil {
		return c.JSON(response.RespondError(err))
	}
	return c.JSON(response.RespondSuccess("request success", location))
}
