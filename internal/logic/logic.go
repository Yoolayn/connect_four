package connect_logic

type Checker struct {
	Color string
}

func (c Checker) Equal(ch Checker) bool {
	return c.Color == ch.Color
}

type row []Checker

func (r *row) claim(c Checker) bool {
	for i, v := range *r {
		if v.Color == "" {
			(*r)[i] = c
			return true
		}
	}
	return false
}

func makeRow() row {
	r := make(row, 6)
	for i := 0; i < 6; i++ {
		r[i] = Checker{}
	}
	return r
}

func (r row) checkWin() (Checker, bool) {
	counter := 0
	var previous Checker
	for _, v := range r {
		if previous.Equal(v) && previous.Color != "" {
			counter++
		} else {
			counter = 1
		}

		previous = v

		if counter >= 4 {
			return previous, true
		}
	}
	return Checker{}, false
}

type Board []row

func MakeBoard() Board {
	b := make(Board, 7)
	for i := 0; i < 7; i++ {
		b[i] = makeRow()
	}
	return b
}

func (b *Board) Claim(c Checker, r int) bool {
	ok := (*b)[r].claim(c)
	return ok
}

func (b Board) CheckWin() (Checker, bool) {
	for _, v := range b {
		if ch, won := v.checkWin(); won {
			return ch, won
		}
	}

	for i := range b[0] {
		counter := 0
		var previous Checker
		for j := range b {
			if previous.Equal(b[j][i]) && previous.Color != "" {
				counter++
			} else {
				counter = 1
			}
			previous = b[j][i]
			if counter >= 4 {
				return previous, true
			}
		}
	}

	for i := 0; i <= len(b)-4; i++ {
		for j := 0; j <= len(b[i])-4; j++ {
			if b[i][j].Color != "" &&
				b[i][j].Equal(b[i+1][j+1]) &&
				b[i+1][j+1].Equal(b[i+2][j+2]) &&
				b[i+2][j+2].Equal(b[i+3][j+3]) {
				return b[i][j], true
			}
		}
	}

	for i := 0; i <= len(b)-4; i++ {
		for j := len(b[i]) - 1; j >= 3; j-- {
			if b[i][j].Color != "" &&
				b[i][j].Equal(b[i+1][j-1]) &&
				b[i+1][j-1].Equal(b[i+2][j-2]) &&
				b[i+2][j-2].Equal(b[i+3][j-3]) {
				return b[i][j], true
			}
		}
	}

	return Checker{}, false
}
