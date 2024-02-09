package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/peterh/liner"
	"golang.org/x/term"
)

type credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func (c credentials) Logged() error {
	if c.Login == "" || c.Password == "" {
		return ErrNotLoggedIn
	}
	return nil
}

func (c credentials) Status() error {
	if err := c.Logged(); err == nil {
		fmt.Println("logged in as", "\""+creds.Name+"\""+" "+"("+creds.Login+")")
		return nil
	} else {
		return err
	}
}

var (
	cmds   = make(map[string]func(string) error, 0)
	height = 0
	buffer bytes.Buffer
	creds  = credentials{}
)

func processing(line string) error {
	line = strings.TrimSpace(line)
	words := strings.SplitN(line, " ", 2)

	cmd := words[0]
	var args string
	if len(words) == 1 {
		args = ""
	} else {
		args = words[1]
	}

	err := dispatch(cmd, args)
	return err
}

func dispatch(cmd, args string) error {
	fn, ok := cmds[cmd]
	if !ok {
		if cmd == "" {
			return nil
		}
		if ok := cmd[0:1] == "/"; ok {
			return ErrCmdNotFound
		} else {
			if err := creds.Logged(); err == nil {
				fmt.Println(creds.Login + ": " + cmd + " " + args)
				return nil
			} else {
				fmt.Println("not logged in")
				return nil
			}
		}
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

	cmds["/search"] = func(args string) error {
		if args == "" {
			return ErrArgsReq
		}
		return search(args)
	}
	cmds["/games"] = func(args string) error {
		return games()
	}
	cmds["/hello"] = func(args string) error {
		if args == "" {
			return ErrArgsReq
		}
		return hello(args)
	}
	cmds["/help"] = func(args string) error {
		return help()
	}
	cmds["/new"] = func(args string) error {
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
	cmds["/login"] = func(args string) error {
		if args == "" {
			return ErrArgsReq
		}
		return login(args)
	}
	cmds["/users"] = func(args string) error {
		return users()
	}
	cmds["/user"] = func(args string) error {
		argSplit := strings.SplitN(args, " ", 3)
		if len(argSplit) < 1 {
			return ErrNotEnoughParams
		}
		switch argSplit[0] {
		case "change":
			fallthrough
		case "update":
			if len(argSplit) != 3 {
				return ErrNotEnoughParams
			}
			switch argSplit[1] {
			case "name":
				return changeName(argSplit[2])
			case "password":
				return changePassword(argSplit[2])
			default:
				return ErrUserUpdate
			}
		case "delete":
			return deleteUser()
		default:
			return ErrUserParams
		}
	}
	cmds["/status"] = func(args string) error {
		return creds.Status()
	}
	cmds["/logout"] = func(args string) error {
		if err := creds.Logged(); err == nil {
			fmt.Println("Logged out")
			creds = credentials{}
			return nil
		} else {
			return ErrNotLoggedIn
		}
	}
	cmds["/join"] = func(args string) error {
		return join()
	}
	cmds["/exit"] = func(args string) error {
		fmt.Println("Exiting...")
		os.Exit(0)
		return nil
	}
	cmds["/op"] = func(args string) error {
		return makeAdmin(args)
	}
	cmds["/deop"] = func(args string) error {
		return removeAdmin(args)
	}
	cmds["/delete"] = func(args string) error {
		return deleteGame()
	}
	cmds["/user"] = func(args string) error {
		return user()
	}
}
