package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/peterh/liner"
	"golang.org/x/term"
)

type credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (c credentials) Logged() error {
	if c.Login == "" || c.Password == "" {
		return ErrNotLoggedIn
	}
	return nil
}

var (
	cmds   = make(map[string]func(string) error, 0)
	height = 0
	buffer bytes.Buffer
	creds  credentials
)

func processing(line string) error {
	line = strings.ToLower(line)
	words := strings.SplitN(line, " ", 2)

	cmd := words[0]
	var args string
	if len(words) == 1 {
		args = ""
	} else {
		args = words[1]
	}

	return dispatch(cmd, args)
}

func dispatch(cmd, args string) error {
	fn, ok := cmds[cmd]
	if !ok {
		return ErrCmdNotFound
	}
	return fn(args)
}

func moveBottom() {
	fmt.Printf("\033[%d;0H", height)
}

func fullClear() {
	fmt.Print("\033[2J\033[H")
}

// func newMessage(content string) {
// 	fmt.Print("\033[s")
// 	fmt.Print("\033[A")
// 	fmt.Println("\n" + content)
// 	fmt.Print(">>= ")
// 	fmt.Print("\033[u")
// }

func prompter() {
	moveBottom()

	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)

	_, err := line.ReadHistory(&buffer)
	if err != nil {
		panic(err)
	}

	if prompt, err := line.Prompt(">>= "); err == nil {
		err := processing(prompt)
		if err != nil {
			fmt.Println(err)
		}
		line.AppendHistory(prompt)
	} else if err == liner.ErrPromptAborted || err.Error() == "EOF" {
		fmt.Println("Exiting...")
		return
	} else {
		fmt.Println("error reading", err)
	}
	defer prompter()

	_, err = line.WriteHistory(&buffer)
	if err != nil {
		panic(err)
	}
}

func main() {
	defer prompter()
	_, h, err := term.GetSize(0)
	if err != nil {
		panic(err)
	}
	height = h

	fullClear()
	moveBottom()

	cmds["search"] = func(args string) error {
		if args == "" {
			return ErrArgsReq
		}
		return search(args)
	}
	cmds["games"] = func(args string) error {
		return games()
	}
	cmds["hello"] = func(args string) error {
		if args == "" {
			return ErrArgsReq
		}
		return hello(args)
	}
	cmds["help"] = func(args string) error {
		return help()
	}
	cmds["new"] = func(args string) error {
		argSplit := strings.SplitN(args, " ", 2)
		switch argSplit[0] {
		case "user":
			if len(argSplit) != 2 {
				return ErrNotEnoughParams
			}
			return newUser(argSplit[1])
		case "game":
			return newGame()
		default:
			return ErrNewParams
		}
	}
	cmds["login"] = func(args string) error {
		if args == "" {
			return ErrArgsReq
		}
		return login(args)
	}
	cmds["users"] = func(args string) error {
		return users()
	}
	cmds["user"] = func(args string) error {
		// user <update|delete> <what> <args>
		argSplit := strings.SplitN(args, " ", 3)
		if len(argSplit) < 1 {
			return ErrNotEnoughParams
		}
		switch argSplit[0] {
		case "update":
			if len(argSplit) != 3 {
				return ErrNotEnoughParams
			}
			switch argSplit[1] {
			case "name":
				return changeName(argSplit[2])
			case "login":
				return ErrNotImplemented
			default:
				return ErrNotImplemented
			}
		case "delete":
			return ErrNotImplemented
		default:
			return ErrNewParams
		}
	}
}
