package tts

import (
	"fmt"
	"time"
	"toby_launcher/apperrors"
	"toby_launcher/config"
	"toby_launcher/core/logger"
)

type Phrase struct {
	Id      int
	Text    string
	Rate    int
	Silence int
}

type SpeechSynthesizer interface {
	Name() string
	CreateNew() (SpeechSynthesizer, error)
	Release()
	Speak(*Phrase) error
	Stop() error
	IsSpeaking() (bool, error)
	SetSpeechRate(rate int) error
	GetSpeechRate() int
	SupportsChangingSpeechRate() bool
	SetLogger(l logger.Logger)
	LogError(err error)
}

type BaseSynthesizer struct {
	logger logger.Logger
}

func (b *BaseSynthesizer) SupportsChangingSpeechRate() bool {
	return true
}

func (b *BaseSynthesizer) SetLogger(l logger.Logger) {
	b.logger = l
}

func (b *BaseSynthesizer) LogError(err error) {
	if b.logger != nil {
		b.logger.Error(err)
	}
}

type TtsManager struct {
	logger                logger.Logger
	currentSynthesizer    SpeechSynthesizer
	availableSynthesizers []SpeechSynthesizer
	phraseCounter         int
	config                *config.TtsConfig
}

func NewTtsManager(cfg *config.TtsConfig, logger logger.Logger) (*TtsManager, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger not specified")
	}
	syns := GetAvailableSynthesizers(logger)
	if len(syns) == 0 {
		return nil, apperrors.New(apperrors.ErrSpeech, "No available speech synthesizers.", nil)
	}
	manager := &TtsManager{
		logger:                logger,
		availableSynthesizers: syns,
		phraseCounter:         0,
		config:                cfg,
	}
	if err := manager.ApplyConfig(); err != nil {
		return nil, err
	}
	return manager, nil
}

func (m *TtsManager) Release() {
	if m.currentSynthesizer != nil {
		m.currentSynthesizer.Release()
		m.currentSynthesizer = nil
	}
	m.availableSynthesizers = nil
}

func (m *TtsManager) Wait(timeout int) error {
	waiting := 0
	waitMs := 200
	if m.currentSynthesizer == nil {
		return nil
	}
	for {
		isSpeaking, err := m.currentSynthesizer.IsSpeaking()
		if err != nil {
			return err
		}
		if !isSpeaking {
			return nil
		}
		time.Sleep(time.Duration(waitMs) * time.Millisecond)
		waiting += waitMs
		if waiting >= timeout {
			return nil
		}
	}
}

func (m *TtsManager) NewPhrase(text string, rate, silence int) *Phrase {
	m.phraseCounter++
	return &Phrase{
		Id:      m.phraseCounter,
		Text:    text,
		Rate:    rate,
		Silence: silence,
	}
}

func (m *TtsManager) Speak(text string) {
	phrase := m.NewPhrase(text, 0, 0)
	m.SpeakPhrase(phrase)
}

func (m *TtsManager) SpeakPhrase(phrase *Phrase) {
	if m.currentSynthesizer == nil {
		m.logger.Error(apperrors.New(apperrors.ErrSpeech, "No speech synthesizer is initialized.", nil))
		return
	}
	if err := m.currentSynthesizer.Speak(phrase); err != nil {
		m.logger.Error(apperrors.New(apperrors.ErrSpeech, err.Error(), nil))
	}
}

func (m *TtsManager) ApplyConfig() error {
	baseSynthName := m.availableSynthesizers[0].Name()
	synthName := m.config.SynthesizerName
	if m.currentSynthesizer == nil && synthName == "" {
		synthName = baseSynthName
	}
	if err := m.SetSynthesizer(synthName); err != nil {
		m.logger.Error(err)
		if synthName != baseSynthName {
			if err := m.SetSynthesizer(baseSynthName); err != nil {
				return err
			}
		}
	}
	if m.currentSynthesizer.SupportsChangingSpeechRate() {
		if err := m.SetSpeechRate(m.config.SpeechRate); err != nil {
			return err
		}
	}
	return nil
}

func (m *TtsManager) SetSynthesizer(synthName string) error {
	if m.currentSynthesizer != nil && m.currentSynthesizer.Name() == synthName {
		return nil
	}
	synth, isFound := m.findSynthesizer(synthName)
	if !isFound {
		return apperrors.New(apperrors.ErrSpeech, "The speech synthesizer \"$synthesizer\" is missing.", map[string]any{"synthesizer": synthName})
	}
	if m.currentSynthesizer != nil {
		m.currentSynthesizer.Release()
	}
	newSynth, err := synth.CreateNew()
	if err != nil {
		return apperrors.New(apperrors.ErrSpeech, err.Error(), nil)
	}
	if m.config.SpeechRate > 0 && newSynth.SupportsChangingSpeechRate() {
		if err := newSynth.SetSpeechRate(m.config.SpeechRate); err != nil {
			return apperrors.New(apperrors.Err, err.Error(), nil)
		}
	}
	newSynth.SetLogger(m.logger)
	m.currentSynthesizer = newSynth
	m.config.SynthesizerName = newSynth.Name()
	m.phraseCounter = 0
	return nil
}

func (m *TtsManager) SetSpeechRate(rate int) error {
	if rate < 0 || rate == m.currentSynthesizer.GetSpeechRate() {
		return nil
	}
	if rate == 0 {
		rate = m.currentSynthesizer.GetSpeechRate()
	}
	if err := m.currentSynthesizer.SetSpeechRate(rate); err != nil {
		return err
	}
	m.config.SpeechRate = m.currentSynthesizer.GetSpeechRate()
	return nil
}

func (m *TtsManager) findSynthesizer(synthName string) (SpeechSynthesizer, bool) {
	var synth SpeechSynthesizer
	isFound := false
	for _, s := range m.availableSynthesizers {
		if s.Name() == synthName {
			synth = s
			isFound = true
			break
		}
	}
	return synth, isFound
}

func (m *TtsManager) AvailableSynthesizers() []SpeechSynthesizer {
	return m.availableSynthesizers
}
