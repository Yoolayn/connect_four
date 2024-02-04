package main

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func decodeBody(c *gin.Context, body interface{}) {
	err := json.NewDecoder(c.Request.Body).Decode(body)
	if err != nil {
		c.AbortWithStatusJSON(newErr(ErrInternal))
		logger.Error("decodeBody", "json decode", err)
		return
	}
	c.Set("decodedbody", body)
}

func decoder(bdy interface{}) func(*gin.Context) {
	return func(c *gin.Context) {
		logger.Debug("decoding body")
		decodeBody(c, bdy)
	}
}

func idToUUID(c *gin.Context, name string) (uuid.UUID, bool) {
	id, ok := c.Params.Get("id")
	if !ok {
		c.AbortWithStatusJSON(newErr(ErrInternal))
		logger.Debug(name, "get login param", "login param not found")
		return uuid.Nil, false
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		c.AbortWithStatusJSON(newErr(ErrParsing))
		logger.Debug(name, "parse uuid", err)
		return uuid.Nil, false
	}

	return uid, true
}

func authorizer(fn func(bdy any) (Credentials, error), admin ...bool) func(*gin.Context) {
	auth := func(c *gin.Context) {
		logger.Debug("auth started")
		authorize(c, fn)
	}
	if len(admin) > 0 && admin[0] {
		auth = func(c *gin.Context) {
			logger.Debug("auth started")
			authorize(c, fn, func(u User) bool {
				return u.IsAdmin
			})
		}
	}
	return auth
}

func authorize(c *gin.Context, fn func(bdy any) (Credentials, error), admin ...func(u User) bool) {
	paramBdy, ok := c.Get("decodedbody")
	if !ok {
		c.AbortWithStatusJSON(newErr(ErrInternal))
		logger.Error("authorize", "get body", ok)
		return
	}

	body, err := fn(paramBdy)
	if err != nil {
		c.AbortWithStatusJSON(newErr(err))
		logger.Error("authorize", "create credentials", err)
		return
	}

	usr, ok := collections["users"].Get(body.Login)
	if !ok {
		c.AbortWithStatusJSON(newErr(ErrUserNotFound))
		logger.Error("authorize", "get login", ok)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(body.Password))
	if err != nil {
		c.AbortWithStatusJSON(newErr(ErrNotAuthorized))
		logger.Error("authorize", "password match", err)
		return
	}

	if len(admin) > 0 {
		ok := admin[0](usr)
		if !ok {
			c.AbortWithStatusJSON(newErr(ErrForbidden))
			logger.Error("authorize", "admin requirement", "failed")
		}
	}
}
