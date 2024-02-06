package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/peterh/liner"
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

	var gms map[uuid.UUID]Game
	err = json.Unmarshal(body, &gms)
	if err != nil {
		return err
	}

	for k, v := range gms {
		fmt.Println("id:", k, ">>=", v)
	}

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
		var body struct {
			Users []User `json:"users"`
			Games []Game `json:"games"`
		}
		err = json.Unmarshal(data, &body)
		if err != nil {
			return err
		}

		if len(body.Users) != 0 {
			for i, v := range body.Users {
				fmt.Println(strconv.Itoa(i+1)+":", v.String())
			}
		} else if len(body.Games) != 0 {
			for i, v := range body.Games {
				fmt.Println(strconv.Itoa(i+1)+":", v.String())
			}
		} else {
			fmt.Println("no results")
		}

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
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return err
		}
		if string(body) == "" {
			body = []byte("NO_NAME")
		}
		creds.Login = payload.Login
		creds.Password = payload.Password
		creds.Name = string(body)
		return creds.Status()
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

	switch r.StatusCode {
	case 409:
		return ErrLoginTaken
	case 201:
		fmt.Println("user created")
		return nil
	default:
		return ErrRequest
	}
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

		var usrs []User
		err = json.Unmarshal(bytes, &usrs)
		if err != nil {
			return nil
		}

		for _, v := range usrs {
			fmt.Println(v.String())
		}

		return nil
	default:
		return ErrUnknown
	}
}

// user update name <new value>
func changeName(args string) error {
	if err := creds.Logged(); err != nil {
		return err
	}

	payload := struct {
		Credentials credentials `json:"credentials"`
		Name        string      `json:"newname"`
	}{
		Credentials: creds,
		Name:        args,
	}

	json, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, baseurl("/users/"+creds.Login+"/name"), bytes.NewBuffer(json))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	r, err := client.Do(req)
	if err != nil {
		return err
	}

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	switch r.StatusCode {
	case 202:
		fmt.Println("name changed to:", args)
		creds.Name = payload.Name
		return nil
	default:
		fmt.Println(string(bytes))
		return ErrUnknown
	}
}

func changePassword(args string) error {
	if err := creds.Logged(); err != nil {
		return err
	}

	payload := struct {
		Credentials credentials `json:"credentials"`
		Password    string      `json:"newpassword"`
	}{
		Credentials: creds,
		Password:    args,
	}

	json, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, baseurl("/users/"+creds.Login+"/password"), bytes.NewBuffer(json))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	r, err := client.Do(req)
	if err != nil {
		return err
	}

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	switch r.StatusCode {
	case 202:
		fmt.Println("password changed")
		creds.Password = payload.Password
		return nil
	default:
		fmt.Println(string(bytes))
		return ErrUnknown
	}
}

func deleteUser() error {
	line := liner.NewLiner()
	defer line.Close()
	line.SetCtrlCAborts(false)

	fmt.Print("\033[B")
	prompt, err := line.Prompt("are you sure? [yes/No]: ")
	fmt.Print("\033[A")
	if err != nil {
		return err
	}

	if prompt != "yes" && prompt != "Yes" {
		return nil
	}

	if err := creds.Logged(); err != nil {
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

	req, err := http.NewRequest(http.MethodDelete, baseurl("/users/"+creds.Login), bytes.NewBuffer(json))
	if err != nil {
		return err
	}

	client := http.Client{}

	r, err := client.Do(req)
	if err != nil {
		return err
	}

	switch r.StatusCode {
	case 200:
		creds = credentials{}
		fmt.Println("deleted user and logged out")
		return nil
	case 404:
		return ErrUserNotFound
	default:
		return ErrRequest
	}
}
