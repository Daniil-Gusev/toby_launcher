package config

import (
	"toby_launcher/apperrors"
	"toby_launcher/core/validation"
)

type ttsConfigData struct {
	SpeechEngine *string `json:"speech_engine"`
	Rate         *int    `json:"rate"`
}

func (d ttsConfigData) validate() (*TtsConfig, error) {
	errStr := "Field \"$field\" is missing."
	if d.SpeechEngine == nil {
		return nil, apperrors.New(apperrors.Err, errStr, map[string]any{
			"field": "tts.speech_engine",
		})
	}
	if d.Rate == nil {
		return nil, apperrors.New(apperrors.Err, errStr, map[string]any{
			"field": "tts.rate",
		})
	}
	rate := 0
	newRate := *d.Rate
	if _, err := validation.IsNumInRange(newRate, 0, 1000); err != nil {
		return nil, apperrors.New(apperrors.Err, "invalid value in field \"$field\": $error", map[string]any{
			"field": "tts.rate",
			"error": err,
		})
	}
	rate = newRate
	cfg := &TtsConfig{
		SynthesizerName: *d.SpeechEngine,
		SpeechRate:      rate,
	}
	return cfg, nil
}

type TtsConfig struct {
	SynthesizerName string
	SpeechRate      int
}

func (c *TtsConfig) save() *ttsConfigData {
	engine := c.SynthesizerName
	rate := c.SpeechRate
	data := &ttsConfigData{
		SpeechEngine: &engine,
		Rate:         &rate,
	}
	return data
}
