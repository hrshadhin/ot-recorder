package http

import (
	"encoding/json"

	"github.com/labstack/echo/v4"
)

type MyJSONSerializer struct {
	echo.DefaultJSONSerializer
}

func (d MyJSONSerializer) Serialize(c echo.Context, i interface{}, indent string) error {
	enc := json.NewEncoder(c.Response())
	enc.SetEscapeHTML(false)

	if indent != "" {
		enc.SetIndent("", indent)
	}

	return enc.Encode(i)
}
