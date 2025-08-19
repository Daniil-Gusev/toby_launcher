package game

import (
	"os/exec"
	"toby_launcher/apperrors"
)

type RawGameData struct {
	Description string   `json:"description"`
	Iwads       []string `json:"iwads"`
	Config      string   `json:"config"`
	Files       []string `json:"files"`
	Params      []string `json:"params"`
}

func (d RawGameData) validate() error {
	if d.Iwads == nil {
		return apperrors.New(apperrors.Err, "field \"iwads\" is missing", nil)
	}
	if len(d.Iwads) == 0 {
		return apperrors.New(apperrors.Err, "field \"iwads\" is empty", nil)
	}
	return nil
}

type RawGamesData map[string]RawGameData

type GameData struct {
	Name        string
	Description string
	Config      string
	Iwads       []string
	Files       []string
	Params      []string
}

type Game struct {
	Info      *GameData
	cmd       *exec.Cmd
	IsRunning bool
}
