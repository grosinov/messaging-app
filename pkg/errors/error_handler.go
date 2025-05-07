package errors

import (
	"errors"
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (e ErrorResponse) Error() string {
	return e.Message
}

func BadRequestError(msg string) ErrorResponse {
	return ErrorResponse{
		Status:  http.StatusBadRequest,
		Message: msg,
	}
}

func NotFoundError(msg string) ErrorResponse {
	return ErrorResponse{
		Status:  http.StatusNotFound,
		Message: msg,
	}
}

func InternalServerError(msg string, err error) ErrorResponse {
	return ErrorResponse{
		Status:  http.StatusInternalServerError,
		Message: fmt.Sprintf("%s: %s", msg, err),
	}
}

func HandleError(w http.ResponseWriter, err error) {
	var er ErrorResponse
	if errors.As(err, &er) {
		http.Error(w, er.Error(), er.Status)
	} else {
		http.Error(w, fmt.Sprintf("Internal Server Error: %v", err), http.StatusInternalServerError)
	}
}
