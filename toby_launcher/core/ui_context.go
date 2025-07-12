package core

import (
	"fmt"
	"strings"
	"toby_launcher/apperrors"
	"toby_launcher/core/logger"
	"toby_launcher/core/tts"
	"toby_launcher/utils"
)

type UiContext struct {
	Console         Console
	ErrorHandler    apperrors.ErrorHandler
	CommandRegistry *CommandRegistry
	Logger          logger.Logger
	TtsManager      *tts.TtsManager
}

func (ui *UiContext) DisplayText(txt string) {
	if err := ui.Console.Write(utils.WrapText(txt, 80)); err != nil {
		fmt.Println(ui.ErrorHandler.Handle(err))
	}
}

func (ui *UiContext) DisplayError(err error) {
	msg := ui.ErrorHandler.Handle(err)
	if msg != "" {
		ui.DisplayText(fmt.Sprintf("%s\r\n", msg))
	}
}

func (ui *UiContext) HandleInput(input string, ctx *AppContext) (State, error) {
	input = strings.TrimSpace(input)
	if cmd, args := ui.CommandRegistry.ParseInput(input); cmd != nil {
		return cmd.Execute(ctx, ui, args)
	}
	state, err := ctx.GetCurrentState()
	if err != nil {
		return nil, err
	}
	return state.Handle(ctx, ui, input)
}
