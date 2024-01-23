package connect_logic

type Checker struct {
	Owner string
	Color string
}

type CheckerErr string

func (c CheckerErr) Error() string {
	return string(c)
}

const ErrRowFull = CheckerErr("cannot claim a field in a row, it's full")

type Row []Checker

func (r *Row) Claim(c Checker) error {
	for i, v := range *r {
		if v.Color == "" && v.Owner == "" {
			(*r)[i] = c
			return nil
		}
	}
	return ErrRowFull
}

type Board []Row

func makeRow() Row {
	r := make(Row, 6)
	for i := 0; i < 6; i++ {
		r[i] = Checker{}
	}
	return r
}

func MakeBoard() Board {
	b := make(Board, 7)
	for i := 0; i < 7; i++ {
		b[i] = makeRow()
	}
	return b
}
