package rest

import (
	"fmt"
	"net/http"
)

type HttpError struct {
	Err        error `json:"-"`
	Code       int   `json:"code"`
	HttpStatus int   `json:"-"`
}

func (e HttpError) Error() string {
	return fmt.Sprintf("%+v", e)
}

func NewHttpError(err error, code int, httpStatus int) HttpError {
	return HttpError{err, code, httpStatus}
}

func CodeError(err error, code int) HttpError {
	return HttpError{err, code, http.StatusInternalServerError}
}

func StatusError(err error, httpStatus int) HttpError {
	return HttpError{err, CodeUnknown, httpStatus}
}
