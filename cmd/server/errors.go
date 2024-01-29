package main

import (
	"errors"
	"net/http"
)

var (
	ErrNotAuthorized = errors.New("user not authorized")
	ErrUserNotFound  = errors.New("user not found")
	ErrInternal      = errors.New("broke")
	ErrType          = errors.New("bad type assertion")
	ErrUserTaken     = errors.New("user name is taken")
	ErrPassTooLong   = errors.New("password is too long")
	ErrParsing       = errors.New("failed to parse uuid")
	ErrGameNotFound  = errors.New("game not found")
)

var errs = map[error]int{
	ErrNotAuthorized: http.StatusUnauthorized,
	ErrUserNotFound:  http.StatusNotFound,
	ErrInternal:      http.StatusInternalServerError,
	ErrType:          http.StatusInternalServerError,
	ErrUserTaken:     http.StatusConflict,
	ErrPassTooLong:   http.StatusBadRequest,
	ErrParsing:       http.StatusBadRequest,
	ErrGameNotFound:  http.StatusNotFound,
}

func newErr(err error) (int, interface{}) {
	return errs[err], struct {
		Code int    `json:"code"`
		Err  string `json:"err"`
	}{
		Code: errs[err],
		Err:  err.Error(),
	}
}
