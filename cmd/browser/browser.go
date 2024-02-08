package main

import (
	"fmt"

	"github.com/google/uuid"
	"gitlab.com/Yoolayn/connect_four/internal/browser"
)

func main() {
	items := []string{uuid.NewString()+ ": game1", uuid.NewString()+ ": game2", uuid.NewString()+ ": game3"}
	selected, err := browser.New("games", items)
	if err != nil {
		panic(err)
	}
	fmt.Println(selected)
}
