package config

import (
	"toby_launcher/apperrors"
)

type gzdoomConfigData struct {
	Params      []string `json:"params"`
	DebugOutput bool
}

func (d gzdoomConfigData) validate() (*GzdoomConfig, error) {
	errStr := "Field \"$field\" is missing."
	if d.Params == nil {
		return nil, apperrors.New(apperrors.Err, errStr, map[string]any{
			"field": "gzdoom.params",
		})
	}
	params := make([]string, 0, 10)
	params = append(params, d.Params...)
	debugOutput := false
	if d.DebugOutput {
		debugOutput = true
	}
	cfg := &GzdoomConfig{
		LaunchParams: params,
		DebugOutput:  debugOutput,
	}
	return cfg, nil
}

type GzdoomConfig struct {
	LaunchParams []string
	DebugOutput  bool
}

func (c *GzdoomConfig) save() *gzdoomConfigData {
	data := &gzdoomConfigData{
		Params: c.LaunchParams,
	}
	if c.DebugOutput {
		data.DebugOutput = c.DebugOutput
	}
	return data
}
