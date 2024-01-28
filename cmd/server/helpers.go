package main

import (
	"encoding/json"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func decodeBody(c *gin.Context, body interface{}) {
	err := json.NewDecoder(c.Request.Body).Decode(body)
	if err != nil {
		c.AbortWithStatusJSON(newErr(ErrInternal))
		log.Error("decodeBody", "json decode", err)
		return
	}
	c.Set("decodedbody", body)
}

func decoder(bdy interface{}) func(*gin.Context) {
	return func(c *gin.Context) {
		decodeBody(c, bdy)
	}
}

func authorizer(fn func(bdy any) (Credentials, error)) func(*gin.Context) {
	return func(c *gin.Context) {
		authorize(c, fn)
	}
}

func authorize(c *gin.Context, fn func(bdy any) (Credentials, error)) {
	paramBdy, ok := c.Get("decodedbody")
	if !ok {
		c.AbortWithStatusJSON(newErr(ErrInternal))
		log.Error("authorize", "get body", ok)
		return
	}

	body, err := fn(paramBdy)
	if err != nil {
		c.AbortWithStatusJSON(newErr(err))
		log.Error("authorize", "create credentials", err)
		return
	}

	usr, ok := users.get(body.Login)
	if !ok {
		c.AbortWithStatusJSON(newErr(ErrUserNotFound))
		log.Error("authorize", "get login", ok)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(body.Password))
	if err != nil {
		c.AbortWithStatusJSON(newErr(ErrNotAuthorized))
		log.Error("authorize", "password match", err)
	}
}
