package app

import (
	"fmt"
	"toby_launcher/core"
	"toby_launcher/core/version"
)

type StartState struct{ core.BaseState }

func (s *StartState) Handle(ctx *core.AppContext, ui *core.UiContext, _ string) (core.State, error) {
	ui.TtsManager.Speak(fmt.Sprintf("Welcome to the %s!", version.AppName))
	return NewMainMenu(ctx, ui), nil
}

func (s *StartState) RequiresInput() bool {
	return false
}

type MainMenuState struct{ core.BaseState }

func (m *MainMenuState) Name() string {
	return "main menu"
}

func NewMainMenu(ctx *core.AppContext, ui *core.UiContext) *core.MenuState {
	parentState := &MainMenuState{}
	options := []core.MenuOption{
		{Id: 0,
			Description: "Exit.",
			NextState:   func() (core.State, error) { return &core.ExitState{}, nil },
		},
		{Id: 1,
			Description: "Play.",
			NextState:   func() (core.State, error) { return &GameSelectionMenuState{}, nil },
		},
		{Id: 2,
			Description: "Settings.",
			NextState:   func() (core.State, error) { return NewSettingsMenu(ctx, ui), nil },
		},
	}
	return core.NewMenu(parentState, options, "")
}

type SettingsMenuState struct{ core.BaseState }

func (m *SettingsMenuState) Name() string {
	return "settings"
}

func NewSettingsMenu(ctx *core.AppContext, ui *core.UiContext) *core.MenuState {
	parrentState := &SettingsMenuState{}
	options := []core.MenuOption{
		{Id: 0,
			Description: "Back.",
			NextState:   func() (core.State, error) { return ctx.GetPreviousState() },
		},
		{Id: 1,
			Description: "Speech settings.",
			NextState:   func() (core.State, error) { return NewSpeechSettingsMenu(ctx, ui), nil },
		},
		{Id: 2,
			Description: "GZDoom settings.",
			NextState:   func() (core.State, error) { return NewGzdoomSettingsMenu(ctx, ui), nil },
		},
	}
	return core.NewMenu(parrentState, options, "")
}
