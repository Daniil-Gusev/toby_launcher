package config

import (
	"toby_launcher/apperrors"
	"toby_launcher/utils/file_utils"
)

type configData struct {
	Tts    *ttsConfigData    `json:"tts"`
	Gzdoom *gzdoomConfigData `json:"gzdoom"`
}

func (d configData) validate() error {
	if d.Tts == nil {
		return apperrors.New(apperrors.Err, "required section \"$section\" is missing.", map[string]any{
			"section": "tts",
		})
	}
	if d.Gzdoom == nil {
		return apperrors.New(apperrors.Err, "required section \"$section\" is missing.", map[string]any{
			"section": "gzdoom",
		})
	}
	return nil
}

// Config contains all application configuration settings.
type Config struct {
	Paths  *PathConfig
	Tts    *TtsConfig
	Gzdoom *GzdoomConfig
}

// NewConfig creates a new Config instance with initialized PathConfig and LanguageConfig.
func NewConfig() (*Config, error) {
	pathConfig, err := NewPathConfig()
	if err != nil {
		return nil, err
	}
	return &Config{
		Paths:  pathConfig,
		Tts:    &TtsConfig{},
		Gzdoom: &GzdoomConfig{},
	}, nil
}

func (c *Config) Load(filePath string) error {
	var rawData configData
	if err := file_utils.LoadData(filePath, &rawData); err != nil {
		return apperrors.New(apperrors.Err, "Error parsing file $file: $error", map[string]any{
			"file":  filePath,
			"error": err,
		})
	}
	if err := rawData.validate(); err != nil {
		return apperrors.New(apperrors.Err, "Error in configuration file $file: $error", map[string]any{
			"file":  filePath,
			"error": err,
		})
	}
	ttsCfg, err := rawData.Tts.validate()
	if err != nil {
		return apperrors.New(apperrors.Err, "error in configuration file $file: $error", map[string]any{
			"file":  filePath,
			"error": err,
		})
	}
	gzdoomCfg, err := rawData.Gzdoom.validate()
	if err != nil {
		return apperrors.New(apperrors.Err, "error in configuration file $file: $error", map[string]any{
			"file":  filePath,
			"error": err,
		})
	}
	c.Tts = ttsCfg
	c.Gzdoom = gzdoomCfg
	return nil
}

func (c *Config) Save() error {
	var cfgData configData
	cfgData.Tts = c.Tts.save()
	cfgData.Gzdoom = c.Gzdoom.save()
	if err := file_utils.SaveData(c.Paths.ConfigFilePath(), cfgData); err != nil {
		return err
	}
	return nil
}
