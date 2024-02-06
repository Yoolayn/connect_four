package main

import (
	"strings"
)

type Checker struct {
	Color string
}

type row []Checker

type Board []row

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Name     string `json:"name"`
	IsAdmin  bool   `json:"admin"`
}

type Player struct {
	User  User   `json:"user"`
	Color string `json:"color"`
}

func (u User) String() string {
	name := u.Name
	login := u.Login
	var str string
	if name == "" && login == "" {
		str = "empty"
	} else {
		str = name + "(" + login + ")"
	}
	if u.IsAdmin {
		str = "Administrator: " + str
	} else {
		str = "User: " + str
	}
	return str
}

type Game struct {
	Board   Board     `json:"board"`
	Title   string    `json:"title"`
	Player1 Player    `json:"player1"`
	Player2 Player    `json:"player2"`
}

func (g Game) String() string {
	title := "Title: " + g.Title

	prefix := "User: "
	if g.Player1.User.IsAdmin {
		prefix = "Administrator: "
	}
	player1 := strings.TrimPrefix(g.Player1.User.String(), prefix)

	player1 = "Player1: " + player1
	if player1 != "Player1: empty" {
		player1 = player1 + " - " + g.Player1.Color
	}

	prefix = "User: "
	if g.Player2.User.IsAdmin {
		prefix = "Administrator: "
	}
	player2 := strings.TrimPrefix(g.Player2.User.String(), prefix)

	player2 = "Player2: " + player2
	if player2 != "Player2: empty" {
		player2 = player2 + " - " + g.Player2.Color
	}

	return title + ", " + player1 + ", " + player2
}
