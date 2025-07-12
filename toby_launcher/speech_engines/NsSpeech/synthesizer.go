//go:build darwin

package NsSpeech

import (
	"fmt"
	"math"
	"toby_launcher/core/tts"
)

func init() {
	tts.RegisterSynthesizer(NewSynthesizer, 0)
}

type Synthesizer struct {
	tts.BaseSynthesizer
	speechRate int
	synth      *nsSpeechSynthesizer
}

func NewSynthesizer() (tts.SpeechSynthesizer, error) {
	synth, err := newNsSpeechSynthesizer()
	if err != nil {
		return nil, err
	}
	rate, err := synth.getRate()
	if err != nil {
		return nil, err
	}
	synthesizer := &Synthesizer{
		speechRate: int(math.Round(rate)),
		synth:      synth,
	}
	return synthesizer, nil
}

func (s *Synthesizer) Name() string {
	return "NsSpeech"
}

func (s *Synthesizer) CreateNew() (tts.SpeechSynthesizer, error) {
	return NewSynthesizer()
}

func (s *Synthesizer) Release() {
	s.synth.free()
	s.synth = nil
}

func (s *Synthesizer) Stop() error {
	isSpeaking, err := s.IsSpeaking()
	if err != nil {
		return err
	}
	if !isSpeaking {
		return nil
	}
	return s.synth.stop()
}

func (s *Synthesizer) IsSpeaking() (bool, error) {
	return s.synth.isSpeaking()
}

func (s *Synthesizer) Speak(phrase *tts.Phrase) error {
	isSpeaking, err := s.IsSpeaking()
	if err != nil {
		return err
	}
	if isSpeaking {
		if err := s.Stop(); err != nil {
			return err
		}
	}
	if err := s.synth.speak(s.handlePhrase(phrase)); err != nil {
		return fmt.Errorf("Failed to speak: %v", err)
	}
	return nil
}

func (s *Synthesizer) handlePhrase(p *tts.Phrase) string {
	text := ""
	if p.Rate > 0 {
		text += fmt.Sprintf("[[rate %d]]", p.Rate)
	}
	if p.Silence > 0 {
		text += fmt.Sprintf("[[slnc %d]]", p.Silence)
	}
	text += p.Text
	return text
}

func (s *Synthesizer) SetSpeechRate(rate int) error {
	if err := s.synth.setRate(rate); err != nil {
		return err
	}
	newRate, err := s.synth.getRate()
	if err != nil {
		return err
	}
	s.speechRate = int(math.Round(newRate))
	return nil
}

func (s *Synthesizer) GetSpeechRate() int {
	return s.speechRate
}
