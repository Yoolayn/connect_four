package main

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
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

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Name     string `json:"name,omitempty"`
	isAdmin  bool
}

func (u User) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Login string `json:"login"`
		Name  string `json:"name,omitempty"`
	}{
		Login: u.Login,
		Name: u.Name,
	})
}

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Users []User

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

func (u Users) get(login string) (User, bool) {
	for _, v := range u {
		if v.Login == login {
			return v, true
		}
	}
	return User{}, false
}

func (u *Users) add(usr User) {
	*u = append(*u, usr)
}
