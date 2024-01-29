package main

import (
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	game "gitlab.com/Yoolayn/connect_four/internal/logic"
)

func addHandlers(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})
	r.GET("/users", getUsers)
	r.GET("/users/:login", getUser)
	r.GET("/games", getGames)
	r.GET("/games/:id", getGame)

	r.POST("/users", decoder(new(User)), newUser)
	r.POST("/games", decoder(new(Credentials)), authorizer(simpleCred), newGame)
	r.POST("/admins/:login", decoder(new(Credentials)), authorizer(simpleCred, true), changeAdmin(true))

	r.PUT("/games/:id/move")
	// r.PUT("/games/:id")
	// r.PUT("/users/:login/password")
	// r.PUT("/users/:login/name")

	r.DELETE("/admins/:login", decoder(new(Credentials)), authorizer(simpleCred, true), changeAdmin(false))

	r.POST("/secretsauce", decoder(new(repeatStruct)), authorizer(func(bdy interface{}) (Credentials, error) {
		body, ok := bdy.(*repeatStruct)
		if !ok {
			return Credentials{}, ErrType
		}
		return body.Credentials, nil
	}), repeat)
}

func simpleCred(bdy interface{}) (Credentials, error) {
	body, ok := bdy.(*Credentials)
	if !ok {
		return Credentials{}, ErrType
	}
	return *body, nil
}

func changeAdmin(to bool) func(c *gin.Context) {
	var name string
	if to {
		name = "setAdmin"
	} else {
		name = "unsetAdmin"
	}
	return func(c *gin.Context) {
		login, ok := c.Params.Get("login")
		if !ok {
			c.AbortWithStatusJSON(newErr(ErrInternal))
			log.Debug(name, "get login param", "login param not found")
			return
		}

		usr, ok := users.get(login)
		if !ok {
			c.AbortWithStatusJSON(newErr(ErrUserNotFound))
			log.Debug(name, "get user", ErrUserNotFound)
			return
		}
		users.Update(usr.MakeAdmin(to))
	}
}

func newGame(c *gin.Context) {
	game := game.MakeBoard()
	id := uuid.New()
	games[id] = Game{
		Board:   game,
		Title:   "",
		Player1: Player{},
		Player2: Player{},
	}
	c.String(http.StatusCreated, id.String())
}

func newUser(c *gin.Context) {
	bdy, ok := c.Get("decodedbody")
	if !ok {
		c.AbortWithStatusJSON(newErr(ErrInternal))
		log.Debug("newUser", "get body", ok)
		return
	}

	body, ok := bdy.(*User)
	if !ok {
		c.AbortWithStatusJSON(newErr(ErrType))
		log.Debug("newUser", "type assertion", ok)
		return
	}

	log.SetLevel(log.DebugLevel)
	log.Debug("newUser", "login", body.Login)
	log.Debug("newUser", "password", body.Password)
	_, ok = users.get(body.Login)
	if ok {
		c.AbortWithStatusJSON(newErr(ErrUserTaken))
		log.Debug("newUser", "user exists", ok)
		return
	}

	err := body.Encrypt()
	if err != nil {
		c.AbortWithStatusJSON(newErr(ErrPassTooLong))
		log.Debug("newUser", "encryption", err)
		return
	}
	body.FixEmpty()

	users.add(*body)
	c.Status(http.StatusCreated)
}

func getGame(c *gin.Context) {
	id, ok := c.Params.Get("id")
	if !ok {
		c.AbortWithStatusJSON(newErr(ErrInternal))
		log.Debug("getGame", "get id", ok)
		return
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		c.AbortWithStatusJSON(newErr(ErrParsing))
		log.Debug("getGame", "parsing", err)
		return
	}

	game, ok := games[uid]
	if !ok {
		c.AbortWithStatusJSON(newErr(ErrGameNotFound))
		log.Debug("getGame", "getting game", ok)
	}
	c.JSON(http.StatusOK, game)
}

func getGames(c *gin.Context) {
	c.JSON(http.StatusOK, games)
}

func getUsers(c *gin.Context) {
	c.JSON(http.StatusOK, users)
}

func getUser(c *gin.Context) {
	login, ok := c.Params.Get("login")
	if !ok {
		c.AbortWithStatusJSON(newErr(ErrInternal))
		log.Debug("getUser", "get id", ok)
		return
	}

	user, ok := users.get(login)
	if !ok {
		c.AbortWithStatusJSON(newErr(ErrUserNotFound))
		log.Debug("getUser", "get user", "user not found")
		return
	}

	c.JSON(http.StatusOK, user)
}

func repeat(c *gin.Context) {
	bdy, ok := c.Get("decodedbody")
	if !ok {
		c.AbortWithStatusJSON(newErr(ErrInternal))
		log.Debug("repeat", "get body", ok)
		return
	}

	body, ok := bdy.(*repeatStruct)
	if !ok {
		c.AbortWithStatusJSON(newErr(ErrType))
		log.Debug("repeat", "type assertion", ok)
		return
	}

	c.JSON(http.StatusOK, body)
}
