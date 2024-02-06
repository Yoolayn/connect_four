package main

import "errors"

var (
	ErrParse           = errors.New("failed to parse results")
	ErrArgsReq         = errors.New("this command requires arguments")
	ErrRequest         = errors.New("failed a request")
	ErrUnknown         = errors.New("unknown error has occured")
	ErrEncoding        = errors.New("failed to marshal json")
	ErrParamReq        = errors.New("search parameter is required")
	ErrNewParams       = errors.New("new command only accepts user/game")
	ErrWrongCode       = errors.New("wrong code")
	ErrLoginTaken      = errors.New("login taken")
	ErrNoResponse      = errors.New("failed to get results")
	ErrUserParams      = errors.New("user command only accepts update/delete")
	ErrUserUpdate      = errors.New("user update only accepts login/password")
	ErrWrongParam      = errors.New("wrong search parameter")
	ErrCmdNotFound     = errors.New("command not found")
	ErrNotLoggedIn     = errors.New("user is not logged in")
	ErrUserNotFound    = errors.New("no user with this name exists")
	ErrWrongPassword   = errors.New("wrong password")
	ErrNotImplemented  = errors.New("not yet implemented")
	ErrNotEnoughParams = errors.New("not enough parameters")
)
