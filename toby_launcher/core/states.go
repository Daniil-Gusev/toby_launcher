package core

import (
	"fmt"
	"sort"
	"toby_launcher/apperrors"
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
	parentState State
	options     []*MenuOption
	optionsMap  map[int]*MenuOption
	header      string
}

func NewMenu(parentState State, options []*MenuOption, header string) *MenuState {
	sortedOptions := make([]*MenuOption, 0, len(options))
	for _, o := range options {
		if o != nil {
			sortedOptions = append(sortedOptions, o)
		}
	}
	sort.Slice(sortedOptions, func(i, j int) bool {
		return sortedOptions[i].Id < sortedOptions[j].Id
	})
	optionsMap := make(map[int]*MenuOption, len(sortedOptions))
	for _, option := range sortedOptions {
		optionsMap[option.Id] = option
	}
	return &MenuState{
		parentState: parentState,
		options:     sortedOptions,
		optionsMap:  optionsMap,
		header:      header,
	}
}

func (m *MenuState) Name() string {
	if m.parentState != nil {
		return m.parentState.Name()
	}
	return "menu"
}

func (m *MenuState) Description() string {
	return "You are in a menu. You need to enter the number of the menu item you wish to select."
}

func (m *MenuState) Display(ctx *AppContext, ui *UiContext) {
	ui.DisplayText(m.header + "\r\n")
	for _, option := range m.options {
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
	option, exists := m.optionsMap[num]
	if !exists {
		ui.DisplayText("There is no such item in the menu.\r\n")
		return m, nil
	}
	return option.NextState()
}

type OptionSwitcher bool

func (s OptionSwitcher) String() string {
	switch s {
	case true:
		return "enable"
	default:
		return "disable"
	}
}

type SwitchOptionState struct {
	BaseState
	option       *bool
	name         string
	changeOption func()
}

func (s *SwitchOptionState) Init(ctx *AppContext, ui *UiContext) (State, error) {
	if s.option == nil {
		err := apperrors.New(apperrors.Err, "Option \"$option\" is not specified.", map[string]any{"option": s.name})
		ui.DisplayError(err)
		return ctx.GetPreviousState()
	}
	return s, nil
}

func (s *SwitchOptionState) Handle(ctx *AppContext, ui *UiContext, input string) (State, error) {
	switcher := OptionSwitcher(*s.option)
	ui.DisplayText(fmt.Sprintf("%s is %vd.\r\n", s.name, !switcher))
	*s.option = !bool(switcher)
	return ctx.GetPreviousState()
}

func (s *SwitchOptionState) RequiresInput() bool {
	return false
}

func NewSwitchMenuOption(id int, name string, option *bool) *MenuOption {
	if option == nil {
		return nil
	}
	return &MenuOption{
		Id:          id,
		Description: "$action $option.",
		Params: func() map[string]any {
			return map[string]any{"action": !OptionSwitcher(*option), "option": name}
		},
		NextState: func() (State, error) {
			return &SwitchOptionState{option: option, name: name}, nil
		},
	}
}
