package config

import (
	"toby_launcher/apperrors"
)

type gzdoomConfigData struct {
	Params []string `json:"params"`
}

func (d gzdoomConfigData) validate() (*GzdoomConfig, error) {
	errStr := "Field \"$field\" is missing."
	if d.Params == nil {
		return nil, apperrors.New(apperrors.Err, errStr, map[string]any{
			"field": "gzdoom.params",
		})
	}
	params := make([]string, 0, 10)
	for _, param := range d.Params {
		params = append(params, param)
	}
	cfg := &GzdoomConfig{
		LaunchParams: params,
	}
	return cfg, nil
}

type GzdoomConfig struct {
	LaunchParams []string
}

func (c *GzdoomConfig) save() *gzdoomConfigData {
	data := &gzdoomConfigData{
		Params: c.LaunchParams,
	}
	return data
}
