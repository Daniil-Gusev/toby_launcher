package main

import (
	"fmt"
	"os"
	"toby_launcher/app"
	"toby_launcher/apperrors"
	"toby_launcher/config"
	"toby_launcher/core"
	"toby_launcher/core/game"
	"toby_launcher/core/logger"
	"toby_launcher/core/tts"
	_ "toby_launcher/speech_engines"
)

func main() {
	errorHandler := &apperrors.StdErrorHandler{}
	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Printf("Failed to initialize configuration: %v\r\n", err)
		return
	}
	logger, err := logger.NewStdLogger(os.Stdout, cfg.Paths.LogFilePath(), errorHandler)
	if err != nil {
		fmt.Printf("Faled to initialize logger: %v\r\n", err)
		return
	}
	defer logger.Release()
	if err := cfg.Load(cfg.Paths.ConfigFilePath()); err != nil {
		logger.Error(err)
	}
	defer func() {
		if err := cfg.Save(); err != nil {
			logger.Error(err)
		}
	}()
	if _, err := cfg.Paths.GzdoomPath(); err != nil {
		logger.Error(err)
		return
	}
	console, err := core.NewReadlineConsole()
	if err != nil {
		logger.Printf("Failed to initialize console: %v\r\n", err)
		return
	}
	defer func() {
		if consoleErr := console.Close(); consoleErr != nil && err == nil {
			err = consoleErr
		}
	}()
	ttsManager, err := tts.NewTtsManager(cfg.Tts, logger)
	if err != nil {
		logger.Error(err)
		return
	}
	defer func() {
		if err := ttsManager.Wait(5000); err != nil {
			logger.Error(err)
		}
		ttsManager.Release()
	}()
	gameManager, err := game.NewGameManager(cfg, logger, ttsManager)
	if err != nil {
		logger.Error(err)
		return
	}
	defer func() {
		if err := gameManager.StopGame(); err != nil {
			logger.Error(err)
		}
	}()
	appCtx := &core.AppContext{
		Config:       cfg,
		StateStack:   core.NewStateStack(),
		GameManager:  gameManager,
		AppIsRunning: true,
	}
	uiCtx := &core.UiContext{
		Console:         console,
		ErrorHandler:    errorHandler,
		CommandRegistry: core.NewCommandRegistry(),
		TtsManager:      ttsManager,
	}
	uiCtx.CommandRegistry.RegisterGlobalCommands(core.DefaultGlobalCommands())
	startState := core.State(&app.StartState{})
	runMainLoop(appCtx, uiCtx, startState)
}

func runMainLoop(appCtx *core.AppContext, uiCtx *core.UiContext, startState core.State) {
	currentState, err := appCtx.GoToState(startState, uiCtx)
	uiCtx.DisplayError(err)
	if currentState == nil {
		return
	}
	for appCtx.AppIsRunning {
		currentState.Display(appCtx, uiCtx)
		input := ""
		if currentState.RequiresInput() {
			buf, inputErr := uiCtx.Console.Read()
			uiCtx.DisplayError(inputErr)
			if appErr, ok := inputErr.(*apperrors.AppError); ok && appErr.Code == apperrors.ErrEOF {
				currentState, err := appCtx.GoToState(&core.ExitState{}, uiCtx)
				uiCtx.DisplayError(err)
				currentState.Display(appCtx, uiCtx)
			}
			input = buf
		}
		nextState, err := uiCtx.HandleInput(input, appCtx)
		uiCtx.DisplayError(err)
		if appErr, ok := err.(*apperrors.AppError); ok && appErr.Code == apperrors.ErrStateStack {
			currentState, err := appCtx.GoToState(startState, uiCtx)
			uiCtx.DisplayError(err)
			currentState.Display(appCtx, uiCtx)
		}
		if nextState != currentState {
			if nextState == startState {
				appCtx.StateStack.Clear()
			}
			currentState, err = appCtx.GoToState(nextState, uiCtx)
			uiCtx.DisplayError(err)
		}
	}
}
