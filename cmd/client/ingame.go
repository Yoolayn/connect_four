package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
	logic "gitlab.com/Yoolayn/connect_four/internal/logic"
	"golang.org/x/term"
)

type Model struct {
	game   Game
	choice int
	player Player
	width  int
	height int
	uuid   uuid.UUID
	error  error
	input  textinput.Model
}

func getter(m Model) (tea.Model, tea.Cmd) {
	res, err := http.Get(baseurl("/games/" + m.uuid.String()))
	if err != nil {
		m.error = err
		return m, tea.Quit
	}

	switch res.StatusCode {
	case http.StatusNotFound:
		m.error = ErrGameNotFound
		return m, tea.Quit
	case http.StatusOK:
	default:
		m.error = ErrUnknown
		return m, tea.Quit
	}

	var game Game
	bities, err := io.ReadAll(res.Body)
	if err != nil {
		m.error = err
		return m, tea.Quit
	}

	err = json.Unmarshal(bities, &game)
	if err != nil {
		m.error = err
		return m, tea.Quit
	}

	m.game = game
	m.error = nil
	return m, nil
}

func InitialModel(p int, g Game, u uuid.UUID) (Model, error) {
	width, height, err := term.GetSize(0)
	if err != nil {
		return Model{}, err
	}
	ti := textinput.New()
	ti.Prompt = "$ "
	ti.Placeholder = g.Title
	ti.CharLimit = 60
	ti.Width = 30
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
		uuid:   u,
		error:  err,
		input:  ti,
	}, nil
}

var helpInGame = lipgloss.JoinVertical(
	lipgloss.Center,
	[]string{
		"0, 1, 2, 3, 4, 5, 6, 7 - make a move",
		"enter - confirm a move",
		"c - change name; m - toggle mqtt",
		"d - delete game",
		"r - refresh; q - quit",
	}...,
)

func (m Model) View() string {
	if m.game.Finished {
		style := lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#9900ff")).
			PaddingTop(5).
			PaddingBottom(5).
			PaddingLeft(10).
			PaddingRight(10)
		winner := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00ff00"))
		loser := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff0000"))
		if m.game.Winner.Color == m.player.Color {
			msg := style.Render(winner.Render("YOU WON!!!!!"))
			b := m.game.Board.Clone().ToTable()
			joined := lipgloss.JoinHorizontal(lipgloss.Center, b, msg)
			return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, joined)
		} else {
			msg := style.Render(loser.Render("you've lost, womp womp"))
			b := m.game.Board.Clone().ToTable()
			joined := lipgloss.JoinHorizontal(lipgloss.Center, b, msg)
			return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, joined)
		}
	}
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
			return []string{"You:", m.game.Player1.User.String(), p1c, "", m.game.Player2.User.String(), p2c, "", "", "", ""}
		} else {
			return []string{"", m.game.Player1.User.String(), "", "You:", m.game.Player2.User.String()}
		}
	}()...)
	finishedProduct := lipgloss.JoinHorizontal(lipgloss.Center, bAndHelp, info)
	finishedForSureNow := lipgloss.JoinVertical(lipgloss.Center, "InGame: "+m.game.Title, err, finishedProduct)

	if m.input.Focused() {
		finishedForSureNow = lipgloss.JoinVertical(lipgloss.Center, finishedForSureNow, "", m.input.View())
	}

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, finishedForSureNow)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.input.Focused() {
		var cmd tea.Cmd
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", tea.KeyEsc.String():
				m.input.Blur()
				m.input.SetValue("")
				return m, nil
			case "enter":
				newName := m.input.Value()
				m.input.Blur()
				m.input.SetValue("")

				payload := struct {
					C     credentials `json:"credentials"`
					Title string      `json:"title"`
				}{
					C:     creds,
					Title: newName,
				}

				bitties, err := json.Marshal(payload)
				if err != nil {
					m.error = err
					return m, nil
				}

				r, err := http.NewRequest(http.MethodPut, baseurl("/games/"+m.uuid.String()), bytes.NewBuffer(bitties))
				if err != nil {
					m.error = err
					return m, nil
				}

				c := &http.Client{}
				res, err := c.Do(r)
				if err != nil {
					m.error = err
					return m, nil
				}

				switch res.StatusCode {
				case http.StatusOK:
					fallthrough
				case http.StatusAccepted:
					model, cmd := getter(m)
					casted := model.(Model)
					return casted, cmd
				case http.StatusInternalServerError:
					m.error = ErrUnknown
					return m, tea.Quit
				case http.StatusNotFound:
					m.error = ErrUserNotFound
					return m, tea.Quit
				case http.StatusUnauthorized:
					m.error = ErrNotLoggedIn
					return m, tea.Quit
				case http.StatusForbidden:
					m.error = ErrForbidden
					return m, tea.Quit
				default:
					m.error = ErrUnknown
					return m, nil
				}
			}
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
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
		case "h":
			if m.choice > 0 {
				m.choice--
			}
		case "l":
			if m.choice < len(m.game.Board) {
				m.choice++
			}
		case "q":

			json, err := json.Marshal(creds)
			if err != nil {
				m.error = err
				return m, tea.Quit
			}

			r, err := http.NewRequest(http.MethodDelete, baseurl("/games/"+m.uuid.String()+"/leave"), bytes.NewBuffer(json))
			if err != nil {
				m.error = err
				return m, tea.Quit
			}

			c := &http.Client{}
			res, err := c.Do(r)
			if err != nil {
				m.error = err
				return m, tea.Quit
			}

			switch res.StatusCode {
			case http.StatusOK:
				m.error = nil
				return m, tea.Quit
			case http.StatusInternalServerError:
				m.error = ErrUnknown
				return m, tea.Quit
			case http.StatusNotFound:
				m.error = ErrUserNotFound
				return m, tea.Quit
			case http.StatusUnauthorized:
				m.error = ErrNotLoggedIn
				return m, tea.Quit
			case http.StatusForbidden:
				m.error = ErrForbidden
				return m, tea.Quit
			default:
				m.error = ErrUnknown
				return m, tea.Quit
			}
		case "enter":
			payload := struct {
				Credentials credentials `json:"credentials"`
				Row         int         `json:"row"`
			}{
				Credentials: creds,
				Row:         m.choice - 1,
			}

			jsn, err := json.Marshal(payload)
			if err != nil {
				m.error = err
				return m, tea.Quit
			}

			r, err := http.NewRequest(http.MethodPut, baseurl("/games/"+m.uuid.String()+"/move"), bytes.NewBuffer(jsn))
			if err != nil {
				m.error = err
				return m, tea.Quit
			}

			c := &http.Client{}
			res, err := c.Do(r)
			if err != nil {
				m.error = err
				return m, tea.Quit
			}
			switch res.StatusCode {
			case http.StatusOK:
				fallthrough
			case http.StatusAccepted:
				model, cmd := getter(m)
				casted := model.(Model)
				casted.choice = 0
				return casted, cmd
			case http.StatusInternalServerError:
				m.error = ErrUnknown
				return m, nil
			case http.StatusNotFound:
				m.error = ErrUserNotFound
				return m, tea.Quit
			case http.StatusUnauthorized:
				m.error = ErrNotLoggedIn
				return m, tea.Quit
			case http.StatusForbidden:
				m.error = ErrForbidden
				return m, tea.Quit
			case http.StatusBadRequest:
				m.error = ErrRowTaken
				return m, nil
			default:
				m.error = ErrUnknown
				return m, nil
			}
		case "r":
			model, cmd := getter(m)
			casted := model.(Model)
			casted.choice = 0
			return casted, cmd
		case "c":
			m.input.Focus()
			return m, textinput.Blink
		case "m":
			if m.error != nil {
				m.error = nil
			} else {
				m.error = ErrNotImplemented
			}
		case "d":
			payload := struct {
				Login    string `json:"login"`
				Password string `json:"password"`
			}{
				Login:    creds.Login,
				Password: creds.Password,
			}

			bitties, err := json.Marshal(payload)
			if err != nil {
				m.error = err
				return m, nil
			}

			req, err := http.NewRequest(http.MethodDelete, baseurl("/games/"+m.uuid.String()), bytes.NewBuffer(bitties))
			if err != nil {
				m.error = err
				return m, nil
			}

			c := &http.Client{}
			res, err := c.Do(req)
			if err != nil {
				m.error = err
				return m, nil
			}

			switch res.StatusCode {
			case http.StatusOK:
				return m, tea.Quit
			case http.StatusForbidden:
				m.error = ErrAdminRequired
				return m, nil
			default:
				m.error = ErrRequest
				return m, nil
			}
		}
	}
	return m, nil
}

func (m Model) Init() tea.Cmd {
	return nil
}
