package main

import (
	"battleship/game"
	"encoding/json"
	"os"
	"path/filepath"
)

// Stats tracks player performance
type Stats struct {
	EasyWins    int `json:"easy_wins"`
	EasyLosses  int `json:"easy_losses"`
	NormalWins  int `json:"normal_wins"`
	NormalLosses int `json:"normal_losses"`
	HardWins    int `json:"hard_wins"`
	HardLosses  int `json:"hard_losses"`
}

// getStatsPath returns the path to the stats file
func getStatsPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".battleship_stats.json"), nil
}

// LoadStats loads stats from disk
func LoadStats() (*Stats, error) {
	path, err := getStatsPath()
	if err != nil {
		return &Stats{}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Stats{}, nil
		}
		return nil, err
	}

	var stats Stats
	if err := json.Unmarshal(data, &stats); err != nil {
		return &Stats{}, nil
	}

	return &stats, nil
}

// SaveStats saves stats to disk
func (s *Stats) Save() error {
	path, err := getStatsPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// RecordWin records a win for the given difficulty
func (s *Stats) RecordWin(difficulty game.Difficulty) {
	switch difficulty {
	case game.Easy:
		s.EasyWins++
	case game.Normal:
		s.NormalWins++
	case game.Hard:
		s.HardWins++
	}
	s.Save()
}

// RecordLoss records a loss for the given difficulty
func (s *Stats) RecordLoss(difficulty game.Difficulty) {
	switch difficulty {
	case game.Easy:
		s.EasyLosses++
	case game.Normal:
		s.NormalLosses++
	case game.Hard:
		s.HardLosses++
	}
	s.Save()
}

// GetWins returns wins for the given difficulty
func (s *Stats) GetWins(difficulty game.Difficulty) int {
	switch difficulty {
	case game.Easy:
		return s.EasyWins
	case game.Normal:
		return s.NormalWins
	case game.Hard:
		return s.HardWins
	}
	return 0
}

// GetLosses returns losses for the given difficulty
func (s *Stats) GetLosses(difficulty game.Difficulty) int {
	switch difficulty {
	case game.Easy:
		return s.EasyLosses
	case game.Normal:
		return s.NormalLosses
	case game.Hard:
		return s.HardLosses
	}
	return 0
}

// GetTotalGames returns total games for the given difficulty
func (s *Stats) GetTotalGames(difficulty game.Difficulty) int {
	return s.GetWins(difficulty) + s.GetLosses(difficulty)
}

// GetWinRate returns win rate as a percentage for the given difficulty
func (s *Stats) GetWinRate(difficulty game.Difficulty) float64 {
	total := s.GetTotalGames(difficulty)
	if total == 0 {
		return 0
	}
	return float64(s.GetWins(difficulty)) / float64(total) * 100
}
