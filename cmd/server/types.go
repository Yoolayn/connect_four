package main

import "golang.org/x/crypto/bcrypt"

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
}

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Users []User

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
