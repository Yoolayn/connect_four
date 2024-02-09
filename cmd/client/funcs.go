package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	"github.com/peterh/liner"
	"gitlab.com/Yoolayn/connect_four/internal/browser"
)

func baseurl(path string) string {
	return "http://localhost:8080" + path
}

func hello(m string) error {
	fmt.Println("welcome", m)
	return nil
}

func join() error {
	if err := creds.Logged(); err != nil {
		return err
	}

	res, err := http.Get(baseurl("/games"))
	if err != nil {
		return err
	}

	switch res.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return ErrUnknown
	default:
		return ErrUnknown
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

	sli := make([]string, 0)
	gmsStr := make(map[string]Game)
	for i, v := range gms {
		gmsStr[i.String()] = v
		sli = append(sli, i.String())
	}

	chosen, err := browser.New("Games:", sli)
	if err != nil {
		return err
	}

	colors := map[string]string{
		"red":    "#ff0000",
		"yellow": "#ffff00",
	}
	colorKeys := make([]string, len(colors))
	i := 0
	for k := range colors {
		colorKeys[i] = k
		i++
	}

	color, err := browser.New("Color:", colorKeys)
	if err != nil {
		return err
	}

	// chosenAsUUID, err := uuid.Parse(chosen)
	// if err != nil {
	// 	return err
	// }

	payload := struct {
		C   credentials `json:"credentials"`
		Col string      `json:"color"`
	}{
		C:   creds,
		Col: colors[color],
	}

	mars, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	r, err := http.Post(baseurl("/games/"+chosen), "application/json", bytes.NewBuffer(mars))
	if err != nil {
		return err
	}

	var bodyResp struct {
		Position int       `json:"position"`
		Game     uuid.UUID `json:"game"`
	}
	switch r.StatusCode {
	case http.StatusBadRequest:
		return ErrGameFull
	case http.StatusNotFound:
		return ErrGameNotFound
	case http.StatusOK:
		earth, err := io.ReadAll(r.Body)
		if err != nil {
			return err
		}
		err = json.Unmarshal(earth, &bodyResp)
		if err != nil {
			return err
		}
	default:
		return ErrUnknown
	}

	resp, err := http.Get(baseurl("/games/" + bodyResp.Game.String()))
	if err != nil {
		return err
	}

	switch res.StatusCode {
	case http.StatusNotFound:
		return ErrGameNotFound
	case http.StatusOK:
	default:
		return ErrUnknown
	}

	var game Game
	bits, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bits, &game)
	if err != nil {
		return err
	}

	start, err := InitialModel(bodyResp.Position, game, bodyResp.Game)
	if err != nil {
		return err
	}

	p := tea.NewProgram(start, tea.WithAltScreen())
	_, err = p.Run()
	if err != nil {
		return err
	}

	return nil
}

func game() error {
	res, err := http.Get(baseurl("/games"))
	if err != nil {
		return err
	}

	switch res.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return ErrUnknown
	default:
		return ErrUnknown
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

	var gmsStr []string
	for k := range gms {
		gmsStr = append(gmsStr, k.String())
	}

	chosen, err := browser.New("Games:", gmsStr)
	if err != nil {
		return err
	}

	r, err := http.Get(baseurl("/games/"+chosen))
	if err != nil {
		return err
	}

	bitties, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	var game Game
	err = json.Unmarshal(bitties, &game)
	if err != nil {
		return err
	}

	fmt.Println(game.String())

	return nil
}

func games() error {
	res, err := http.Get(baseurl("/games"))
	if err != nil {
		return err
	}

	switch res.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return ErrUnknown
	default:
		return ErrUnknown
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
		switch response.StatusCode {
		case http.StatusOK:
		case http.StatusAccepted:
		case http.StatusNotFound:
			return ErrUnknown
		default:
			return ErrUnknown
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
		if k == "/op" || k == "/deop" {
			continue
		}
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
	case http.StatusNotFound:
		return ErrUserNotFound
	case http.StatusUnauthorized:
		return ErrWrongPassword
	case http.StatusOK:
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
	case http.StatusConflict:
		return ErrLoginTaken
	case http.StatusCreated:
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
	case http.StatusCreated:
		id, err := io.ReadAll(r.Body)
		if err != nil {
			return ErrUnknown
		}
		fmt.Println("new game created with id", string(id), "and title \"New Game\"")
		return nil
	case http.StatusNotFound:
		return ErrUserNotFound
	default:
		return ErrUnknown
	}
}

func user() error {
	r, err := http.Get(baseurl("/users"))
	if err != nil {
		return err
	}

	switch r.StatusCode {
	case http.StatusOK:
	default:
		return ErrUnknown
	}

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	var usrs []User
	err = json.Unmarshal(bytes, &usrs)
	if err != nil {
		return nil
	}

	var usrKeys []string
	for _, v := range usrs {
		usrKeys = append(usrKeys, v.Login)
	}

	chosen, err := browser.New("Users:", usrKeys)
	if err != nil {
		return err
	}

	res, err := http.Get(baseurl("/users/" + chosen))
	if err != nil {
		return err
	}

	switch res.StatusCode {
	case http.StatusOK:
	default:
		return ErrRequest
	}

	bitties, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var usr User
	err = json.Unmarshal(bitties, &usr)
	if err != nil {
		return err
	}

	fmt.Println(usr.String())

	return nil
}

func users() error {
	r, err := http.Get(baseurl("/users"))
	if err != nil {
		return err
	}
	switch r.StatusCode {
	case http.StatusOK:
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
	case http.StatusAccepted:
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
	case http.StatusAccepted:
		fmt.Println("password changed")
		creds.Password = payload.Password
		return nil
	default:
		fmt.Println(string(bytes))
		return ErrUnknown
	}
}

func deleteUser() error {
	if err := creds.Logged(); err != nil {
		return err
	}

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
	case http.StatusOK:
		creds = credentials{}
		fmt.Println("deleted user and logged out")
		return nil
	case http.StatusNotFound:
		return ErrUserNotFound
	default:
		return ErrRequest
	}
}

func makeAdmin(args string) error {
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

	resp, err := http.Post(baseurl("/admins/"+args), "application/json", bytes.NewBuffer(json))
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		fallthrough
	case http.StatusAccepted:
		fmt.Println("Privileges elevated")
		return nil
	case http.StatusForbidden:
		return ErrAdminRequired
	default:
		return ErrUnknown
	}
}

func removeAdmin(args string) error {
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

	req, err := http.NewRequest(http.MethodDelete, baseurl("/admins/"+args), bytes.NewBuffer(json))
	if err != nil {
		return err
	}

	c := http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return err
	}

	switch res.StatusCode {
	case http.StatusOK:
		fallthrough
	case http.StatusAccepted:
		fmt.Println("Privileges de elevated")
		return nil
	case http.StatusForbidden:
		return ErrAdminRequired
	default:
		return ErrUnknown
	}
}

func deleteGame() error {
	if err := creds.Logged(); err != nil {
		return err
	}

	r, err := http.Get(baseurl("/games"))
	if err != nil {
		return err
	}

	switch r.StatusCode {
	case http.StatusOK:
	case http.StatusAccepted:
	case http.StatusNotFound:
		return ErrUnknown
	default:
		return ErrUnknown
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	var gms map[uuid.UUID]Game
	err = json.Unmarshal(body, &gms)
	if err != nil {
		return err
	}

	var gmsKeys []string
	gmsStr := make(map[string]Game)
	for k, v := range gms {
		gmsStr[k.String()] = v
		gmsKeys = append(gmsKeys, k.String())
	}

	chosen, err := browser.New("Choose to delete:", gmsKeys)
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

	bitties, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodDelete, baseurl("/games/"+chosen), bytes.NewBuffer(bitties))
	if err != nil {
		return err
	}

	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return err
	}

	switch res.StatusCode {
	case http.StatusOK:
		fmt.Println("game deleted")
	case http.StatusForbidden:
		return ErrAdminRequired
	default:
		return ErrRequest
	}

	return nil
}
