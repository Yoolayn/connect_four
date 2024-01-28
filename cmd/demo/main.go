package main

import (
	"fmt"

	connect_logic "gitlab.com/Yoolayn/connect_four/internal/logic"
)

func main() {
	b := connect_logic.MakeBoard()

	red := connect_logic.Checker{Color: "red"}
	blue := connect_logic.Checker{Color: "blue"}

	for i := 0; i < 4; i++ {
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

		b.Claim(blue, 1)
		fmt.Println("the", blue.Color, "is making a move and takes claims a checker on row 1!")
		if win() {
			break
		}
	}
}
