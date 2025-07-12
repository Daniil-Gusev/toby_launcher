package core

import (
	"fmt"
	"sort"
	"toby_launcher/core/validation"
	"toby_launcher/utils"
)

type ExitState struct{ BaseState }

func (e *ExitState) Name() string {
	return "exit"
}

func (e *ExitState) Display(_ *AppContext, ui *UiContext) {
	ui.TtsManager.Speak("Good bye!")
	ui.DisplayText("Good bye!\r\n")
}

func (e *ExitState) Handle(ctx *AppContext, _ *UiContext, _ string) (State, error) {
	ctx.AppIsRunning = false
	return e, nil
}

func (e *ExitState) RequiresInput() bool {
	return false
}

type ConfirmationDialogState struct {
	BaseState
	message   string
	nextState State
}

func NewConfirmationDialog(nextState State, message string) *ConfirmationDialogState {
	if message == "" {
		message = "Are you sure?"
	}
	return &ConfirmationDialogState{
		nextState: nextState,
		message:   message,
	}
}

func (d *ConfirmationDialogState) Name() string {
	return "confirmation dialog"
}

func (d *ConfirmationDialogState) Description() string {
	return "You are in a confirmation dialog for your last action. You need to confirm or cancel it."
}

func (d *ConfirmationDialogState) Display(_ *AppContext, ui *UiContext) {
	ui.DisplayText(d.message + "\r\n")
}

func (d *ConfirmationDialogState) Handle(ctx *AppContext, ui *UiContext, input string) (State, error) {
	ui.DisplayText("You need to confirm or cancel your choice (yes/no)." + "\r\n")
	return d, nil
}

func (d *ConfirmationDialogState) Commands() []Command {
	return []Command{
		&ConfirmCommand{},
		&CancelCommand{},
	}
}

type MenuOption struct {
	Id          int
	Description string
	Params      func() map[string]any
	NextState   func() (State, error)
}

type MenuState struct {
	BaseState
	ParentState State
	Options     []MenuOption
	OptionsMap  map[int]MenuOption
	Header      string
}

func NewMenu(parentState State, options []MenuOption, header string) *MenuState {
	sort.Slice(options, func(i, j int) bool {
		return options[i].Id < options[j].Id
	})
	optionsMap := make(map[int]MenuOption, len(options))
	for _, option := range options {
		optionsMap[option.Id] = option
	}
	return &MenuState{
		ParentState: parentState,
		Options:     options,
		OptionsMap:  optionsMap,
		Header:      header,
	}
}

func (m *MenuState) Name() string {
	return "menu"
}

func (m *MenuState) Description() string {
	return "You are in a menu. You need to enter the number of the menu item you wish to select."
}

func (m *MenuState) Display(ctx *AppContext, ui *UiContext) {
	ui.DisplayText(m.Header + "\r\n")
	for _, option := range m.Options {
		desc := option.Description
		if option.Params != nil {
			desc = utils.SubstituteParams(desc, option.Params())
		}
		ui.DisplayText(fmt.Sprintf("%d. %s\r\n", option.Id, desc))
	}
	ui.DisplayText("Make your choice.\r\n")
}

func (m *MenuState) Handle(ctx *AppContext, ui *UiContext, input string) (State, error) {
	num, err := validation.ParseInt(input)
	if err != nil {
		return m, err
	}
	option, exists := m.OptionsMap[num]
	if !exists {
		ui.DisplayText("There is no such item in the menu.\r\n")
		return m, nil
	}
	return option.NextState()
}
