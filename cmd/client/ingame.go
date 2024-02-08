package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	logic "gitlab.com/Yoolayn/connect_four/internal/logic"
	"golang.org/x/term"
)

type Model struct {
	game   Game
	choice int
	player Player
	width  int
	height int
	error  error
}

func InitialModel(p int, g Game) (Model, error) {
	width, height, err := term.GetSize(0)
	if err != nil {
		return Model{}, err
	}
	return Model{
		game:   g,
		choice: 0,
		player: func() Player {
			if p == 1 {
				return g.Player1
			} else {
				return g.Player2
			}
		}(),
		width:  width,
		height: height,
	}, nil
}

var helpInGame = lipgloss.JoinVertical(lipgloss.Center, []string{"0, 1, 2, 3, 4, 5, 6, 7 - make a move", "enter - confirm a move", "r - refresh; q - quit"}...)

func (m Model) View() string {
	b := m.game.Board.Clone()
	if m.choice != 0 {
		b.Claim(logic.NewChecker("#9900ff"), m.choice-1)
	}
	bStr := b.ToTable()
	var err string
	if m.error != nil {
		err = m.error.Error()
	}
	bAndHelp := lipgloss.JoinVertical(lipgloss.Center, bStr, helpInGame)
	info := lipgloss.JoinVertical(lipgloss.Center, func() []string {
		p1c := func() string {
			color := m.game.Player1.Color
			if color == "#ff0000" {
				return "red"
			} else if color == "#ffff00" {
				return "yellow"
			}
			return ""
		}()
		p2c := func() string {
			color := m.game.Player2.Color
			if color == "#ff0000" {
				return "red"
			} else if color == "#ffff00" {
				return "yellow"
			}
			return ""
		}()
		if m.game.Player1.User.Login == m.player.User.Login {
			return []string{"You:", m.game.Player1.User.String(), p1c, "", "", m.game.Player2.User.String(), p2c}
		} else {
			return []string{"", m.game.Player1.User.String(), "", "You:", m.game.Player2.User.String()}
		}
	}()...)
	finishedProduct := lipgloss.JoinHorizontal(lipgloss.Center, bAndHelp, info)
	finishedForSureNow := lipgloss.JoinVertical(lipgloss.Center, "InGame", err, finishedProduct)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, finishedForSureNow)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "0":
			m.choice = 0
		case "1":
			m.choice = 1
		case "2":
			m.choice = 2
		case "3":
			m.choice = 3
		case "4":
			m.choice = 4
		case "5":
			m.choice = 5
		case "6":
			m.choice = 6
		case "7":
			m.choice = 7
		}
	}
	return m, nil
}

func (m Model) Init() tea.Cmd {
	return nil
}
