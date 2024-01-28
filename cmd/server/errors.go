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
)

var errs = map[error]int{
	ErrNotAuthorized: http.StatusUnauthorized,
	ErrUserNotFound:  http.StatusNotFound,
	ErrInternal:      http.StatusInternalServerError,
	ErrType:          http.StatusInternalServerError,
	ErrUserTaken:     http.StatusConflict,
	ErrPassTooLong:   http.StatusBadRequest,
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
