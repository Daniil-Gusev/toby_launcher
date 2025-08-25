package app

import (
	"fmt"
	"strings"
	"toby_launcher/core"
	"toby_launcher/core/game"
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
			NextState:   ctx.GetPreviousState,
		},
		{Id: 1,
			Description: "Change additional GZDoom launch parametrs.",
			NextState:   func() (core.State, error) { return &ChangeLaunchParamsState{}, nil },
		},
		{Id: 2,
			Description: "Change video backend ($backend).",
			Params:      func() map[string]any { return map[string]any{"backend": ctx.GameManager.Params.VideoBackend()} },
			NextState:   func() (core.State, error) { return NewVideoBackendMenu(ctx, ui), nil },
		},
		core.NewSwitchMenuOption(3, "music", ctx.GameManager.Params.MusicPtr()),
		core.NewSwitchMenuOption(4, "sound affects", ctx.GameManager.Params.SoundfxPtr()),
		core.NewSwitchMenuOption(5, "debug output", &ctx.Config.Gzdoom.DebugOutput),
		core.NewSwitchMenuOption(6, "logging", &ctx.Config.Gzdoom.Logging),
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
	if len(ctx.Config.Gzdoom.AdditionalLaunchParams) > 0 {
		ui.DisplayText(fmt.Sprintf("Current value: %s\r\n", strings.Join(ctx.Config.Gzdoom.AdditionalLaunchParams, "; ")))
	}
}

func (s *ChangeLaunchParamsState) Handle(ctx *core.AppContext, ui *core.UiContext, input string) (core.State, error) {
	if input == "" {
		ctx.Config.Gzdoom.AdditionalLaunchParams = []string{}
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
	ctx.Config.Gzdoom.AdditionalLaunchParams = params
	ui.DisplayText(fmt.Sprintf("The following launch parameters are set: %s.\r\n", strings.Join(params, "; ")))
	return ctx.GetPreviousState()
}

func (s *ChangeLaunchParamsState) Commands() []core.Command {
	return []core.Command{&core.BackCommand{}}
}

type VideoBackendMenuState struct{ core.BaseState }

func (s *VideoBackendMenuState) Name() string {
	return "video backend menu"
}

func NewVideoBackendMenu(ctx *core.AppContext, ui *core.UiContext) *core.MenuState {
	parrentState := &VideoBackendMenuState{}
	options := make([]*core.MenuOption, 0, 5)
	options = append(options, &core.MenuOption{
		Id:          0,
		Description: "Back.",
		NextState:   ctx.GetPreviousState,
	})
	optNum := 1
	for _, b := range game.VideoBackends {
		backend := b
		options = append(options, &core.MenuOption{
			Id:          optNum,
			Description: backend.String(),
			NextState:   func() (core.State, error) { return &SelectVideoBackendState{backend: backend}, nil },
		})
		optNum += 1
	}
	return core.NewMenu(parrentState, options, "")
}

type SelectVideoBackendState struct {
	core.BaseState
	backend game.VidBackend
}

func (s *SelectVideoBackendState) Handle(ctx *core.AppContext, ui *core.UiContext, input string) (core.State, error) {
	if err := ctx.GameManager.Params.SetVideoBackend(s.backend); err != nil {
		ui.DisplayError(err)
	} else {
		ui.DisplayText(fmt.Sprintf("You have selected video backend: %s.\r\n", s.backend))
	}
	return ctx.GetStateFromDeep(2)
}

func (s *SelectVideoBackendState) RequiresInput() bool {
	return false
}
