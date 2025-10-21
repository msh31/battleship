package game

import (
	"math/rand"
	"time"
)

// GamePhase represents the current phase of the game
type GamePhase int

const (
	MainMenuPhase GamePhase = iota
	PlacementPhase
	PlayerTurnPhase
	ComputerTurnPhase
	GameOverPhase
)

// Game represents the game state
type Game struct {
	PlayerBoard      *Board
	ComputerBoard    *Board
	Phase            GamePhase
	BoardSize        int
	CurrentShip      int // For placement phase
	ShipTypes        []ShipType
	Winner           string
	LastMessage      string
	ClaudeThinking   string
	Random           *rand.Rand
}

// Claude thinking messages
var thinkingMessages = []string{
	"Pondering",
	"Hatching a plan",
	"Simmering",
	"Meandering",
	"Contemplating",
	"Channelling",
	"Strategizing",
	"Calculating",
	"Analyzing",
	"Reasoning",
	"Deliberating",
	"Ruminating",
}

// NewGame creates a new game
func NewGame(boardSize int) *Game {
	g := &Game{
		PlayerBoard:   NewBoard(boardSize),
		ComputerBoard: NewBoard(boardSize),
		Phase:         PlacementPhase,
		BoardSize:     boardSize,
		CurrentShip:   0,
		ShipTypes: []ShipType{
			Carrier,
			Battleship,
			Cruiser,
			Submarine,
			Destroyer,
		},
		Random: rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	// Place computer ships randomly
	g.placeComputerShips()

	return g
}

// placeComputerShips randomly places all ships for the computer
func (g *Game) placeComputerShips() {
	for _, shipType := range g.ShipTypes {
		ship := NewShip(shipType)
		placed := false

		for !placed {
			row := g.Random.Intn(g.BoardSize)
			col := g.Random.Intn(g.BoardSize)
			orientation := Orientation(g.Random.Intn(2))

			pos := Position{Row: row, Col: col}
			placed = g.ComputerBoard.PlaceShip(ship, pos, orientation)
		}
	}
}

// PlacePlayerShip places the current ship for the player
func (g *Game) PlacePlayerShip(pos Position, orientation Orientation) bool {
	if g.Phase != PlacementPhase || g.CurrentShip >= len(g.ShipTypes) {
		return false
	}

	ship := NewShip(g.ShipTypes[g.CurrentShip])
	if g.PlayerBoard.PlaceShip(ship, pos, orientation) {
		g.CurrentShip++
		if g.CurrentShip >= len(g.ShipTypes) {
			g.Phase = PlayerTurnPhase
			g.LastMessage = "All ships placed! Your turn to attack!"
		}
		return true
	}

	return false
}

// GetCurrentShipForPlacement returns the ship currently being placed
func (g *Game) GetCurrentShipForPlacement() *Ship {
	if g.CurrentShip >= len(g.ShipTypes) {
		return nil
	}
	return NewShip(g.ShipTypes[g.CurrentShip])
}

// PlayerAttack performs a player attack on the computer's board
func (g *Game) PlayerAttack(pos Position) bool {
	if g.Phase != PlayerTurnPhase {
		return false
	}

	hit, ship := g.ComputerBoard.Attack(pos)

	if hit {
		if ship != nil && ship.IsSunk() {
			g.LastMessage = "Hit! You sunk the enemy's " + ship.Name + "!"
		} else {
			g.LastMessage = "Hit!"
		}

		if g.ComputerBoard.AllShipsSunk() {
			g.Phase = GameOverPhase
			g.Winner = "Player"
			g.LastMessage = "Victory! You sunk all enemy ships!"
			return true
		}
	} else {
		g.LastMessage = "Miss!"
	}

	g.Phase = ComputerTurnPhase
	g.ClaudeThinking = g.GetRandomThinkingMessage()
	return true
}

// GetRandomThinkingMessage returns a random thinking message for Claude
func (g *Game) GetRandomThinkingMessage() string {
	return thinkingMessages[g.Random.Intn(len(thinkingMessages))]
}

// ComputerAttack performs a computer attack on the player's board
func (g *Game) ComputerAttack() {
	if g.Phase != ComputerTurnPhase {
		return
	}

	// Simple AI: random attacks on unattacked cells
	var pos Position
	found := false

	for !found {
		row := g.Random.Intn(g.BoardSize)
		col := g.Random.Intn(g.BoardSize)
		pos = Position{Row: row, Col: col}

		cell := g.PlayerBoard.GetCell(pos)
		if cell != Hit && cell != Miss {
			found = true
		}
	}

	hit, ship := g.PlayerBoard.Attack(pos)

	if hit {
		if ship != nil && ship.IsSunk() {
			g.LastMessage = "Computer sunk your " + ship.Name + "!"
		} else {
			g.LastMessage = "Computer hit your ship!"
		}

		if g.PlayerBoard.AllShipsSunk() {
			g.Phase = GameOverPhase
			g.Winner = "Computer"
			g.LastMessage = "Defeat! All your ships were sunk!"
			return
		}
	} else {
		g.LastMessage = "Computer missed!"
	}

	g.Phase = PlayerTurnPhase
}
