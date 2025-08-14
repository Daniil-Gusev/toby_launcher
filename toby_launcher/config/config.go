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
	errStr := "required section \"$section\" is missing"
	if d.Tts == nil {
		return apperrors.New(apperrors.Err, errStr, map[string]any{
			"section": "tts",
		})
	}
	if d.Gzdoom == nil {
		return apperrors.New(apperrors.Err, errStr, map[string]any{
			"section": "gzdoom",
		})
	}
	return nil
}

type Config struct {
	Paths  *PathConfig
	Tts    *TtsConfig
	Gzdoom *GzdoomConfig
}

func NewConfig() (*Config, error) {
	pathConfig, err := NewPathConfig()
	if err != nil {
		return nil, err
	}
	return &Config{
		Paths:  pathConfig,
		Tts:    &TtsConfig{},
		Gzdoom: NewGzdoomConfig(),
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
	if err := c.Tts.load(rawData.Tts); err != nil {
		return apperrors.New(apperrors.Err, "error in configuration file $file: $error", map[string]any{
			"file":  filePath,
			"error": err,
		})
	}
	if err := c.Gzdoom.load(rawData.Gzdoom); err != nil {
		return apperrors.New(apperrors.Err, "error in configuration file $file: $error", map[string]any{
			"file":  filePath,
			"error": err,
		})
	}
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
