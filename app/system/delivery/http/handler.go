package http

import (
	"net/http"
	"ot-recorder/app/response"
	"ot-recorder/app/system/usecase"

	"github.com/labstack/echo/v4"
)

// SystemHandler  represent the httphandler for system
type SystemHandler struct {
	Usecase usecase.SystemUsecase
}

// NewSystemHandler will initialize the system related endpoints
func NewSystemHandler(e *echo.Echo, us usecase.SystemUsecase) {
	handler := &SystemHandler{
		Usecase: us,
	}

	e.GET("/", handler.Root)
	e.GET("/health", handler.Health)
}

// Root will let you know, whoami
func (sh *SystemHandler) Root(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "OwnTracks Recorder API!"})
}

// Health will let you know the heart beats
func (sh *SystemHandler) Health(c echo.Context) error {
	err := sh.Usecase.GetHealth()
	if err != nil {
		return c.JSON(response.RespondError(err))
	}

	return c.JSON(http.StatusOK, &response.Response{Message: "I'm healthy :)"})
}
