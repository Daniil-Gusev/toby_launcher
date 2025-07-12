package core

import (
	"fmt"
	"strings"
	"time"
	"toby_launcher/core/version"
)

type QuitCommand struct{ BaseCommand }

func (c *QuitCommand) Name() string {
	return "quit"
}

func (c *QuitCommand) Description() string {
	return "Immediately terminates the program."
}

func (c *QuitCommand) Aliases() []string {
	return []string{"exit", "terminate"}
}

func (c *QuitCommand) Execute(ctx *AppContext, ui *UiContext, args []string) (State, error) {
	if len(args) > 1 && args[1] == "force" {
		return &ExitState{}, nil
	}
	state, err := ctx.GetCurrentState()
	if err != nil {
		return state, nil
	}
	if _, ok := state.(*ConfirmationDialogState); ok {
		return state, nil
	}
	return NewConfirmationDialog(&ExitState{}, "Are you sure you want to immediately terminate the program?"), nil
}

type HelpCommand struct{ BaseCommand }

func (c *HelpCommand) Name() string {
	return "help"
}

func (c *HelpCommand) Description() string {
	return "Displays help information."
}

func (c *HelpCommand) Aliases() []string {
	return []string{"?", "info"}
}

func (c *HelpCommand) Execute(ctx *AppContext, ui *UiContext, args []string) (State, error) {
	state, err := ctx.GetCurrentState()
	if err != nil {
		return nil, err
	}
	desc := state.Description()
	if desc == "" {
		ui.DisplayText("Help for this state not found.\r\n")
	} else {
		ui.DisplayText(desc + "\r\n")
	}
	ui.DisplayText("The following commands are available to you:\r\n")
	for _, cmd := range ui.CommandRegistry.GetLocalCommands() {
		ui.DisplayText(fmt.Sprintf("%s: (%s).\r\n%s\r\n", cmd.Name(), strings.Join(cmd.Aliases(), ", "), cmd.Description()))
	}
	for _, cmd := range ui.CommandRegistry.GetGlobalCommands() {
		ui.DisplayText(fmt.Sprintf("%s: (%s).\r\n%s\r\n", cmd.Name(), strings.Join(cmd.Aliases(), ", "), cmd.Description()))
	}
	return state, nil
}

type BackCommand struct{ BaseCommand }

func (c *BackCommand) Name() string {
	return "back"
}

func (c *BackCommand) Description() string {
	return "Returns to the previous step."
}

func (c *BackCommand) Execute(ctx *AppContext, ui *UiContext, args []string) (State, error) {
	state, err := ctx.GetPreviousState()
	if err != nil {
		return nil, err
	}
	return state, nil
}

type VersionCommand struct{ BaseCommand }

func (c *VersionCommand) Name() string {
	return "version"
}

func (c *VersionCommand) Description() string {
	return "Displays the current version of the application."
}

func (c *VersionCommand) Execute(ctx *AppContext, ui *UiContext, args []string) (State, error) {
	var displayTime string
	builtTime, err := time.Parse(time.RFC3339, version.BuildTime)
	if err != nil {
		displayTime = version.BuildTime
	} else {
		displayTime = builtTime.Format("02.01.2006 15:04:05")
	}
	versionMsg := fmt.Sprintf("%s version: %s, built: %s.\r\n", version.AppName, version.Version, displayTime)
	ui.DisplayText(versionMsg)
	return ctx.GetCurrentState()
}

type ConfirmCommand struct{ BaseCommand }

func (c *ConfirmCommand) Name() string {
	return "confirm"
}

func (c *ConfirmCommand) Description() string {
	return "Confirms the specified action."
}

func (c *ConfirmCommand) Aliases() []string {
	return []string{"yes"}
}

func (c *ConfirmCommand) Execute(ctx *AppContext, ui *UiContext, args []string) (State, error) {
	currentState, err := ctx.GetCurrentState()
	if err != nil {
		return currentState, err
	}
	confirmationState, ok := currentState.(*ConfirmationDialogState)
	if !ok {
		previousState, _ := ctx.GetPreviousState()
		ui.DisplayText("Incorrect confirmation dialog!\r\n")
		return previousState, nil
	}
	return confirmationState.nextState, nil
}

type CancelCommand struct{ BaseCommand }

func (c *CancelCommand) Name() string {
	return "cancel"
}

func (c *CancelCommand) Description() string {
	return "Cancels the specified action."
}

func (c *CancelCommand) Aliases() []string {
	return []string{"no"}
}

func (c *CancelCommand) Execute(ctx *AppContext, ui *UiContext, args []string) (State, error) {
	return ctx.GetPreviousState()
}
