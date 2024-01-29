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
	ErrForbidden     = errors.New("forbidden!")
	ErrNotInGame     = errors.New("user is not in the game")
	ErrOutOfBounds   = errors.New("selected row would be out of bounds")
	ErrFieldTaken    = errors.New("field is taken")
	ErrUpdateFailed  = errors.New("update in database failed")
	ErrAddFailed     = errors.New("failed adding to database")
	ErrGetAllFailed  = errors.New("failed getting resources from database")
	ErrGameFull      = errors.New("game is full!")
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
	ErrForbidden:     http.StatusForbidden,
	ErrNotInGame:     http.StatusBadRequest,
	ErrOutOfBounds:   http.StatusBadRequest,
	ErrFieldTaken:    http.StatusBadRequest,
	ErrUpdateFailed:  http.StatusInternalServerError,
	ErrAddFailed:     http.StatusInternalServerError,
	ErrGetAllFailed:  http.StatusInternalServerError,
	ErrGameFull:      http.StatusBadRequest,
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
