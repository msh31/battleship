package game

// CellState represents the state of a cell on the board
type CellState int

const (
	Empty CellState = iota
	ShipCell
	Miss
	Hit
)

// Board represents a game board
type Board struct {
	Size  int
	Grid  [][]CellState
	Ships []*Ship
}

// NewBoard creates a new board of the given size
func NewBoard(size int) *Board {
	grid := make([][]CellState, size)
	for i := range grid {
		grid[i] = make([]CellState, size)
	}

	return &Board{
		Size:  size,
		Grid:  grid,
		Ships: make([]*Ship, 0),
	}
}

// IsValidPosition checks if a position is within board bounds
func (b *Board) IsValidPosition(pos Position) bool {
	return pos.Row >= 0 && pos.Row < b.Size && pos.Col >= 0 && pos.Col < b.Size
}

// CanPlaceShip checks if a ship can be placed at the given position
func (b *Board) CanPlaceShip(pos Position, length int, orientation Orientation) bool {
	positions := b.getShipPositions(pos, length, orientation)

	for _, p := range positions {
		if !b.IsValidPosition(p) {
			return false
		}
		if b.Grid[p.Row][p.Col] == ShipCell {
			return false
		}
	}

	return true
}

// PlaceShip places a ship on the board
func (b *Board) PlaceShip(ship *Ship, pos Position, orientation Orientation) bool {
	if !b.CanPlaceShip(pos, ship.Length, orientation) {
		return false
	}

	positions := b.getShipPositions(pos, ship.Length, orientation)
	ship.Positions = positions

	for _, p := range positions {
		b.Grid[p.Row][p.Col] = ShipCell
	}

	b.Ships = append(b.Ships, ship)
	return true
}

// getShipPositions returns all positions a ship would occupy
func (b *Board) getShipPositions(pos Position, length int, orientation Orientation) []Position {
	positions := make([]Position, length)

	for i := 0; i < length; i++ {
		if orientation == Horizontal {
			positions[i] = Position{Row: pos.Row, Col: pos.Col + i}
		} else {
			positions[i] = Position{Row: pos.Row + i, Col: pos.Col}
		}
	}

	return positions
}

// Attack performs an attack at the given position
func (b *Board) Attack(pos Position) (bool, *Ship) {
	if !b.IsValidPosition(pos) {
		return false, nil
	}

	cell := b.Grid[pos.Row][pos.Col]

	if cell == Miss || cell == Hit {
		return false, nil // Already attacked
	}

	if cell == ShipCell {
		b.Grid[pos.Row][pos.Col] = Hit
		// Find which ship was hit
		for _, ship := range b.Ships {
			if ship.Hit(pos) {
				return true, ship
			}
		}
		return true, nil
	}

	b.Grid[pos.Row][pos.Col] = Miss
	return false, nil
}

// AllShipsSunk returns true if all ships on the board are sunk
func (b *Board) AllShipsSunk() bool {
	for _, ship := range b.Ships {
		if !ship.IsSunk() {
			return false
		}
	}
	return len(b.Ships) > 0
}

// GetCell returns the state of a cell (for opponent tracking)
func (b *Board) GetCell(pos Position) CellState {
	if !b.IsValidPosition(pos) {
		return Empty
	}
	return b.Grid[pos.Row][pos.Col]
}
