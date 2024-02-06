package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func baseurl(path string) string {
	return "http://localhost:8080" + path
}

func hello(m string) error {
	fmt.Println("welcome", m)
	return nil
}

func games() error {
	res, err := http.Get(baseurl("/games"))
	if err != nil {
		return err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	return nil
}

func search(args string) error {
	split := strings.SplitN(args, " ", 2)
	if len(split) < 2 {
		return ErrParamReq
	}

	pattern := strings.ReplaceAll(split[1], " ", "%20")

	switch split[0] {
	case "user":
		fallthrough
	case "users":
		fallthrough
	case "game":
		fallthrough
	case "games":
		response, err := http.Get(baseurl("/search?") + split[0] + "=" + pattern)
		if err != nil {
			return err
		}
		data, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	default:
		return ErrWrongParam
	}
	return nil
}

func help() error {
	fmt.Println("available commands:")
	var iter int
	for k := range cmds {
		iter++
		if iter == len(cmds) {
			fmt.Print(k)
			break
		}
		fmt.Printf("%s, ", k)
	}
	fmt.Print("\n")
	return nil
}

func login(args string) error {
	argSplit := strings.Split(args, " ")
	if len(argSplit) != 2 {
		return ErrWrongParam
	}

	payload := struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}{
		Login:    argSplit[0],
		Password: argSplit[1],
	}

	json, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	r, err := http.Post(baseurl("/login"), "application/json", bytes.NewBuffer(json))
	if err != nil {
		return err
	}

	switch r.StatusCode {
	case 404:
		return ErrUserNotFound
	case 401:
		return ErrWrongPassword
	case 200:
		creds.Login = payload.Login
		creds.Password = payload.Password
		fmt.Println("logged in")
		fmt.Printf("%#v\n", creds)
		return nil
	default:
		return ErrUnknown
	}
}

// new user <login> <password> <name>
func newUser(args string) error {
	argSplit := strings.SplitN(args, " ", 3)
	if len(argSplit) < 2 {
		return ErrNotEnoughParams
	}

	payload := struct {
		Login    string `json:"login"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}{
		Login:    argSplit[0],
		Password: argSplit[1],
	}

	if len(argSplit) != 2 {
		payload.Name = argSplit[2]
	}

	json, err := json.Marshal(payload)
	if err != nil {
		return ErrEncoding
	}

	r, err := http.Post(baseurl("/users"), "application/json", bytes.NewBuffer(json))
	if err != nil {
		return err
	}

	if r.StatusCode == 409 {
		return ErrLoginTaken
	}

	if r.StatusCode != 201 {
		return ErrRequest
	}

	fmt.Println(r.Status)

	return nil
}

func newGame() error {
	err := creds.Logged()
	if err != nil {
		return err
	}

	payload := struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}{
		Login:    creds.Login,
		Password: creds.Password,
	}

	json, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	r, err := http.Post(baseurl("/games"), "application/json", bytes.NewBuffer(json))
	if err != nil {
		return err
	}

	switch r.StatusCode {
	case 201:
		id, err := io.ReadAll(r.Body)
		if err != nil {
			return ErrUnknown
		}
		fmt.Println("new game created with id", string(id), "and title \"New Game\"")
		return nil
	case 404:
		return ErrUserNotFound
	default:
		return ErrUnknown
	}
}

func users() error {
	r, err := http.Get(baseurl("/users"))
	if err != nil {
		return ErrRequest
	}
	switch r.StatusCode {
	case 200:
		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			return err
		}
		fmt.Println(string(bytes))
		return nil
	default:
		return ErrUnknown
}
}
