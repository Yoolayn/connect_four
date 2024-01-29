package main

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	game "gitlab.com/Yoolayn/connect_four/internal/logic"
	"golang.org/x/crypto/bcrypt"
)

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

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Name     string `json:"name"`
	IsAdmin  bool   `json:"isadmin"`
}

type Player struct {
	User  User   `json:"user"`
	Color string `json:"color"`
}

type Game struct {
	Board   game.Board `json:"board"`
	Title   string     `json:"title"`
	Player1 Player     `json:"player1"`
	Player2 Player     `json:"player2"`
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

func (u Users) get(login string) (*User, bool) {
	for _, v := range u {
		if v.Login == login {
			return &v, true
		}
	}
	return &User{}, false
}

func (u *Users) Update(user User) {
	for i, v := range *u {
		if v.Login == user.Login {
			(*u)[i] = user
		}
	}
}

func (u *Users) add(usr User) {
	*u = append(*u, usr)
}
