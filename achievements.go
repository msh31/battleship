package main

import (
	"battleship/game"
	"encoding/json"
	"os"
	"path/filepath"
)

type Achievement struct {
	ID          string
	Name        string
	Description string
	Unlocked    bool
}

type Achievements struct {
	PerfectGame     bool `json:"perfect_game"`      // Win without any of your ships sunk
	Sharpshooter    bool `json:"sharpshooter"`      // Win with 90%+ accuracy
	ComebackKing    bool `json:"comeback_king"`     // Win after losing 4 ships
	FirstBlood      bool `json:"first_blood"`       // Win your first game
	HardcoreVictor  bool `json:"hardcore_victor"`   // Beat Hard difficulty
	SalvoMaster     bool `json:"salvo_master"`      // Win in Salvo mode
	Efficient       bool `json:"efficient"`         // Win in under 50 shots
	LuckyShot       bool `json:"lucky_shot"`        // Sink a ship on first hit
	Domination      bool `json:"domination"`        // Win with all ships intact
	SmallBoardWin   bool `json:"small_board_win"`   // Win on 8x8 board
	LargeBoardWin   bool `json:"large_board_win"`   // Win on 12x12 board
}

func LoadAchievements() *Achievements {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return &Achievements{}
	}

	filePath := filepath.Join(homeDir, ".battleship_achievements.json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return &Achievements{}
	}

	var achievements Achievements
	if err := json.Unmarshal(data, &achievements); err != nil {
		return &Achievements{}
	}

	return &achievements
}

func (a *Achievements) Save() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	filePath := filepath.Join(homeDir, ".battleship_achievements.json")
	data, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

func (a *Achievements) GetAll() []Achievement {
	return []Achievement{
		{
			ID:          "first_blood",
			Name:        "First Blood",
			Description: "Win your first game",
			Unlocked:    a.FirstBlood,
		},
		{
			ID:          "perfect_game",
			Name:        "Perfect Game",
			Description: "Win without losing any ships",
			Unlocked:    a.PerfectGame,
		},
		{
			ID:          "domination",
			Name:        "Domination",
			Description: "Win with all ships at full health",
			Unlocked:    a.Domination,
		},
		{
			ID:          "sharpshooter",
			Name:        "Sharpshooter",
			Description: "Win with 90%+ accuracy",
			Unlocked:    a.Sharpshooter,
		},
		{
			ID:          "comeback_king",
			Name:        "Comeback King",
			Description: "Win after losing 4 ships",
			Unlocked:    a.ComebackKing,
		},
		{
			ID:          "hardcore_victor",
			Name:        "Hardcore Victor",
			Description: "Beat Hard difficulty",
			Unlocked:    a.HardcoreVictor,
		},
		{
			ID:          "salvo_master",
			Name:        "Salvo Master",
			Description: "Win in Salvo mode",
			Unlocked:    a.SalvoMaster,
		},
		{
			ID:          "efficient",
			Name:        "Efficient",
			Description: "Win in under 50 shots",
			Unlocked:    a.Efficient,
		},
		{
			ID:          "small_board_win",
			Name:        "Compact Commander",
			Description: "Win on 8x8 board",
			Unlocked:    a.SmallBoardWin,
		},
		{
			ID:          "large_board_win",
			Name:        "Admiral of the Seas",
			Description: "Win on 12x12 board",
			Unlocked:    a.LargeBoardWin,
		},
	}
}

func (a *Achievements) CheckAndUnlock(g *game.Game) []Achievement {
	newlyUnlocked := []Achievement{}

	if g.Winner != "Player" {
		return newlyUnlocked
	}

	// First Blood - first win
	if !a.FirstBlood {
		a.FirstBlood = true
		newlyUnlocked = append(newlyUnlocked, Achievement{
			ID:          "first_blood",
			Name:        "First Blood",
			Description: "Win your first game",
			Unlocked:    true,
		})
	}

	// Perfect Game - no ships sunk
	sunkShips := 0
	for _, ship := range g.PlayerBoard.Ships {
		if ship.IsSunk() {
			sunkShips++
		}
	}

	if sunkShips == 0 && !a.PerfectGame {
		a.PerfectGame = true
		newlyUnlocked = append(newlyUnlocked, Achievement{
			ID:          "perfect_game",
			Name:        "Perfect Game",
			Description: "Win without losing any ships",
			Unlocked:    true,
		})
	}

	// Domination - all ships at full health
	anyDamage := false
	for _, ship := range g.PlayerBoard.Ships {
		for _, hit := range ship.Hits {
			if hit {
				anyDamage = true
				break
			}
		}
		if anyDamage {
			break
		}
	}

	if !anyDamage && !a.Domination {
		a.Domination = true
		newlyUnlocked = append(newlyUnlocked, Achievement{
			ID:          "domination",
			Name:        "Domination",
			Description: "Win with all ships at full health",
			Unlocked:    true,
		})
	}

	// Comeback King - win after losing 4 ships
	if sunkShips >= 4 && !a.ComebackKing {
		a.ComebackKing = true
		newlyUnlocked = append(newlyUnlocked, Achievement{
			ID:          "comeback_king",
			Name:        "Comeback King",
			Description: "Win after losing 4 ships",
			Unlocked:    true,
		})
	}

	// Hardcore Victor - beat Hard difficulty
	if g.Difficulty == game.Hard && !a.HardcoreVictor {
		a.HardcoreVictor = true
		newlyUnlocked = append(newlyUnlocked, Achievement{
			ID:          "hardcore_victor",
			Name:        "Hardcore Victor",
			Description: "Beat Hard difficulty",
			Unlocked:    true,
		})
	}

	// Salvo Master - win in Salvo mode
	if g.SalvoMode && !a.SalvoMaster {
		a.SalvoMaster = true
		newlyUnlocked = append(newlyUnlocked, Achievement{
			ID:          "salvo_master",
			Name:        "Salvo Master",
			Description: "Win in Salvo mode",
			Unlocked:    true,
		})
	}

	// Board size achievements
	if g.BoardSize == 8 && !a.SmallBoardWin {
		a.SmallBoardWin = true
		newlyUnlocked = append(newlyUnlocked, Achievement{
			ID:          "small_board_win",
			Name:        "Compact Commander",
			Description: "Win on 8x8 board",
			Unlocked:    true,
		})
	}

	if g.BoardSize == 12 && !a.LargeBoardWin {
		a.LargeBoardWin = true
		newlyUnlocked = append(newlyUnlocked, Achievement{
			ID:          "large_board_win",
			Name:        "Admiral of the Seas",
			Description: "Win on 12x12 board",
			Unlocked:    true,
		})
	}

	// Save if any new achievements
	if len(newlyUnlocked) > 0 {
		a.Save()
	}

	return newlyUnlocked
}
