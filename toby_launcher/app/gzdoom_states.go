package app

import (
	"fmt"
	"strings"
	"toby_launcher/core"
)

type GzdoomSettingsMenuState struct{ core.BaseState }

func (m *GzdoomSettingsMenuState) Name() string {
	return "gzdoom settings menu"
}

func NewGzdoomSettingsMenu(ctx *core.AppContext, ui *core.UiContext) *core.MenuState {
	parrentState := &GzdoomSettingsMenuState{}
	options := []*core.MenuOption{
		{Id: 0,
			Description: "Back.",
			NextState:   func() (core.State, error) { return ctx.GetPreviousState() },
		},
		{Id: 1,
			Description: "Change GZDoom launch parametrs.",
			NextState:   func() (core.State, error) { return &ChangeLaunchParamsState{}, nil },
		},
		core.NewSwitchMenuOption(2, &ctx.Config.Gzdoom.DebugOutput, "debug output"),
	}
	return core.NewMenu(parrentState, options, "")
}

type ChangeLaunchParamsState struct{ core.BaseState }

func (s *ChangeLaunchParamsState) Name() string {
	return "change launch params"
}

func (s *ChangeLaunchParamsState) Description() string {
	return "You need to specify the parameters that will be passed to GZDoom when starting any game. The separator between the parameters is a semicolon. To reset the parameters, press \"enter\"."
}

func (s *ChangeLaunchParamsState) Display(ctx *core.AppContext, ui *core.UiContext) {
	ui.DisplayText("Enter the desired GZDoom launch parameters, separating them with semicolons.\r\n")
	if len(ctx.Config.Gzdoom.LaunchParams) > 0 {
		ui.DisplayText(fmt.Sprintf("Current value: %s\r\n", strings.Join(ctx.Config.Gzdoom.LaunchParams, "; ")))
	}
}

func (s *ChangeLaunchParamsState) Handle(ctx *core.AppContext, ui *core.UiContext, input string) (core.State, error) {
	if input == "" {
		ctx.Config.Gzdoom.LaunchParams = []string{}
		ui.DisplayText("Launch parameters have bin reset.\r\n")
		return ctx.GetPreviousState()
	}
	rawParams := strings.Split(input, ";")
	params := make([]string, 0, len(rawParams))
	for _, p := range rawParams {
		param := strings.TrimSpace(p)
		if param != "" {
			params = append(params, param)
		}
	}
	if len(params) == 0 {
		ui.DisplayText("Launch parameters remain unchanged.\r\n")
		return ctx.GetPreviousState()
	}
	ctx.Config.Gzdoom.LaunchParams = params
	ui.DisplayText(fmt.Sprintf("The following launch parameters are set: %s.\r\n", strings.Join(params, "; ")))
	return ctx.GetPreviousState()
}

func (s *ChangeLaunchParamsState) Command() []core.Command {
	return []core.Command{&core.BackCommand{}}
}
