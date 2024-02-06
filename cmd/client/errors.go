package main

import "errors"

var (
	ErrParse           = errors.New("failed to parse results")
	ErrArgsReq         = errors.New("this command requires arguments")
	ErrUnknown         = errors.New("unknown error has occured")
	ErrRequest         = errors.New("failed a request")
	ErrParamReq        = errors.New("search parameter is required")
	ErrEncoding        = errors.New("failed to marshal json")
	ErrNewParams       = errors.New("new command only accepts user/game")
	ErrWrongCode       = errors.New("wrong code")
	ErrWrongParam      = errors.New("wrong search parameter")
	ErrLoginTaken      = errors.New("login taken")
	ErrNoResponse      = errors.New("failed to get results")
	ErrCmdNotFound     = errors.New("command not found")
	ErrNotLoggedIn     = errors.New("user is not logged in")
	ErrUserNotFound    = errors.New("no user with this name exists")
	ErrWrongPassword   = errors.New("wrong password")
	ErrNotImplemented  = errors.New("not yet implemented")
	ErrNotEnoughParams = errors.New("not enough parameters")
)
