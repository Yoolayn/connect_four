package browser

import (
	"errors"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

var ErrNothingChosen = errors.New("no entry was selected, aborting...")

type Model struct {
	items    []string
	cursor   int
	selected string
	error    error
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.error = ErrNothingChosen
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case " ", "enter":
			m.selected = m.items[m.cursor]
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	str := "Games:\n\n"

	for i, v := range m.items {
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}

		str += fmt.Sprintf("%s %s\n", cursor, v)
	}

	return str
}

func New(items []string, altScreen bool) (string, error) {
	p := func() *tea.Program {
		m := Model{
			items:    items,
			cursor:   0,
			error:    nil,
			selected: "",
		}
		if altScreen {
			return tea.NewProgram(m, tea.WithAltScreen())
		} else {
			return tea.NewProgram(m)
		}
	}()
	m, err := p.Run()
	if err != nil {
		return "", err
	}

	if m.(Model).error != nil {
		return "", m.(Model).error
	}

	return m.(Model).selected, nil
}
