//go:build windows

package nvda

import (
	"fmt"
	"time"
	"toby_launcher/core/tts"
)

func init() {
	tts.RegisterSynthesizer(NewSynthesizer, 1)
}

type Synthesizer struct {
	tts.BaseSynthesizer
	synth      *nvdaSynthesizer
	isSpeaking bool
}

func NewSynthesizer() (tts.SpeechSynthesizer, error) {
	synth, err := newNvdaSynthesizer()
	if err != nil {
		return nil, err
	}
	return &Synthesizer{
		synth:      synth,
		isSpeaking: false,
	}, nil
}

func (s *Synthesizer) Name() string {
	return "nvda"
}

func (s *Synthesizer) CreateNew() (tts.SpeechSynthesizer, error) {
	return NewSynthesizer()
}

func (s *Synthesizer) Release() {
	if s.synth != nil {
		if err := s.Stop(); err != nil {
			s.LogError(err)
		}
		s.synth.free()
		s.synth = nil
	}
}

func (s *Synthesizer) Stop() error {
	if s.synth == nil {
		return fmt.Errorf("Nvda synthesizer is not initialized")
	}
	if err := s.synth.stop(); err != nil {
		return err
	}
	s.isSpeaking = false
	return nil
}

func (s *Synthesizer) Speak(phrase *tts.Phrase) error {
	if s.synth == nil {
		return fmt.Errorf("Nvda synthesizer is not initialized")
	}
	if phrase.Text == "" {
		return fmt.Errorf("no text to speak has been specified")
	}
	if phrase.Silence > 0 {
		phraseCopy := *phrase
		go func(phrase *tts.Phrase) {
			time.Sleep(time.Duration(phrase.Silence) * time.Millisecond)
			phrase.Silence = 0
			if err := s.Speak(phrase); err != nil {
				s.LogError(err)
			}
		}(&phraseCopy)
		return nil
	}
	if err := s.synth.stop(); err != nil {
		return err
	}
	s.isSpeaking = true
	go func() {
		if err := s.synth.speak(phrase.Text); err != nil {
			if s.isSpeaking {
				s.LogError(err)
			}
		}
		s.isSpeaking = false
	}()
	return nil
}

func (s *Synthesizer) SetSpeechRate(rate int) error {
	return fmt.Errorf("The speech rate can be changed in the nvda settings.")
}

func (s *Synthesizer) GetSpeechRate() int {
	return 0
}

func (s *Synthesizer) SupportsChangingSpeechRate() bool {
	return false
}

func (s *Synthesizer) IsSpeaking() (bool, error) {
	if s.synth == nil {
		return false, fmt.Errorf("Nvda synthesizer is not initialized")
	}
	return s.isSpeaking, nil
}
