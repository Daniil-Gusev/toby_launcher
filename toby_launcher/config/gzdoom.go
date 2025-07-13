package config

import (
	"toby_launcher/apperrors"
)

type gzdoomConfigData struct {
	Params      []string `json:"params"`
	DebugOutput bool     `json:"debug_output"`
	Logging     bool     `json:"logging"`
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
	debugOutput := d.DebugOutput
	logging := d.Logging
	cfg := &GzdoomConfig{
		LaunchParams: params,
		DebugOutput:  debugOutput,
		Logging:      logging,
	}
	return cfg, nil
}

type GzdoomConfig struct {
	LaunchParams []string
	DebugOutput  bool
	Logging      bool
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
