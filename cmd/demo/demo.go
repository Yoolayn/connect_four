package main

import (
	"fmt"
	"time"

	connect_logic "gitlab.com/Yoolayn/connect_four/internal/logic"
)

func main() {
	b := connect_logic.MakeBoard()

	red := connect_logic.Checker{Color: "red"}
	blue := connect_logic.Checker{Color: "blue"}

	for i := 0; i < 4; i++ {
		time.Sleep(500 * time.Millisecond)
		win := func() bool {
			winner, won := b.CheckWin()
			if won {
				fmt.Println("the", winner.Color, "has won!")
			}
			return won
		}

		b.Claim(red, 0)
		fmt.Println("the", red.Color, "is making a move and takes claims a checker on row 0!")
		if win() {
			break
		}

		time.Sleep(500 * time.Millisecond)

		b.Claim(blue, 1)
		fmt.Println("the", blue.Color, "is making a move and takes claims a checker on row 1!")
		if win() {
			break
		}
	}
}
