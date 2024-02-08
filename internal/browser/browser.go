package browser

import (
	"errors"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

var (
	ErrNothingChosen = errors.New("no entry was selected, aborting...")
	height           int
	width            int
)

var style = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("63")).
	Padding(5).
	PaddingTop(2).
	PaddingBottom(2)

type Model struct {
	items    []string
	cursor   int
	selected string
	error    error
	title    string
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
	str := []string{m.title, "", ""}

	for i, v := range m.items {
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}

		str = append(str, fmt.Sprintf("%s %s", cursor, v))
	}

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, style.Render(
		lipgloss.JoinVertical(lipgloss.Center, str...),
	))
}

func New(title string, items []string) (string, error) {
	var err error
	width, height, err = term.GetSize(0)
	if err != nil {
		return "", err
	}

	p := tea.NewProgram(Model{
		items:    items,
		cursor:   0,
		error:    nil,
		selected: "",
		title:    title,
	}, tea.WithAltScreen())
	m, err := p.Run()
	if err != nil {
		return "", err
	}

	if m.(Model).error != nil {
		return "", m.(Model).error
	}

	return m.(Model).selected, nil
}
