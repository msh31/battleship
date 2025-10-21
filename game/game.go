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

// Difficulty represents the AI difficulty level
type Difficulty int

const (
	Easy Difficulty = iota
	Normal
	Hard
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
	Difficulty       Difficulty
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
			g.LastMessage = "Hit! You sunk Captain Claude's " + ship.Name + "!"
		} else {
			g.LastMessage = "Hit!"
		}

		if g.ComputerBoard.AllShipsSunk() {
			g.Phase = GameOverPhase
			g.Winner = "Player"
			g.LastMessage = "Victory! You sunk Captain Claude's fleet!"
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

	var pos Position

	// Choose attack strategy based on difficulty
	switch g.Difficulty {
	case Easy:
		pos = g.easyAIAttack()
	case Normal:
		pos = g.normalAIAttack()
	case Hard:
		pos = g.hardAIAttack()
	}

	hit, ship := g.PlayerBoard.Attack(pos)

	if hit {
		if ship != nil && ship.IsSunk() {
			g.LastMessage = "Claude sunk your " + ship.Name + "!"
		} else {
			g.LastMessage = "Claude hit your ship!"
		}

		if g.PlayerBoard.AllShipsSunk() {
			g.Phase = GameOverPhase
			g.Winner = "Claude"
			g.LastMessage = "Defeat! All your ships were sunk!"
			return
		}
	} else {
		g.LastMessage = "Claude missed!"
	}

	g.Phase = PlayerTurnPhase
}

// easyAIAttack implements easy difficulty - random attacks
func (g *Game) easyAIAttack() Position {
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

	return pos
}

// normalAIAttack implements normal difficulty - hunts around hits
func (g *Game) normalAIAttack() Position {
	// First, look for existing hits to follow up on
	for row := 0; row < g.BoardSize; row++ {
		for col := 0; col < g.BoardSize; col++ {
			if g.PlayerBoard.Grid[row][col] == Hit {
				// Found a hit, try adjacent cells
				adjacents := []Position{
					{Row: row - 1, Col: col},
					{Row: row + 1, Col: col},
					{Row: row, Col: col - 1},
					{Row: row, Col: col + 1},
				}

				// Shuffle adjacents for variety
				for i := range adjacents {
					j := g.Random.Intn(i + 1)
					adjacents[i], adjacents[j] = adjacents[j], adjacents[i]
				}

				for _, adj := range adjacents {
					if g.PlayerBoard.IsValidPosition(adj) {
						cell := g.PlayerBoard.GetCell(adj)
						if cell != Hit && cell != Miss {
							return adj
						}
					}
				}
			}
		}
	}

	// No hits to follow up on, attack randomly
	return g.easyAIAttack()
}

// hardAIAttack implements hard difficulty - smart pattern hunting and direction following
func (g *Game) hardAIAttack() Position {
	// Look for hits in a line (ship orientation detected)
	for row := 0; row < g.BoardSize; row++ {
		for col := 0; col < g.BoardSize; col++ {
			if g.PlayerBoard.Grid[row][col] == Hit {
				// Check horizontal line
				if col+1 < g.BoardSize && g.PlayerBoard.Grid[row][col+1] == Hit {
					// Found horizontal ship, extend in both directions
					// Try right first
					if col+2 < g.BoardSize {
						adj := Position{Row: row, Col: col + 2}
						cell := g.PlayerBoard.GetCell(adj)
						if cell != Hit && cell != Miss {
							return adj
						}
					}
					// Try left
					if col-1 >= 0 {
						adj := Position{Row: row, Col: col - 1}
						cell := g.PlayerBoard.GetCell(adj)
						if cell != Hit && cell != Miss {
							return adj
						}
					}
				}

				// Check vertical line
				if row+1 < g.BoardSize && g.PlayerBoard.Grid[row+1][col] == Hit {
					// Found vertical ship, extend in both directions
					// Try down first
					if row+2 < g.BoardSize {
						adj := Position{Row: row + 2, Col: col}
						cell := g.PlayerBoard.GetCell(adj)
						if cell != Hit && cell != Miss {
							return adj
						}
					}
					// Try up
					if row-1 >= 0 {
						adj := Position{Row: row - 1, Col: col}
						cell := g.PlayerBoard.GetCell(adj)
						if cell != Hit && cell != Miss {
							return adj
						}
					}
				}
			}
		}
	}

	// No line detected, use normal mode's adjacent hunting
	for row := 0; row < g.BoardSize; row++ {
		for col := 0; col < g.BoardSize; col++ {
			if g.PlayerBoard.Grid[row][col] == Hit {
				adjacents := []Position{
					{Row: row - 1, Col: col},
					{Row: row + 1, Col: col},
					{Row: row, Col: col - 1},
					{Row: row, Col: col + 1},
				}

				for _, adj := range adjacents {
					if g.PlayerBoard.IsValidPosition(adj) {
						cell := g.PlayerBoard.GetCell(adj)
						if cell != Hit && cell != Miss {
							return adj
						}
					}
				}
			}
		}
	}

	// No hits to follow, use checkerboard pattern for efficient hunting
	for row := 0; row < g.BoardSize; row++ {
		for col := 0; col < g.BoardSize; col++ {
			if (row+col)%2 == 0 { // Checkerboard pattern
				pos := Position{Row: row, Col: col}
				cell := g.PlayerBoard.GetCell(pos)
				if cell != Hit && cell != Miss {
					return pos
				}
			}
		}
	}

	// Checkerboard exhausted, fill in remaining cells
	return g.easyAIAttack()
}
