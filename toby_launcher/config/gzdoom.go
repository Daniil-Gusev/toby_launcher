package config

import (
	"toby_launcher/apperrors"
)

type GzdoomParams map[string]any

type gzdoomConfigData struct {
	Params           GzdoomParams `json:"params"`
	AdditionalParams []string     `json:"additional_params"`
	DebugOutput      bool         `json:"debug_output"`
	Logging          bool         `json:"logging"`
}

func (d *gzdoomConfigData) validate() error {
	errStr := "Field \"$field\" is missing."
	if d.Params == nil {
		return apperrors.New(apperrors.Err, errStr, map[string]any{
			"field": "gzdoom.params",
		})
	}
	if d.AdditionalParams == nil {
		return apperrors.New(apperrors.Err, errStr, map[string]any{
			"field": "gzdoom.additional_params",
		})
	}

	return nil
}

type GzdoomConfig struct {
	GameParams             GzdoomParams
	AdditionalLaunchParams []string
	DebugOutput            bool
	Logging                bool
}

func NewGzdoomConfig() *GzdoomConfig {
	return &GzdoomConfig{
		GameParams:             make(map[string]any, 5),
		AdditionalLaunchParams: make([]string, 0, 5),
	}
}

func (c *GzdoomConfig) load(data *gzdoomConfigData) error {
	if err := data.validate(); err != nil {
		return err
	}
	c.GameParams = data.Params
	c.AdditionalLaunchParams = data.AdditionalParams
	c.Logging = data.Logging
	c.DebugOutput = data.DebugOutput
	return nil
}

func (c *GzdoomConfig) save() *gzdoomConfigData {
	data := &gzdoomConfigData{
		Params:           c.GameParams,
		AdditionalParams: c.AdditionalLaunchParams,
	}
	if c.DebugOutput {
		data.DebugOutput = c.DebugOutput
	}
	if c.Logging {
		data.Logging = c.Logging
	}
	return data
}
