package main

import (
	"fmt"

	connect_logic "gitlab.com/Yoolayn/connect_four/internal/logic"
)

func main() {
	b := connect_logic.MakeBoard()
	fmt.Println(b[0][0])
}
