package main

import (
	"battleship/game"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Model represents the bubbletea model for the game
type Model struct {
	game               *game.Game
	cursorRow          int
	cursorCol          int
	shipOrientation    game.Orientation
	showHelp           bool
	computerThinking   bool
	width              int
	height             int
}

// computerTurnMsg is sent after a delay to simulate computer thinking
type computerTurnMsg struct{}

func computerTurn() tea.Msg {
	time.Sleep(800 * time.Millisecond)
	return computerTurnMsg{}
}

// InitialModel creates the initial model
func InitialModel() Model {
	return Model{
		game:            game.NewGame(10),
		cursorRow:       0,
		cursorCol:       0,
		shipOrientation: game.Horizontal,
		showHelp:        true,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case computerTurnMsg:
		if m.game.Phase == game.ComputerTurnPhase {
			m.game.ComputerAttack()
			m.computerThinking = false
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "h":
			m.showHelp = !m.showHelp
			return m, nil

		case "r":
			// Reset game
			m.game = game.NewGame(10)
			m.cursorRow = 0
			m.cursorCol = 0
			m.shipOrientation = game.Horizontal
			m.showHelp = true
			m.computerThinking = false
			return m, nil

		case "up", "w":
			if m.cursorRow > 0 {
				m.cursorRow--
			}
			return m, nil

		case "down", "s":
			if m.cursorRow < m.game.BoardSize-1 {
				m.cursorRow++
			}
			return m, nil

		case "left", "a":
			if m.cursorCol > 0 {
				m.cursorCol--
			}
			return m, nil

		case "right", "d":
			if m.cursorCol < m.game.BoardSize-1 {
				m.cursorCol++
			}
			return m, nil

		case "o", "O":
			// Toggle ship orientation during placement
			if m.game.Phase == game.PlacementPhase {
				if m.shipOrientation == game.Horizontal {
					m.shipOrientation = game.Vertical
				} else {
					m.shipOrientation = game.Horizontal
				}
			}
			return m, nil

		case " ", "enter":
			return m.handleAction()
		}
	}

	return m, nil
}

// handleAction handles the action button (space/enter)
func (m Model) handleAction() (tea.Model, tea.Cmd) {
	pos := game.Position{Row: m.cursorRow, Col: m.cursorCol}

	switch m.game.Phase {
	case game.PlacementPhase:
		m.game.PlacePlayerShip(pos, m.shipOrientation)
		return m, nil

	case game.PlayerTurnPhase:
		if m.game.PlayerAttack(pos) {
			if m.game.Phase == game.ComputerTurnPhase {
				m.computerThinking = true
				return m, computerTurn
			}
		}
		return m, nil

	case game.GameOverPhase:
		// Could restart on enter
		return m, nil
	}

	return m, nil
}

// View renders the model
func (m Model) View() string {
	return renderGame(m)
}
