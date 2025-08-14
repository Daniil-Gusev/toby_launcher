package config

import (
	"toby_launcher/apperrors"
)

type gzdoomConfigData struct {
	Params      []string `json:"params"`
	DebugOutput bool     `json:"debug_output"`
	Logging     bool     `json:"logging"`
}

func (d *gzdoomConfigData) validate() error {
	errStr := "Field \"$field\" is missing."
	if d.Params == nil {
		return apperrors.New(apperrors.Err, errStr, map[string]any{
			"field": "gzdoom.params",
		})
	}
	return nil
}

type GzdoomConfig struct {
	LaunchParams []string
	DebugOutput  bool
	Logging      bool
}

func NewGzdoomConfig() *GzdoomConfig {
	return &GzdoomConfig{
		LaunchParams: make([]string, 0, 10),
	}
}

func (c *GzdoomConfig) load(data *gzdoomConfigData) error {
	if err := data.validate(); err != nil {
		return err
	}
	c.LaunchParams = data.Params
	c.Logging = data.Logging
	c.DebugOutput = data.DebugOutput
	return nil
}

func (c *GzdoomConfig) save() *gzdoomConfigData {
	data := &gzdoomConfigData{
		Params: c.LaunchParams,
	}
	if c.DebugOutput {
		data.DebugOutput = c.DebugOutput
	}
	if c.Logging {
		data.Logging = c.Logging
	}
	return data
}
