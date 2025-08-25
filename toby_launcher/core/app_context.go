package core

import (
	"toby_launcher/apperrors"
	"toby_launcher/config"
	"toby_launcher/core/game"
)

type AppContext struct {
	Config       *config.Config
	StateStack   *StateStack
	GameManager  *game.GameManager
	AppIsRunning bool
}

func (app *AppContext) GetCurrentState() (State, error) {
	if app.StateStack.IsEmpty() {
		return nil, apperrors.New(apperrors.ErrStateStack, "state stack is empty.", nil)
	}
	return app.StateStack.Peek(), nil
}

func (app *AppContext) GetPreviousState() (State, error) {
	return app.GetStateFromDeep(1)
}

func (app *AppContext) GetStateFromDeep(depth int) (State, error) {
	if app.StateStack.IsEmpty() {
		return nil, apperrors.New(apperrors.ErrStateStack, "state stack is empty", nil)
	}
	if len(app.StateStack.states) < 1+depth {
		return nil, apperrors.New(apperrors.ErrStateStack, "state stack is insufficient.", nil)
	}
	app.StateStack.Pop()
	currentDepth := 0
	for {
		state := app.StateStack.Pop()
		if state.RequiresInput() {
			currentDepth += 1
			if currentDepth < depth {
				continue
			}
			return state, nil
		}
	}
}

func (app *AppContext) GoToState(nextState State, ui *UiContext) (State, error) {
	app.StateStack.Push(nextState)
	if newState, err := nextState.Init(app, ui); err != nil {
		if newState != nextState {
			app.StateStack.Pop()
			app.StateStack.Push(newState)
		}
		return newState, err
	}
	ui.CommandRegistry.RegisterLocalCommands(nextState.Commands())
	return nextState, nil
}
