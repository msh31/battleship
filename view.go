package main

import (
	"battleship/game"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors
	oceanBlue    = lipgloss.Color("#0066CC")
	darkBlue     = lipgloss.Color("#003366")
	shipGray     = lipgloss.Color("#666666")
	hitRed       = lipgloss.Color("#CC0000")
	missWhite    = lipgloss.Color("#AAAAAA")
	cursorYellow = lipgloss.Color("#FFCC00")
	successGreen = lipgloss.Color("#00CC00")
	titleCyan    = lipgloss.Color("#00CCCC")

	// Styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(titleCyan).
			Padding(1, 2).
			Border(lipgloss.DoubleBorder()).
			BorderForeground(titleCyan)

	boardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(oceanBlue).
			Padding(1, 2)

	cellStyle = lipgloss.NewStyle().
			Width(3).
			Align(lipgloss.Center)

	waterStyle = cellStyle.Copy().
			Foreground(oceanBlue).
			Background(darkBlue)

	shipStyle = cellStyle.Copy().
			Foreground(shipGray).
			Background(darkBlue).
			Bold(true)

	hitStyle = cellStyle.Copy().
			Foreground(hitRed).
			Background(darkBlue).
			Bold(true)

	missStyle = cellStyle.Copy().
			Foreground(missWhite).
			Background(darkBlue)

	cursorStyle = cellStyle.Copy().
			Foreground(cursorYellow).
			Background(darkBlue).
			Bold(true)

	grayCursorStyle = cellStyle.Copy().
			Foreground(lipgloss.Color("#888888")).
			Background(darkBlue).
			Bold(true)

	messageStyle = lipgloss.NewStyle().
			Foreground(successGreen).
			Bold(true).
			Padding(0, 2)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Padding(1, 2)

	headerStyle = lipgloss.NewStyle().
			Foreground(titleCyan).
			Bold(true).
			Padding(0, 2)
)

func renderGame(m Model) string {
	var sb strings.Builder

	// Title
	title := titleStyle.Render("⚓ BATTLESHIP ⚓")
	sb.WriteString(title)
	sb.WriteString("\n\n")

	// Game phase message
	sb.WriteString(renderPhaseMessage(m))
	sb.WriteString("\n")

	// Render boards side by side
	switch m.game.Phase {
	case game.PlacementPhase:
		sb.WriteString(renderPlacementBoard(m))
	case game.PlayerTurnPhase, game.ComputerTurnPhase:
		sb.WriteString(renderBattleBoards(m))
	case game.GameOverPhase:
		sb.WriteString(renderGameOver(m))
	}

	// Help text
	if m.showHelp {
		sb.WriteString("\n")
		sb.WriteString(renderHelp(m))
	}

	return sb.String()
}

func renderPhaseMessage(m Model) string {
	msg := ""

	switch m.game.Phase {
	case game.PlacementPhase:
		ship := m.game.GetCurrentShipForPlacement()
		if ship != nil {
			orientation := "Horizontal"
			if m.shipOrientation == game.Vertical {
				orientation = "Vertical"
			}
			msg = fmt.Sprintf("Place your %s (Length: %d) - Orientation: %s",
				ship.Name, ship.Length, orientation)
		}
	case game.PlayerTurnPhase:
		msg = "Your turn! Select a target and fire!"
	case game.ComputerTurnPhase:
		msg = "Claude is thinking..."
	case game.GameOverPhase:
		msg = fmt.Sprintf("Game Over! %s wins!", m.game.Winner)
	}

	if m.game.LastMessage != "" && m.game.Phase != game.PlacementPhase {
		msg = m.game.LastMessage
	}

	return messageStyle.Render(msg)
}

func renderPlacementBoard(m Model) string {
	var sb strings.Builder

	sb.WriteString(headerStyle.Render("Your Fleet"))
	sb.WriteString("\n\n")

	// Render column headers
	sb.WriteString("    ")
	for col := 0; col < m.game.BoardSize; col++ {
		sb.WriteString(fmt.Sprintf(" %c  ", 'A'+col))
	}
	sb.WriteString("\n")

	// Render board
	for row := 0; row < m.game.BoardSize; row++ {
		// Row number
		sb.WriteString(fmt.Sprintf("%2d  ", row+1))

		for col := 0; col < m.game.BoardSize; col++ {
			pos := game.Position{Row: row, Col: col}
			cell := m.game.PlayerBoard.GetCell(pos)

			// Check if this is a preview position for the current ship
			isPreview := false
			isCursor := row == m.cursorRow && col == m.cursorCol

			if isCursor {
				ship := m.game.GetCurrentShipForPlacement()
				if ship != nil {
					if m.game.PlayerBoard.CanPlaceShip(pos, ship.Length, m.shipOrientation) {
						// Show preview
						for i := 0; i < ship.Length; i++ {
							previewRow := row
							previewCol := col
							if m.shipOrientation == game.Horizontal {
								previewCol += i
							} else {
								previewRow += i
							}
							if previewRow == row && previewCol == col {
								isPreview = true
								break
							}
						}
					}
				}
			}

			// Check if any upcoming preview cell matches this position
			if !isPreview && !isCursor {
				ship := m.game.GetCurrentShipForPlacement()
				if ship != nil {
					cursorPos := game.Position{Row: m.cursorRow, Col: m.cursorCol}
					if m.game.PlayerBoard.CanPlaceShip(cursorPos, ship.Length, m.shipOrientation) {
						for i := 0; i < ship.Length; i++ {
							previewRow := m.cursorRow
							previewCol := m.cursorCol
							if m.shipOrientation == game.Horizontal {
								previewCol += i
							} else {
								previewRow += i
							}
							if previewRow == row && previewCol == col {
								isPreview = true
								break
							}
						}
					}
				}
			}

			cellStr := renderCell(cell, isCursor, isPreview, false)
			sb.WriteString(cellStr)
		}
		sb.WriteString("\n")
	}

	return boardStyle.Render(sb.String())
}

func renderBattleBoards(m Model) string {
	playerBoard := renderPlayerBoard(m)
	enemyBoard := renderEnemyBoard(m)

	return lipgloss.JoinHorizontal(lipgloss.Top, playerBoard, "  ", enemyBoard)
}

func renderPlayerBoard(m Model) string {
	var sb strings.Builder

	sb.WriteString(headerStyle.Render("Your Fleet"))
	sb.WriteString("\n\n")

	// Column headers
	sb.WriteString("    ")
	for col := 0; col < m.game.BoardSize; col++ {
		sb.WriteString(fmt.Sprintf(" %c  ", 'A'+col))
	}
	sb.WriteString("\n")

	// Board
	for row := 0; row < m.game.BoardSize; row++ {
		sb.WriteString(fmt.Sprintf("%2d  ", row+1))

		for col := 0; col < m.game.BoardSize; col++ {
			pos := game.Position{Row: row, Col: col}
			cell := m.game.PlayerBoard.GetCell(pos)
			cellStr := renderCell(cell, false, false, true)
			sb.WriteString(cellStr)
		}
		sb.WriteString("\n")
	}

	return boardStyle.Render(sb.String())
}

func renderEnemyBoard(m Model) string {
	var sb strings.Builder

	sb.WriteString(headerStyle.Render("Enemy Waters"))
	sb.WriteString("\n\n")

	// Column headers
	sb.WriteString("    ")
	for col := 0; col < m.game.BoardSize; col++ {
		sb.WriteString(fmt.Sprintf(" %c  ", 'A'+col))
	}
	sb.WriteString("\n")

	// Board
	for row := 0; row < m.game.BoardSize; row++ {
		sb.WriteString(fmt.Sprintf("%2d  ", row+1))

		for col := 0; col < m.game.BoardSize; col++ {
			pos := game.Position{Row: row, Col: col}
			cell := m.game.ComputerBoard.GetCell(pos)

			isCursor := row == m.cursorRow && col == m.cursorCol
			cellStr := renderCell(cell, isCursor, false, false)
			sb.WriteString(cellStr)
		}
		sb.WriteString("\n")
	}

	return boardStyle.Render(sb.String())
}

func renderCell(cell game.CellState, isCursor bool, isPreview bool, showShips bool) string {
	symbol := "~"

	switch cell {
	case game.Empty:
		symbol = "~"
		if isCursor {
			return cursorStyle.Render("[" + symbol + "]")
		}
		if isPreview {
			return shipStyle.Render("▓")
		}
		return waterStyle.Render(" " + symbol + " ")

	case game.ShipCell:
		if showShips {
			symbol = "█"
			return shipStyle.Render(" " + symbol + " ")
		}
		// For hidden enemy ships, show cursor if applicable
		symbol = "~"
		if isCursor {
			return cursorStyle.Render("[" + symbol + "]")
		}
		return waterStyle.Render(" " + symbol + " ")

	case game.Hit:
		symbol = "X"
		if isCursor {
			return grayCursorStyle.Render("[" + symbol + "]")
		}
		return hitStyle.Render(" " + symbol + " ")

	case game.Miss:
		symbol = "○"
		if isCursor {
			return grayCursorStyle.Render("[" + symbol + "]")
		}
		return missStyle.Render(" " + symbol + " ")
	}

	return waterStyle.Render(" ~ ")
}

func renderGameOver(m Model) string {
	var sb strings.Builder

	// Show both boards
	sb.WriteString(renderBattleBoards(m))
	sb.WriteString("\n\n")

	// Game over message
	gameOverMsg := ""
	if m.game.Winner == "Player" {
		gameOverMsg = "🎉 VICTORY! You sunk all enemy ships! 🎉"
	} else {
		gameOverMsg = "💥 DEFEAT! All your ships were sunk! 💥"
	}

	gameOverStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFCC00")).
		Padding(1, 2).
		Border(lipgloss.DoubleBorder())

	sb.WriteString(gameOverStyle.Render(gameOverMsg))

	return sb.String()
}

func renderHelp(m Model) string {
	var sb strings.Builder

	sb.WriteString("Controls:\n")

	switch m.game.Phase {
	case game.PlacementPhase:
		sb.WriteString("  Arrow Keys/WASD - Move cursor\n")
		sb.WriteString("  O - Toggle orientation (Horizontal/Vertical)\n")
		sb.WriteString("  Space/Enter - Place ship\n")
	case game.PlayerTurnPhase, game.ComputerTurnPhase:
		sb.WriteString("  Arrow Keys/WASD - Move cursor\n")
		sb.WriteString("  Space/Enter - Fire!\n")
	case game.GameOverPhase:
		sb.WriteString("  R - Restart game\n")
	}

	sb.WriteString("  H - Toggle help\n")
	sb.WriteString("  Q - Quit\n")

	return helpStyle.Render(sb.String())
}
