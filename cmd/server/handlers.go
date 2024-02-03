package main

import (
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	connectFour "gitlab.com/Yoolayn/connect_four/internal/logic"
)

func addHandlers(r *gin.Engine) {
	r.Use(static.Serve("/", static.LocalFile("./client/dist", true)))

	r.GET("/users", getUsers)
	r.GET("/users/:login", getUser)
	r.GET("/games", getGames)
	r.GET("/games/:id", getGame)

	r.POST("/users", decoder(new(User)), newUser)
	r.POST("/games", decoder(new(Credentials)), authorizer(simpleCred), newGame)
	r.POST("/admins/:login", decoder(new(Credentials)), authorizer(simpleCred, true), changeAdmin(true))
	r.POST("/games/:id", decoder(new(Join)), authorizer(func(bdy interface{}) (Credentials, error) {
		body, ok := bdy.(*Join)
		if !ok {
			return Credentials{}, ErrType
		}
		return body.Credentials, nil
		}), joinGame)

	r.PUT("/games/:id/move", decoder(new(Move)), authorizer(func(bdy interface{}) (Credentials, error) {
		body, ok := bdy.(*Move)
		if !ok {
			return Credentials{}, ErrType
		}

		return body.Credentials, nil
	}), makeMove)
	r.PUT("/games/:id", decoder(new(Title)), authorizer(func(bdy interface{}) (Credentials, error) {
		title, ok := bdy.(*Title)
		if !ok {
			return Credentials{}, ErrType
		}
		return title.Credentials, nil
	}), updateTitle)
	// r.PUT("/users/:login/password")
	// r.PUT("/users/:login/name")

	r.DELETE("/admins/:login", decoder(new(Credentials)), authorizer(simpleCred, true), changeAdmin(false))
	r.DELETE("/users/:login", decoder(new(Credentials)), authorizer(simpleCred), deleteUser)
	// r.DELETE("/games/:id")
	// r.DELETE("/games/:id/leave")

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})
	r.POST("/authtest", decoder(new(repeatStruct)), authorizer(func(bdy interface{}) (Credentials, error) {
		body, ok := bdy.(*repeatStruct)
		if !ok {
			return Credentials{}, ErrType
		}
		return body.Credentials, nil
	}), repeat)
}

func joinGame(c *gin.Context) {
	uid, ok := idToUUID(c, "joinGame")
	if !ok {
		return
	}

	game, ok := games[uid]
	if !ok {
		c.AbortWithStatusJSON(newErr(ErrGameNotFound))
		log.Debug("joinGame", "get game", ErrGameNotFound)
		return
	}

	body, ok := c.Get("decodedbody")
	if !ok {
		c.AbortWithStatusJSON(newErr(ErrInternal))
		log.Debug("joinGame", "get body", "failed to get the body")
		return
	}

	bdy, ok := body.(Credentials)
	if !ok {
		c.AbortWithStatusJSON(newErr(ErrType))
		log.Debug("joinGame", "type cast", ErrType)
		return
	}

	usr, ok := collections["users"].Get(bdy.Login)
	if !ok {
		c.AbortWithStatusJSON(newErr(ErrInternal))
		log.Debug("joinGame", "get user", "failed getting user from db")
		return
	}

	response := struct {
		Position int       `json:"position"`
		Game     uuid.UUID `json:"game"`
	}{}

	if game.Player1.User.Login == "" {
		game.Player1.User = usr
		response.Position = 1
	} else if game.Player2.User.Login == "" {
		game.Player2.User = usr
		response.Position = 2
	} else {
		c.AbortWithStatusJSON(newErr(ErrGameFull))
		return
	}

	response.Game = uid

	c.JSON(http.StatusOK, response)

}

func updateTitle(c *gin.Context) {
	uid, ok := idToUUID(c, "updateTitle")
	if !ok {
		return
	}

	body, ok := c.Get("decodedbody")
	if !ok {
		c.AbortWithStatusJSON(newErr(ErrInternal))
		logger.Debug("updateTitle", "get body", ok)
		return
	}

	bdy, ok := body.(*Title)
	if !ok {
		c.AbortWithStatusJSON(newErr(ErrType))
		logger.Debug("updateTitle", "convert type", ok)
		return
	}

	game, ok := games[uid]
	if !ok {
		c.AbortWithStatusJSON(newErr(ErrGameNotFound))
		logger.Debug("updateTitle", "find game", ErrGameNotFound)
		return
	}

	if bdy.Credentials.Login != game.Player1.User.Login || bdy.Credentials.Login != game.Player2.User.Login {
		c.AbortWithStatusJSON(newErr(ErrNotInGame))
		logger.Debug("updateTitle", "change title", ErrNotInGame)
		return
	}

	game.Title = bdy.Title
	games[uid] = game
	c.Status(http.StatusAccepted)
}

func deleteUser(c *gin.Context) {
	login, ok := c.Params.Get("login")
	if !ok {
		c.AbortWithStatusJSON(newErr(ErrInternal))
		logger.Debug("deleteUser", "get login param", "login param not found")
		return
	}

	ok = collections["users"].Delete(login)
	if !ok {
		c.AbortWithStatusJSON(newErr(ErrUserNotFound))
		logger.Debug("deleteUser", "find user to delete", ok)
		return
	}

	c.Status(http.StatusOK)
}

func makeMove(c *gin.Context) {
	uid, ok := idToUUID(c, "makeMove")
	if !ok {
		return
	}

	game, ok := games[uid]
	if !ok {
		c.AbortWithStatusJSON(newErr(ErrGameNotFound))
		logger.Debug("makeMove", "find game", ok)
		return
	}

	body, ok := c.Get("decodedbody")
	if !ok {
		c.AbortWithStatusJSON(newErr(ErrInternal))
		logger.Debug("makeMove", "get body", ok)
		return
	}

	bdy, ok := body.(*Move)
	if !ok {
		c.AbortWithStatusJSON(newErr(ErrType))
		logger.Debug("makeMove", "convert type", ok)
		return
	}

	if bdy.Row > 6 || bdy.Row < 0 {
		c.AbortWithStatusJSON(newErr(ErrOutOfBounds))
		return
	}

	var chkr connectFour.Checker
	if bdy.Credentials.Login == game.Player1.User.Login {
		chkr = connectFour.Checker{Color: game.Player1.Color}
	} else if bdy.Credentials.Login == game.Player2.User.Login {
		chkr = connectFour.Checker{Color: game.Player2.Color}
	}

	if chkr.Color == "" {
		c.AbortWithStatusJSON(newErr(ErrNotInGame))
		return
	}

	ok = game.Board.Claim(chkr, bdy.Row)
	if !ok {
		c.AbortWithStatusJSON(newErr(ErrFieldTaken))
		return
	}

	games[uid] = game
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
			logger.Debug(name, "get login param", "login param not found")
			return
		}

		usr, ok := collections["users"].Get(login)
		if !ok {
			c.AbortWithStatusJSON(newErr(ErrUserNotFound))
			logger.Debug(name, "get user", ErrUserNotFound)
			return
		}

		admin := usr.MakeAdmin(to)
		ok = collections["users"].Update(admin.Login, admin)
		if !ok {
			c.AbortWithStatusJSON(newErr(ErrUpdateFailed))
			logger.Debug(name, "update user", ErrUpdateFailed)
			return
		}
	}
}

func newGame(c *gin.Context) {
	game := connectFour.MakeBoard()
	id := uuid.New()
	games[id] = Game{
		Board:   game,
		Title:   "New Game",
		Player1: Player{},
		Player2: Player{},
	}
	c.String(http.StatusCreated, id.String())
}

func newUser(c *gin.Context) {
	bdy, ok := c.Get("decodedbody")
	if !ok {
		c.AbortWithStatusJSON(newErr(ErrInternal))
		logger.Debug("newUser", "get body", ok)
		return
	}

	body, ok := bdy.(*User)
	if !ok {
		c.AbortWithStatusJSON(newErr(ErrType))
		logger.Debug("newUser", "type assertion", ok)
		return
	}

	logger.Debug("newUser", "login", body.Login)
	logger.Debug("newUser", "password", body.Password)

	_, ok = collections["users"].Get(body.Login)
	if ok {
		c.AbortWithStatusJSON(newErr(ErrUserTaken))
		logger.Debug("newUser", "user exists", ok)
		return
	}

	err := body.Encrypt()
	if err != nil {
		c.AbortWithStatusJSON(newErr(ErrPassTooLong))
		logger.Debug("newUser", "encryption", err)
		return
	}
	body.FixEmpty()

	ok = collections["users"].Add(*body)
	if !ok {
		c.AbortWithStatusJSON(newErr(ErrAddFailed))
		return
	}

	c.Status(http.StatusCreated)
}

func getGame(c *gin.Context) {
	uid, ok := idToUUID(c, "getGame")
	if !ok {
		return
	}

	game, ok := games[uid]
	if !ok {
		c.AbortWithStatusJSON(newErr(ErrGameNotFound))
		logger.Debug("getGame", "getting game", ok)
	}
	c.JSON(http.StatusOK, game)
}

func getGames(c *gin.Context) {
	c.JSON(http.StatusOK, games)
}

func getUsers(c *gin.Context) {
	var users []User
	ok := collections["users"].GetAll(&users)
	if !ok {
		c.AbortWithStatusJSON(newErr(ErrGetAllFailed))
		return
	}

	c.JSON(http.StatusOK, users)
}

func getUser(c *gin.Context) {
	login, ok := c.Params.Get("login")
	if !ok {
		c.AbortWithStatusJSON(newErr(ErrInternal))
		logger.Debug("getUser", "get id", ok)
		return
	}

	user, ok := collections["users"].Get(login)
	if !ok {
		c.AbortWithStatusJSON(newErr(ErrUserNotFound))
		logger.Debug("getUser", "get user", "user not found")
		return
	}

	c.JSON(http.StatusOK, user)
}

func repeat(c *gin.Context) {
	bdy, ok := c.Get("decodedbody")
	if !ok {
		c.AbortWithStatusJSON(newErr(ErrInternal))
		logger.Debug("repeat", "get body", ok)
		return
	}

	body, ok := bdy.(*repeatStruct)
	if !ok {
		c.AbortWithStatusJSON(newErr(ErrType))
		logger.Debug("repeat", "type assertion", ok)
		return
	}

	c.JSON(http.StatusOK, body)
}
