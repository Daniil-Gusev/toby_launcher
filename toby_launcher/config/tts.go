package config

import (
	"toby_launcher/apperrors"
	"toby_launcher/core/validation"
)

type ttsConfigData struct {
	SpeechEngine *string `json:"speech_engine"`
	Rate         *int    `json:"rate"`
}

func (d *ttsConfigData) validate() error {
	errStr := "Field \"$field\" is missing."
	if d.SpeechEngine == nil {
		return apperrors.New(apperrors.Err, errStr, map[string]any{
			"field": "tts.speech_engine",
		})
	}
	if d.Rate == nil {
		return apperrors.New(apperrors.Err, errStr, map[string]any{
			"field": "tts.rate",
		})
	}
	rate := *d.Rate
	if _, err := validation.IsNumInRange(rate, 0, 1000); err != nil {
		return apperrors.New(apperrors.Err, "invalid value in field \"$field\": $error", map[string]any{
			"field": "tts.rate",
			"error": err,
		})
	}
	return nil
}

type TtsConfig struct {
	SynthesizerName string
	SpeechRate      int
}

func (c *TtsConfig) load(data *ttsConfigData) error {
	if err := data.validate(); err != nil {
		return err
	}
	c.SpeechRate = *data.Rate
	c.SynthesizerName = *data.SpeechEngine
	return nil
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
