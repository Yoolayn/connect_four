package main

import (
	"fmt"
	"time"

	"gitlab.com/Yoolayn/connect_four/internal/browser"
)

func noAlt() {
	items := []string{"game1", "game2", "game3"}
	selected, err := browser.New(items, false)
	if err != nil {
		panic(err)
	}
	fmt.Println(selected)
}

func alt() {
	items := []string{"game1", "game2", "game3"}
	selected, err := browser.New(items, true)
	if err != nil {
		panic(err)
	}
	fmt.Println(selected)
}

func main() {
	fmt.Println("version without alt screen:")
	time.Sleep(1*time.Second)
	noAlt()
	time.Sleep(1*time.Second)
	fmt.Println("version with alt screen:")
	time.Sleep(1*time.Second)
	alt()
}
