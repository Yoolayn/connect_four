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
			Checker{},
			Checker{},
			Checker{},
			Checker{},
			Checker{},
			Checker{},
		}

		if !slices.Equal(got, want) {
			t.Errorf("got %#v, but want %#v", got, want)
		}
	})

	t.Run("claim", func(t *testing.T) {
		got := makeRow()
		want := row{
			Checker{"red"},
			Checker{},
			Checker{},
			Checker{},
			Checker{},
			Checker{},
		}
		_ = got.claim(Checker{"red"})

		assertRow(t, got, want)
	})

	t.Run("claim with existing", func(t *testing.T) {
		got := row{
			Checker{"red"},
			Checker{},
			Checker{},
			Checker{},
			Checker{},
			Checker{},
		}
		want := row{
			Checker{"red"},
			Checker{"red"},
			Checker{},
			Checker{},
			Checker{},
			Checker{},
		}
		_ = got.claim(Checker{"red"})

		assertRow(t, got, want)
	})

	t.Run("claim full", func(t *testing.T) {
		r := row{
			Checker{"red"},
			Checker{"red"},
			Checker{"red"},
			Checker{"red"},
			Checker{"red"},
			Checker{"red"},
		}
		cpy := row{
			Checker{"red"},
			Checker{"red"},
			Checker{"red"},
			Checker{"red"},
			Checker{"red"},
			Checker{"red"},
		}
		got := r.claim(Checker{"red"})

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
			Checker{},
			Checker{"red"},
			Checker{"red"},
			Checker{"red"},
			Checker{"red"},
			Checker{},
		}
		winner, got := r.checkWin()

		if !got {
			t.Error("this row should have a winner")
		}

		if !winner.Equal(Checker{"red"}) {
			t.Error("red should have been a winner")
		}
	})

	t.Run("enemy block", func(t *testing.T) {
		r := row{
			Checker{"red"},
			Checker{"red"},
			Checker{"blue"},
			Checker{"red"},
			Checker{"red"},
			Checker{"red"},
		}
		winner, got := r.checkWin()

		if got {
			t.Error("there shouldn't have been a winner")
		}

		if !winner.Equal(Checker{}) {
			t.Error("winner should be an empty Checker")
		}
	})

	t.Run("winnable with enemy", func(t *testing.T) {
		r := row{
			Checker{"red"},
			Checker{"red"},
			Checker{"blue"},
			Checker{"blue"},
			Checker{"blue"},
			Checker{"blue"},
		}
		winner, got := r.checkWin()

		if !got {
			t.Error("there should be a winner")
		}

		if !winner.Equal(Checker{"blue"}) {
			t.Error("the winner should be blue")
		}
	})
}

func TestBoard(t *testing.T) {
	t.Run("maker", func(t *testing.T) {
		got := MakeBoard()
		want := Board{
			row{Checker{}, Checker{}, Checker{}, Checker{}, Checker{}, Checker{}},
			row{Checker{}, Checker{}, Checker{}, Checker{}, Checker{}, Checker{}},
			row{Checker{}, Checker{}, Checker{}, Checker{}, Checker{}, Checker{}},
			row{Checker{}, Checker{}, Checker{}, Checker{}, Checker{}, Checker{}},
			row{Checker{}, Checker{}, Checker{}, Checker{}, Checker{}, Checker{}},
			row{Checker{}, Checker{}, Checker{}, Checker{}, Checker{}, Checker{}},
			row{Checker{}, Checker{}, Checker{}, Checker{}, Checker{}, Checker{}},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("wanted %#v, but got %#v", want, got)
		}
	})

	t.Run("claim", func(t *testing.T) {
		got := MakeBoard()
		want := Board{
			row{Checker{},      Checker{}, Checker{}, Checker{}, Checker{}, Checker{}},
			row{Checker{},      Checker{}, Checker{}, Checker{}, Checker{}, Checker{}},
			row{Checker{"red"}, Checker{}, Checker{}, Checker{}, Checker{}, Checker{}},
			row{Checker{},      Checker{}, Checker{}, Checker{}, Checker{}, Checker{}},
			row{Checker{},      Checker{}, Checker{}, Checker{}, Checker{}, Checker{}},
			row{Checker{},      Checker{}, Checker{}, Checker{}, Checker{}, Checker{}},
			row{Checker{},      Checker{}, Checker{}, Checker{}, Checker{}, Checker{}},
		}
		ok := got.Claim(Checker{"red"}, 2)

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

		if !winner.Equal(Checker{}) {
			t.Error("no one wins here")
		}
	})

	t.Run("row win", func(t *testing.T) {
		r := Board{
			row{Checker{}, Checker{"red"}, Checker{"red"}, Checker{"red"}, Checker{"red"}, Checker{}},
			row{Checker{}, Checker{}, Checker{}, Checker{}, Checker{}, Checker{}},
			row{Checker{}, Checker{}, Checker{}, Checker{}, Checker{}, Checker{}},
			row{Checker{}, Checker{}, Checker{}, Checker{}, Checker{}, Checker{}},
			row{Checker{}, Checker{}, Checker{}, Checker{}, Checker{}, Checker{}},
			row{Checker{}, Checker{}, Checker{}, Checker{}, Checker{}, Checker{}},
			row{Checker{}, Checker{}, Checker{}, Checker{}, Checker{}, Checker{}},
		}

		winner, won := r.CheckWin()

		if !won {
			t.Error("there was supposed to be a winner")
		}

		if !winner.Equal(Checker{"red"}) {
			t.Error("red was supposed to win")
		}
	})

	t.Run("column win", func(t *testing.T) {
		r := Board{
			row{Checker{}, Checker{}, Checker{}, Checker{}, Checker{}, Checker{}},
			row{Checker{"red"}, Checker{}, Checker{}, Checker{}, Checker{}, Checker{}},
			row{Checker{"red"}, Checker{}, Checker{}, Checker{}, Checker{}, Checker{}},
			row{Checker{"red"}, Checker{}, Checker{}, Checker{}, Checker{}, Checker{}},
			row{Checker{"red"}, Checker{}, Checker{}, Checker{}, Checker{}, Checker{}},
			row{Checker{"red"}, Checker{}, Checker{}, Checker{}, Checker{}, Checker{}},
			row{Checker{}, Checker{}, Checker{}, Checker{}, Checker{}, Checker{}},
		}

		winner, won := r.CheckWin()

		if !won {
			t.Error("there was supposed to be a winner")
		}

		if !winner.Equal(Checker{"red"}) {
			t.Error("red was supposed to win this")
		}
	})

	t.Run("diagonal win", func(t *testing.T) {
		r := Board{
			row{Checker{}, Checker{},      Checker{},      Checker{},      Checker{},      Checker{}},
			row{Checker{}, Checker{"red"}, Checker{},      Checker{},      Checker{},      Checker{}},
			row{Checker{}, Checker{},      Checker{"red"}, Checker{},      Checker{},      Checker{}},
			row{Checker{}, Checker{},      Checker{},      Checker{"red"}, Checker{},      Checker{}},
			row{Checker{}, Checker{},      Checker{},      Checker{},      Checker{"red"}, Checker{}},
			row{Checker{}, Checker{},      Checker{},      Checker{},      Checker{},      Checker{}},
			row{Checker{}, Checker{},      Checker{},      Checker{},      Checker{},      Checker{}},
		}

		winner, won := r.CheckWin()

		if !won {
			t.Error("there was supposed to be a winner")
		}

		if !winner.Equal(Checker{"red"}) {
			t.Error("red was supposed to win this")
		}
	})

	t.Run("diagonal win opposite", func(t *testing.T) {
		r := Board{
			row{Checker{}, Checker{},      Checker{},      Checker{},      Checker{},      Checker{}},
			row{Checker{}, Checker{},      Checker{},      Checker{},      Checker{"red"}, Checker{}},
			row{Checker{}, Checker{},      Checker{},      Checker{"red"}, Checker{},      Checker{}},
			row{Checker{}, Checker{},      Checker{"red"}, Checker{},      Checker{},      Checker{}},
			row{Checker{}, Checker{"red"}, Checker{},      Checker{},      Checker{},      Checker{}},
			row{Checker{}, Checker{},      Checker{},      Checker{},      Checker{},      Checker{}},
			row{Checker{}, Checker{},      Checker{},      Checker{},      Checker{},      Checker{}},
		}

		winner, won := r.CheckWin()

		if !won {
			t.Error("there was supposed to be a winner")
		}

		if !winner.Equal(Checker{"red"}) {
			t.Error("red was supposed to win this")
		}
	})
}
