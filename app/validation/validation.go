package validation

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

const two = 2

func Validate(r interface{}) (bool, error) {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", two)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	err := validate.Struct(r)
	if err != nil {
		return false, err
	}

	return true, nil
}

func FormatErrors(err error) (map[string]interface{}, error) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := make(map[string]interface{}, len(ve))

		for _, fe := range ve {
			fieldName := fe.Field()
			out[fieldName] = msgForField(fe)
		}

		return out, nil
	}

	return nil, err
}

func msgForField(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	case "min":
		m := fmt.Sprintf("Must be at least %v", fe.Param())
		if fe.Type().String() == "string" {
			return m + " characters"
		}

		return m
	case "max":
		m := fmt.Sprintf("Not more than %v", fe.Param())
		if fe.Type().String() == "string" {
			return m + " characters"
		}

		return m
	}

	return "unknown error"
}
