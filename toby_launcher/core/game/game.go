package game

import (
	"os/exec"
)

type RawGameData struct {
	Description string   `json:"description"`
	Iwad        string   `json:"iwad"`
	Config      string   `json:"config"`
	Files       []string `json:"files"`
	Params      []string `json:"params"`
}

type RawGamesData map[string]GameData

type GameData struct {
	Name        string
	Description string
	Config      string
	Iwad        string
	Files       []string
	Params      []string
}

type Game struct {
	Info      *GameData
	cmd       *exec.Cmd
	IsRunning bool
}
