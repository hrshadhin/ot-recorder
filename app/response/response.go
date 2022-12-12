package response

import "net/http"

type Response struct {
	Message string      `json:"message,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func RespondSuccess(msg string, data interface{}) (int, Response) {
	return http.StatusOK, Response{
		Message: msg,
		Data:    data,
	}
}

func RespondEmpty() (int, interface{}) {
	return http.StatusOK, []string{}
}
