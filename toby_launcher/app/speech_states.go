package app

import (
	"fmt"
	"toby_launcher/core"
	"toby_launcher/core/validation"
)

type SpeechSettingsMenuState struct{ core.BaseState }

func (m *SpeechSettingsMenuState) Name() string {
	return "speech settings menu"
}

func NewSpeechSettingsMenu(ctx *core.AppContext, ui *core.UiContext) *core.MenuState {
	parrentState := &SpeechSettingsMenuState{}
	options := []*core.MenuOption{
		{Id: 0,
			Description: "Back",
			NextState:   ctx.GetPreviousState,
		},
		{Id: 1,
			Description: "Change speech synthesizer ($synthesizer).",
			Params: func() map[string]any {
				return map[string]any{"synthesizer": ctx.Config.Tts.SynthesizerName}
			},
			NextState: func() (core.State, error) { return NewSynthesizerSelectionMenu(ctx, ui), nil },
		},
		{Id: 2,
			Description: "Change speech rate ($rate).",
			Params: func() map[string]any {
				return map[string]any{"rate": ctx.Config.Tts.SpeechRate}
			},
			NextState: func() (core.State, error) { return &SelectSpeechRateState{}, nil },
		},
	}
	return core.NewMenu(parrentState, options, "")
}

type SynthesizerSelectionMenuState struct{ core.BaseState }

func (m *SynthesizerSelectionMenuState) Name() string {
	return "list synthesizers menu"
}

func NewSynthesizerSelectionMenu(ctx *core.AppContext, ui *core.UiContext) *core.MenuState {
	parentState := &SynthesizerSelectionMenuState{}
	syns := ui.TtsManager.AvailableSynthesizers()
	options := make([]*core.MenuOption, 0, 1+len(syns))
	options = append(options, &core.MenuOption{
		Id:          0,
		Description: "Back.",
		NextState:   ctx.GetPreviousState,
	})
	for i, s := range syns {
		synth := s
		opt := &core.MenuOption{
			Id:          i + 1,
			Description: synth.Name() + ".",
			NextState: func() (core.State, error) {
				if err := ui.TtsManager.SetSynthesizer(synth.Name()); err != nil {
					ui.DisplayError(err)
					return ctx.GetCurrentState()
				}
				msg := fmt.Sprintf("%s selected.", synth.Name())
				ui.DisplayText(msg + "\r\n")
				ui.TtsManager.Speak(msg)
				return ctx.GetCurrentState()
			},
		}
		options = append(options, opt)
	}
	return core.NewMenu(parentState, options, "")
}

type SelectSpeechRateState struct{ core.BaseState }

func (s *SelectSpeechRateState) Name() string {
	return "select speech rate"
}

func (s *SelectSpeechRateState) Description() string {
	return "You need to enter the number of words that corresponds to your desired speech speed in words per minute."
}

func (s *SelectSpeechRateState) Display(ctx *core.AppContext, ui *core.UiContext) {
	ui.DisplayText("Enter your desired speech rate.\r\n")
	if ctx.Config.Tts.SpeechRate > 0 {
		ui.DisplayText(fmt.Sprintf("Current value: %d.\r\n", ctx.Config.Tts.SpeechRate))
	}
}

func (s *SelectSpeechRateState) Handle(ctx *core.AppContext, ui *core.UiContext, input string) (core.State, error) {
	rate, err := validation.ParseIntInRange(input, 0, 1000)
	if err != nil {
		return s, err
	}
	if err := ui.TtsManager.SetSpeechRate(rate); err != nil {
		ui.DisplayError(err)
		return ctx.GetPreviousState()
	}
	msg := fmt.Sprintf("You have selected speech rate: %d.\r\n", ctx.Config.Tts.SpeechRate)
	ui.DisplayText(msg)
	ui.TtsManager.Speak(msg)
	return ctx.GetPreviousState()
}

func (s *SelectSpeechRateState) Commands() []core.Command {
	return []core.Command{
		&core.BackCommand{},
	}
}
