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

func authorize(c *gin.Context) {
	bdy, ok := c.Get("decodedbody")
	if !ok {
		c.AbortWithStatusJSON(newErr(ErrInternal))
		log.Error("authorize", "get body", ok)
		return
	}
	body, ok := bdy.(*repeatStruct)
	if !ok {
		c.AbortWithStatusJSON(newErr(ErrType))
		log.Error("authorize", "type assertion", ok)
		return
	}

	usr, ok := users.get(body.Credentials.Login)
	if !ok {
		c.AbortWithStatusJSON(newErr(ErrUserNotFound))
		log.Error("authorize", "get login", ok)
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(body.Credentials.Password))
	if err != nil {
		c.AbortWithStatusJSON(newErr(ErrNotAuthorized))
		log.Error("authorize", "password match", err)
	}
}
