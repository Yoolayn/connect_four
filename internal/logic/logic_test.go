package connect_logic

import (
	"slices"
	"testing"
)

func TestErr(t *testing.T) {
	t.Run("")
}

func TestRow(t *testing.T) {
	t.Run("make default row", func(t *testing.T) {
		got := makeRow()
		want := Row{
			Checker{},
			Checker{},
			Checker{},
			Checker{},
			Checker{},
			Checker{},
		}
		if !slices.Equal(got, want) {
			t.Errorf("got %#v, want %#v", got, want)
		}
	})
	t.Run("claim on empty", func(t *testing.T) {
		row := makeRow()
		checker := Checker{"test", "red"}
		err := row.Claim(checker)
		if err != nil {
			t.Errorf("got an error %q", err)
		}
	})
	t.Run("claim on full", func(t *testing.T) {
		enemyCh := Checker{"enemyCh", "Red"}
		row := Row{
			enemyCh,
			enemyCh,
			enemyCh,
			enemyCh,
			enemyCh,
			enemyCh,
		}
		err := row.Claim(Checker{"test", "red"})
		if err != ErrRowFull {
			t.Errorf("expected: %q", ErrRowFull)
		}
	})
}
