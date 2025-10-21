package game

// ShipType represents different types of ships
type ShipType int

const (
	Carrier ShipType = iota
	Battleship
	Cruiser
	Submarine
	Destroyer
)

// Ship represents a ship on the board
type Ship struct {
	Type      ShipType
	Length    int
	Positions []Position
	Hits      []bool
	Name      string
}

// Position represents a coordinate on the board
type Position struct {
	Row int
	Col int
}

// Orientation represents ship placement direction
type Orientation int

const (
	Horizontal Orientation = iota
	Vertical
)

// NewShip creates a new ship of the given type
func NewShip(shipType ShipType) *Ship {
	ship := &Ship{
		Type: shipType,
	}

	switch shipType {
	case Carrier:
		ship.Length = 5
		ship.Name = "Carrier"
	case Battleship:
		ship.Length = 4
		ship.Name = "Battleship"
	case Cruiser:
		ship.Length = 3
		ship.Name = "Cruiser"
	case Submarine:
		ship.Length = 3
		ship.Name = "Submarine"
	case Destroyer:
		ship.Length = 2
		ship.Name = "Destroyer"
	}

	ship.Hits = make([]bool, ship.Length)
	return ship
}

// IsSunk returns true if all positions of the ship have been hit
func (s *Ship) IsSunk() bool {
	for _, hit := range s.Hits {
		if !hit {
			return false
		}
	}
	return true
}

// Hit marks a position on the ship as hit and returns true if successful
func (s *Ship) Hit(pos Position) bool {
	for i, p := range s.Positions {
		if p.Row == pos.Row && p.Col == pos.Col {
			s.Hits[i] = true
			return true
		}
	}
	return false
}
