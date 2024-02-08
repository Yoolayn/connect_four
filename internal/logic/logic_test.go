package connect_logic

import (
	"reflect"
	"slices"
	"testing"
)

func assertRow(t *testing.T, got, want row) {
	t.Helper()

	if !slices.Equal(got, want) {
		t.Errorf("wanted %#v, but got %#v", want, got)
	}
}

func TestRow(t *testing.T) {
	t.Run("maker", func(t *testing.T) {
		got := makeRow()
		want := row{
			NewChecker(""),
			NewChecker(""),
			NewChecker(""),
			NewChecker(""),
			NewChecker(""),
			NewChecker(""),
		}

		if !slices.Equal(got, want) {
			t.Errorf("got %#v, but want %#v", got, want)
		}
	})

	t.Run("claim", func(t *testing.T) {
		got := makeRow()
		want := row{
			NewChecker("#ff0000"),
			NewChecker(""),
			NewChecker(""),
			NewChecker(""),
			NewChecker(""),
			NewChecker(""),
		}
		_ = got.claim(NewChecker("#ff0000"))

		assertRow(t, got, want)
	})

	t.Run("claim with existing", func(t *testing.T) {
		got := row{
			NewChecker("#ff0000"),
			NewChecker(""),
			NewChecker(""),
			NewChecker(""),
			NewChecker(""),
			NewChecker(""),
		}
		want := row{
			NewChecker("#ff0000"),
			NewChecker("#ff0000"),
			NewChecker(""),
			NewChecker(""),
			NewChecker(""),
			NewChecker(""),
		}
		_ = got.claim(NewChecker("#ff0000"))

		assertRow(t, got, want)
	})

	t.Run("claim full", func(t *testing.T) {
		r := row{
			NewChecker("#ff0000"),
			NewChecker("#ff0000"),
			NewChecker("#ff0000"),
			NewChecker("#ff0000"),
			NewChecker("#ff0000"),
			NewChecker("#ff0000"),
		}
		cpy := row{
			NewChecker("#ff0000"),
			NewChecker("#ff0000"),
			NewChecker("#ff0000"),
			NewChecker("#ff0000"),
			NewChecker("#ff0000"),
			NewChecker("#ff0000"),
		}
		got := r.claim(NewChecker("#ff0000"))

		if got {
			t.Errorf("wanted %#v, but got %#v", false, got)
		}

		if !slices.Equal(r, cpy) {
			t.Error("row wasn't supposed to be changed")
		}
	})

	t.Run("no win", func(t *testing.T) {
		r := makeRow()
		_, got := r.checkWin()

		if got {
			t.Errorf("there shouldn't be a winner in an empty row, but got %#v", got)
		}
	})

	t.Run("winnable", func(t *testing.T) {
		r := row{
			NewChecker(""),
			NewChecker("#ff0000"),
			NewChecker("#ff0000"),
			NewChecker("#ff0000"),
			NewChecker("#ff0000"),
			NewChecker(""),
		}
		winner, got := r.checkWin()

		if !got {
			t.Error("this row should have a winner")
		}

		if !winner.Equal(NewChecker("#ff0000")) {
			t.Error("red should have been a winner")
		}
	})

	t.Run("enemy block", func(t *testing.T) {
		r := row{
			NewChecker("#ff0000"),
			NewChecker("#ff0000"),
			NewChecker("#0000ff"),
			NewChecker("#ff0000"),
			NewChecker("#ff0000"),
			NewChecker("#ff0000"),
		}
		winner, got := r.checkWin()

		if got {
			t.Error("there shouldn't have been a winner")
		}

		if !winner.Equal(NewChecker("")) {
			t.Error("winner should be an empty Checker")
		}
	})

	t.Run("winnable with enemy", func(t *testing.T) {
		r := row{
			NewChecker("#ff0000"),
			NewChecker("#ff0000"),
			NewChecker("#0000ff"),
			NewChecker("#0000ff"),
			NewChecker("#0000ff"),
			NewChecker("#0000ff"),
		}
		winner, got := r.checkWin()

		if !got {
			t.Error("there should be a winner")
		}

		if !winner.Equal(NewChecker("#0000ff")) {
			t.Error("the winner should be blue")
		}
	})
}

func TestBoard(t *testing.T) {
	t.Run("maker", func(t *testing.T) {
		got := MakeBoard()
		want := Board{
			row{NewChecker(""), NewChecker(""), NewChecker(""), NewChecker(""), NewChecker(""), NewChecker("")},
			row{NewChecker(""), NewChecker(""), NewChecker(""), NewChecker(""), NewChecker(""), NewChecker("")},
			row{NewChecker(""), NewChecker(""), NewChecker(""), NewChecker(""), NewChecker(""), NewChecker("")},
			row{NewChecker(""), NewChecker(""), NewChecker(""), NewChecker(""), NewChecker(""), NewChecker("")},
			row{NewChecker(""), NewChecker(""), NewChecker(""), NewChecker(""), NewChecker(""), NewChecker("")},
			row{NewChecker(""), NewChecker(""), NewChecker(""), NewChecker(""), NewChecker(""), NewChecker("")},
			row{NewChecker(""), NewChecker(""), NewChecker(""), NewChecker(""), NewChecker(""), NewChecker("")},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("wanted %#v, but got %#v", want, got)
		}
	})

	t.Run("claim", func(t *testing.T) {
		got := MakeBoard()
		want := Board{
			row{NewChecker(""),      NewChecker(""), NewChecker(""), NewChecker(""), NewChecker(""), NewChecker("")},
			row{NewChecker(""),      NewChecker(""), NewChecker(""), NewChecker(""), NewChecker(""), NewChecker("")},
			row{NewChecker("#ff0000"), NewChecker(""), NewChecker(""), NewChecker(""), NewChecker(""), NewChecker("")},
			row{NewChecker(""),      NewChecker(""), NewChecker(""), NewChecker(""), NewChecker(""), NewChecker("")},
			row{NewChecker(""),      NewChecker(""), NewChecker(""), NewChecker(""), NewChecker(""), NewChecker("")},
			row{NewChecker(""),      NewChecker(""), NewChecker(""), NewChecker(""), NewChecker(""), NewChecker("")},
			row{NewChecker(""),      NewChecker(""), NewChecker(""), NewChecker(""), NewChecker(""), NewChecker("")},
		}
		ok := got.Claim(NewChecker("#ff0000"), 2)

		if !ok {
			t.Error("tried to claim a field, but failed")
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("ended up with %#v, but wanted %#v", got, want)
		}
	})

	t.Run("no win", func(t *testing.T) {
		r := MakeBoard()

		winner, won := r.CheckWin()

		if won {
			t.Error("there wasn't supposed to be a winner")
		}

		if !winner.Equal(NewChecker("")) {
			t.Error("no one wins here")
		}
	})

	t.Run("row win", func(t *testing.T) {
		r := Board{
			row{NewChecker(""), NewChecker("#ff0000"), NewChecker("#ff0000"), NewChecker("#ff0000"), NewChecker("#ff0000"), NewChecker("")},
			row{NewChecker(""), NewChecker(""), NewChecker(""), NewChecker(""), NewChecker(""), NewChecker("")},
			row{NewChecker(""), NewChecker(""), NewChecker(""), NewChecker(""), NewChecker(""), NewChecker("")},
			row{NewChecker(""), NewChecker(""), NewChecker(""), NewChecker(""), NewChecker(""), NewChecker("")},
			row{NewChecker(""), NewChecker(""), NewChecker(""), NewChecker(""), NewChecker(""), NewChecker("")},
			row{NewChecker(""), NewChecker(""), NewChecker(""), NewChecker(""), NewChecker(""), NewChecker("")},
			row{NewChecker(""), NewChecker(""), NewChecker(""), NewChecker(""), NewChecker(""), NewChecker("")},
		}

		winner, won := r.CheckWin()

		if !won {
			t.Error("there was supposed to be a winner")
		}

		if !winner.Equal(NewChecker("#ff0000")) {
			t.Error("#ff0000 was supposed to win")
		}
	})

	t.Run("column win", func(t *testing.T) {
		r := Board{
			row{NewChecker(""), NewChecker(""), NewChecker(""), NewChecker(""), NewChecker(""), NewChecker("")},
			row{NewChecker("#ff0000"), NewChecker(""), NewChecker(""), NewChecker(""), NewChecker(""), NewChecker("")},
			row{NewChecker("#ff0000"), NewChecker(""), NewChecker(""), NewChecker(""), NewChecker(""), NewChecker("")},
			row{NewChecker("#ff0000"), NewChecker(""), NewChecker(""), NewChecker(""), NewChecker(""), NewChecker("")},
			row{NewChecker("#ff0000"), NewChecker(""), NewChecker(""), NewChecker(""), NewChecker(""), NewChecker("")},
			row{NewChecker("#ff0000"), NewChecker(""), NewChecker(""), NewChecker(""), NewChecker(""), NewChecker("")},
			row{NewChecker(""), NewChecker(""), NewChecker(""), NewChecker(""), NewChecker(""), NewChecker("")},
		}

		winner, won := r.CheckWin()

		if !won {
			t.Error("there was supposed to be a winner")
		}

		if !winner.Equal(NewChecker("#ff0000")) {
			t.Error("red was supposed to win this")
		}
	})

	t.Run("diagonal win", func(t *testing.T) {
		r := Board{
			row{NewChecker(""), NewChecker(""),      NewChecker(""),      NewChecker(""),      NewChecker(""),      NewChecker("")},
			row{NewChecker(""), NewChecker("#ff0000"), NewChecker(""),      NewChecker(""),      NewChecker(""),      NewChecker("")},
			row{NewChecker(""), NewChecker(""),      NewChecker("#ff0000"), NewChecker(""),      NewChecker(""),      NewChecker("")},
			row{NewChecker(""), NewChecker(""),      NewChecker(""),      NewChecker("#ff0000"), NewChecker(""),      NewChecker("")},
			row{NewChecker(""), NewChecker(""),      NewChecker(""),      NewChecker(""),      NewChecker("#ff0000"), NewChecker("")},
			row{NewChecker(""), NewChecker(""),      NewChecker(""),      NewChecker(""),      NewChecker(""),      NewChecker("")},
			row{NewChecker(""), NewChecker(""),      NewChecker(""),      NewChecker(""),      NewChecker(""),      NewChecker("")},
		}

		winner, won := r.CheckWin()

		if !won {
			t.Error("there was supposed to be a winner")
		}

		if !winner.Equal(NewChecker("#ff0000")) {
			t.Error("red was supposed to win this")
		}
	})

	t.Run("diagonal win opposite", func(t *testing.T) {
		r := Board{
			row{NewChecker(""), NewChecker(""),        NewChecker(""),        NewChecker(""),        NewChecker(""),        NewChecker("")},
			row{NewChecker(""), NewChecker(""),        NewChecker(""),        NewChecker(""),        NewChecker("#ff0000"), NewChecker("")},
			row{NewChecker(""), NewChecker(""),        NewChecker(""),        NewChecker("#ff0000"), NewChecker(""),        NewChecker("")},
			row{NewChecker(""), NewChecker(""),        NewChecker("#ff0000"), NewChecker(""),        NewChecker(""),        NewChecker("")},
			row{NewChecker(""), NewChecker("#ff0000"), NewChecker(""),        NewChecker(""),        NewChecker(""),        NewChecker("")},
			row{NewChecker(""), NewChecker(""),        NewChecker(""),        NewChecker(""),        NewChecker(""),        NewChecker("")},
			row{NewChecker(""), NewChecker(""),        NewChecker(""),        NewChecker(""),        NewChecker(""),        NewChecker("")},
		}

		winner, won := r.CheckWin()

		if !won {
			t.Error("there was supposed to be a winner")
		}

		if !winner.Equal(NewChecker("#ff0000")) {
			t.Error("red was supposed to win this")
		}
	})
}
