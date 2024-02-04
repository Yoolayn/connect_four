package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	game "gitlab.com/Yoolayn/connect_four/internal/logic"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type collection struct {
	c    *mongo.Collection
	name string
}

func (c collection) Update(login string, u User) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	_, err := c.c.UpdateOne(ctx, bson.M{"login": login}, bson.M{"$set": u.ToDB()})
	if err != nil {
		logger.Debug("update in database", c.name, err)
		return false
	}

	return true
}

func (c collection) Get(login string) (User, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	var usr User
	err := collections["users"].c.FindOne(ctx, bson.M{"login": login}).Decode(&usr)
	if err != nil {
		logger.Debug("Get from database", c.name, err)
		return usr, false
	}

	return usr, true
}

func (c collection) Add(u User) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	_, err := c.c.InsertOne(ctx, u.ToDB())
	if err != nil {
		logger.Debug("Add to database", c.name, err)
		return false
	}

	return true
}

func (c collection) GetAll(result interface{}) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	cursor, err := c.c.Find(ctx, bson.M{})
	if err != nil {
		logger.Debug("GetAll", c.name, err)
		return false
	}

	err = cursor.All(ctx, result)
	if err != nil {
		logger.Debug("GetAll", c.name, err)
		return false
	}

	return true
}

func (c collection) Delete(login string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	r, err := c.c.DeleteOne(ctx, bson.M{"login": login})
	if err != nil {
		log.Debug("Delete from database", c.name, err)
		return false
	}

	if r.DeletedCount != 1 {
		log.Debug("Delete from database", c.name, "not found")
		return false
	}
	return true
}

type repeatStruct struct {
	Rest struct {
		Hello string `json:"hello"`
	} `json:"rest"`
	Credentials struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	} `json:"credentials"`
}

type Move struct {
	Credentials Credentials `json:"credentials"`
	Row         int         `json:"row"`
}

type Title struct {
	Credentials Credentials `json:"credentials"`
	Title       string      `json:"title"`
}

type Join struct {
	Credentials Credentials `json:"credentials"`
	Color       string      `json:"color"`
}

type NewName struct {
	Credentials Credentials `json:"credentials"`
	NewName     string      `json:"newname"`
}

type User struct {
	Login    string `json:"login" bson:"login"`
	Password string `json:"password" bson:"password"`
	Name     string `json:"name" bson:"name"`
	IsAdmin  bool   `json:"isadmin" bson:"isadmin"`
}

type Player struct {
	User  User   `json:"user"`
	Color string `json:"color"`
}

type NewPassword struct {
	Credentials Credentials `json:"credentials"`
	NewPassword string      `json:"newpassword"`
}

type Game struct {
	Board   game.Board `json:"board"`
	Title   string     `json:"title"`
	Player1 Player     `json:"player1"`
	Player2 Player     `json:"player2"`
}

func (u User) ToDB() struct {
	Login    string `json:"login" bson:"login"`
	Name     string `json:"name" bson:"name"`
	Password string `json:"password" bson:"password"`
	IsAdmin  bool   `json:"isadmin" bson:"isadmin"`
} {
	return struct {
		Login    string `json:"login" bson:"login"`
		Name     string `json:"name" bson:"name"`
		Password string `json:"password" bson:"password"`
		IsAdmin  bool   `json:"isadmin" bson:"isadmin"`
	}{
		Login:    u.Login,
		Name:     u.Name,
		Password: u.Password,
		IsAdmin:  u.IsAdmin,
	}
}

func (u User) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Login   string `json:"login"`
		Name    string `json:"name,omitempty"`
		IsAdmin bool   `json:"admin"`
	}{
		Login:   u.Login,
		Name:    u.Name,
		IsAdmin: u.IsAdmin,
	})
}

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Users []User

func (u User) MakeAdmin(status bool) User {
	u.IsAdmin = status
	return u
}

func (u *User) FixEmpty() {
	if u.Name == "" {
		id := uuid.NewString()[:8]
		u.Name = fmt.Sprintf("anonymous-%s", id)
	}
}

func (u *User) Encrypt() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}
