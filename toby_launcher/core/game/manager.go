package game

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"toby_launcher/apperrors"
	"toby_launcher/config"
	"toby_launcher/core/logger"
	"toby_launcher/core/tts"
	"toby_launcher/utils/file_utils"
)

// GameManager manages game profiles and their execution.
type GameManager struct {
	logger        logger.Logger
	config        *config.Config
	tts           *tts.TtsManager
	games         []*GameData
	currentGame   *Game
	textProcessor *TextProcessor
}

// NewGameManager creates a new GameManager instance.
func NewGameManager(cfg *config.Config, logger logger.Logger, tts *tts.TtsManager) (*GameManager, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger not specified")
	}
	manager := &GameManager{
		logger:        logger,
		config:        cfg,
		tts:           tts,
		games:         make([]*GameData, 0, 10),
		textProcessor: NewTextProcessor(cfg, logger, tts),
	}
	if err := manager.loadGames(); err != nil {
		return nil, err
	}
	return manager, nil
}

// loadGames loads game profiles from games.json.
func (m *GameManager) loadGames() error {
	gamesPath := m.config.Paths.GamesPath()
	var gamesData RawGamesData
	if err := file_utils.LoadData(gamesPath, &gamesData); err != nil {
		return apperrors.New(apperrors.Err, "Failed to load games: $error", map[string]any{"error": err})
	}
	for n, g := range gamesData {
		game := &GameData{
			Name:        n,
			Description: g.Description,
			Config:      g.Config,
			Iwad:        g.Iwad,
			Files:       g.Files,
			Params:      g.Params,
		}
		m.games = append(m.games, game)
	}
	m.sortGames()
	return nil
}

func (m *GameManager) sortGames() {
	sort.Slice(m.games, func(i, j int) bool { return m.games[i].Name < m.games[j].Name })
}

// AvailableGames returns the list of available games.
func (m *GameManager) AvailableGames() []*GameData {
	return m.games
}

// StartGame launches a game with the given name.
func (m *GameManager) StartGame(gameData *GameData) error {
	if m.currentGame != nil && m.currentGame.IsRunning {
		return apperrors.New(apperrors.Err, "Another game is already running", nil)
	}
	gzdoomPath, err := m.config.Paths.GzdoomPath()
	if err != nil {
		return apperrors.New(apperrors.Err, "Failed to find gzdoom: $error", map[string]any{"error": err})
	}
	args := m.buildGameArgs(gameData)
	cmd := exec.Command(gzdoomPath, args...)
	cmd.Stdout = m.textProcessor // TextProcessor will handle output
	cmd.Stderr = m.textProcessor
	cmd.Env = append(os.Environ(), fmt.Sprintf("DOOMWADDIR=%s", m.config.Paths.FilesDir))
	m.currentGame = &Game{
		Info:      gameData,
		cmd:       cmd,
		IsRunning: true,
	}
	m.logger.DebugPrintf("Running %v\r\n", strings.Join(cmd.Args, " "))
	if err := cmd.Start(); err != nil {
		m.currentGame = nil
		return apperrors.New(apperrors.Err, "Failed to start game: $error", map[string]any{"error": err})
	}
	go m.handleGameProcess()
	return nil
}

// StopGame stops the currently running game.
func (m *GameManager) StopGame() error {
	if m.currentGame == nil || !m.currentGame.IsRunning {
		return nil
	}
	if err := m.currentGame.cmd.Process.Kill(); err != nil {
		return apperrors.New(apperrors.Err, "Failed to stop game: $error", map[string]any{"error": err})
	}
	m.currentGame.IsRunning = false
	m.currentGame = nil
	return nil
}

func (m *GameManager) GameIsRunning() bool {
	if m.currentGame == nil {
		return false
	}
	return m.currentGame.IsRunning
}

// handleGameProcess manages the game process lifecycle and output processing.
func (m *GameManager) handleGameProcess() {
	if m.currentGame == nil {
		return
	}
	err := m.currentGame.cmd.Wait()
	if m.currentGame != nil {
		if err != nil && m.currentGame.IsRunning {
			m.logger.Error(apperrors.New(apperrors.Err, "Game process error: $error", map[string]any{"error": err}))
		}
		m.tts.Speak("Game finished.")
		m.logger.Printf("Game finished.\r\n")
		m.currentGame.IsRunning = false
		m.currentGame = nil
	}
}

// buildGameArgs constructs the command-line arguments for gzdoom.
func (m *GameManager) buildGameArgs(data *GameData) []string {
	args := make([]string, 0, 5+len(data.Files)*2+len(data.Params)*2+len(m.config.Gzdoom.LaunchParams)*2)
	args = append(args, "-stdout")
	if m.config.Gzdoom.Logging {
		args = append(args, "+logfile", m.config.Paths.GzdoomLogFilePath())
	}
	if len(m.config.Gzdoom.LaunchParams) > 0 {
		for _, param := range m.config.Gzdoom.LaunchParams {
			args = append(args, strings.Split(param, " ")...)
		}
	}
	if len(data.Params) > 0 {
		for _, param := range data.Params {
			args = append(args, strings.Split(param, " ")...)
		}
	}
	if data.Config != "" {
		configPath := m.config.Paths.GameFilePath(data.Config)
		if file_utils.Exists(configPath) {
			args = append(args, "-config", configPath)
		} else {
			m.logger.Printf("Warning: configuration file %s for game %s is not found.\r\n", configPath, data.Name)
		}
	}
	if data.Iwad != "" {
		iwadPath := m.config.Paths.GameFilePath(data.Iwad)
		if file_utils.Exists(iwadPath) {
			args = append(args, "-iwad", data.Iwad)
		} else {
			m.logger.Printf("Warning: iwad file %s for game %s is not found.\r\n", iwadPath, data.Name)
		}
	}
	files := make([]string, 0, len(data.Files))
	for _, file := range data.Files {
		filePath := m.config.Paths.GameFilePath(file)
		if file_utils.Exists(filePath) {
			files = append(files, file)
		} else {
			m.logger.Printf("Warning: additional file %s for game %s is not found.\r\n", filePath, data.Name)
		}
	}
	if len(files) > 0 {
		args = append(args, "-file")
		args = append(args, files...)
	}
	return args
}
