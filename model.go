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
	menuSelection      int
	selectedDifficulty game.Difficulty
}

// computerTurnMsg is sent after a delay to simulate computer thinking
type computerTurnMsg struct{}

func computerTurn() tea.Msg {
	time.Sleep(800 * time.Millisecond)
	return computerTurnMsg{}
}

// InitialModel creates the initial model
func InitialModel() Model {
	g := game.NewGame(10)
	g.Phase = game.MainMenuPhase
	return Model{
		game:            g,
		cursorRow:       0,
		cursorCol:       0,
		shipOrientation: game.Horizontal,
		showHelp:        false,
		menuSelection:   0,
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
			if m.game.Phase == game.MainMenuPhase {
				if m.menuSelection > 0 {
					m.menuSelection--
				}
			} else if m.cursorRow > 0 {
				m.cursorRow--
			}
			return m, nil

		case "down", "s":
			if m.game.Phase == game.MainMenuPhase {
				if m.menuSelection < 2 { // 0 = Difficulty, 1 = Play, 2 = Quit
					m.menuSelection++
				}
			} else if m.cursorRow < m.game.BoardSize-1 {
				m.cursorRow++
			}
			return m, nil

		case "left", "a":
			if m.game.Phase == game.MainMenuPhase && m.menuSelection == 0 {
				// Cycle difficulty left
				if m.selectedDifficulty == game.Easy {
					m.selectedDifficulty = game.Hard
				} else {
					m.selectedDifficulty--
				}
			} else if m.cursorCol > 0 {
				m.cursorCol--
			}
			return m, nil

		case "right", "d":
			if m.game.Phase == game.MainMenuPhase && m.menuSelection == 0 {
				// Cycle difficulty right
				if m.selectedDifficulty == game.Hard {
					m.selectedDifficulty = game.Easy
				} else {
					m.selectedDifficulty++
				}
			} else if m.cursorCol < m.game.BoardSize-1 {
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
	case game.MainMenuPhase:
		if m.menuSelection == 0 {
			// Difficulty selection - do nothing, just cycle with arrow keys
			return m, nil
		} else if m.menuSelection == 1 {
			// Start new game
			m.game = game.NewGame(10)
			m.game.Difficulty = m.selectedDifficulty
			m.cursorRow = 0
			m.cursorCol = 0
			m.shipOrientation = game.Horizontal
			m.showHelp = true
			m.computerThinking = false
		} else {
			// Quit
			return m, tea.Quit
		}
		return m, nil

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
