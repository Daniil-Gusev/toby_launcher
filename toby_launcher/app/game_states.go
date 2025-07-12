package app

import (
	"fmt"
	"time"
	"toby_launcher/apperrors"
	"toby_launcher/core"
	"toby_launcher/core/game"
	"toby_launcher/core/validation"
)

type InitGameState struct {
	core.BaseState
	game *game.GameData
}

func (s *InitGameState) Name() string {
	return "init game state"
}

func (s *InitGameState) Init(ctx *core.AppContext, ui *core.UiContext) (core.State, error) {
	if s.game == nil {
		ui.DisplayError(apperrors.New(apperrors.Err, "Error: The game was not specified during initialization.", nil))
		return ctx.GetPreviousState()
	}
	return s, nil
}

func (s *InitGameState) Display(ctx *core.AppContext, ui *core.UiContext) {
	ui.DisplayText("Loading GZDoom...\r\n")
}

func (s *InitGameState) Handle(ctx *core.AppContext, ui *core.UiContext, input string) (core.State, error) {
	if err := ctx.GameManager.StartGame(s.game); err != nil {
		ui.DisplayError(err)
		return ctx.GetPreviousState()
	}
	msg := fmt.Sprintf("Game starting: %s. Good luck!\r\n", s.game.Name)
	ui.DisplayText(msg)
	ui.TtsManager.Speak(msg)
	return &GameState{}, nil
}

func (s *InitGameState) RequiresInput() bool {
	return false
}

type GameState struct{ core.BaseState }

func (s *GameState) Name() string {
	return "game"
}

func (s *GameState) Description() string {
	return "You are in the game."
}

func (s *GameState) Handle(ctx *core.AppContext, ui *core.UiContext, input string) (core.State, error) {
	time.Sleep(2000 * time.Millisecond)
	if !ctx.GameManager.GameIsRunning() {
		return ctx.GetPreviousState()
	}
	return s, nil
}

func (s *GameState) RequiresInput() bool {
	return false
}

type GameSelectionMenuState struct {
	core.BaseState
}

func (m *GameSelectionMenuState) Name() string {
	return "game selection menu"
}

func (m *GameSelectionMenuState) Description() string {
	return "You are in the game selection menu. You need to enter the number of the game you want to launch."
}

func (m *GameSelectionMenuState) Display(ctx *core.AppContext, ui *core.UiContext) {
	ui.DisplayText("0. Back.\r\n\r\n")
	ui.DisplayText("The following games are available to you:\r\n\r\n")
	games := ctx.GameManager.AvailableGames()
	for i, game := range games {
		ui.DisplayText(fmt.Sprintf("%d. %s.\r\n%s\r\n\r\n", i+1, game.Name, game.Description))
	}
	if len(games) == 0 {
		ui.DisplayText("No games are currently available.\r\n\r\n")
	}
	ui.DisplayText("Make your choice.\r\n")
}

func (m *GameSelectionMenuState) Handle(ctx *core.AppContext, ui *core.UiContext, input string) (core.State, error) {
	option, err := validation.ParseInt(input)
	if err != nil {
		return m, err
	}
	games := ctx.GameManager.AvailableGames()
	maxOption := len(games)
	if option < 0 || option > maxOption {
		ui.DisplayText("There is no such item in the menu.\r\n")
		return m, nil
	}
	if option == 0 {
		return ctx.GetPreviousState()
	}
	game := games[option-1]
	ui.DisplayText(fmt.Sprintf("You have chosen a game: %s.\r\n", game.Name))
	return &InitGameState{game: game}, nil
}
