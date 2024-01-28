package main

import (
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
)

func addHandlers(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})
	r.GET("/users", getUsers)

	r.POST("/users", decoder(new(User)), newUser)
	r.POST("/secretsauce", decoder(new(repeatStruct)), authorize, repeat)
}

func newUser(c *gin.Context) {
	bdy, ok := c.Get("decodedbody")
	if !ok {
		c.AbortWithStatusJSON(newErr(ErrInternal))
		log.Error("newUser", "get body", ok)
		return
	}

	body, ok := bdy.(*User)
	if !ok {
		c.AbortWithStatusJSON(newErr(ErrType))
		log.Error("newUser", "type assertion", ok)
		return
	}

	log.SetLevel(log.DebugLevel)
	log.Debug("newUser", "login", body.Login)
	_, ok = users.get(body.Login)
	if ok {
		c.AbortWithStatusJSON(newErr(ErrUserTaken))
		log.Error("newUser", "user exists", ok)
		return
	}

	err := body.Encrypt()
	if err != nil {
		c.AbortWithStatusJSON(newErr(ErrPassTooLong))
		log.Error("NewUser", "encryption", err)
		return
	}

	users.add(*body)
	c.Status(http.StatusCreated)
}

func getUsers(c *gin.Context) {
	c.JSON(http.StatusOK, users)
}

func repeat(c *gin.Context) {
	bdy, ok := c.Get("decodedbody")
	if !ok {
		c.AbortWithStatusJSON(newErr(ErrInternal))
		log.Error("repeat", "get body", ok)
		return
	}

	body, ok := bdy.(*repeatStruct)
	if !ok {
		c.AbortWithStatusJSON(newErr(ErrType))
		log.Error("repeat", "type assertion", ok)
		return
	}

	c.JSON(http.StatusOK, body)
}
