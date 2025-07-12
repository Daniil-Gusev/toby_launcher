//go:build windows

package sapi

import (
	"fmt"
	"toby_launcher/core/tts"
)

func init() {
	tts.RegisterSynthesizer(NewSynthesizer, 0)
}

type Synthesizer struct {
	tts.BaseSynthesizer
	speechRate int
	synth      *sapiSynthesizer
	isSpeaking bool
}

func NewSynthesizer() (tts.SpeechSynthesizer, error) {
	synth, err := newSapiSynthesizer()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize SAPI synthesizer: %v", err)
	}
	rate, err := synth.getSpeechRate()
	if err != nil {
		synth.release()
		return nil, fmt.Errorf("failed to get speech rate: %v", err)
	}
	synthesizer := &Synthesizer{
		synth:      synth,
		isSpeaking: false,
	}
	synthesizer.speechRate = synthesizer.sapiRateToRate(rate)
	return synthesizer, nil
}

func (s *Synthesizer) Name() string {
	return "sapi (unstable)"
}

func (s *Synthesizer) CreateNew() (tts.SpeechSynthesizer, error) {
	return NewSynthesizer()
}

func (s *Synthesizer) Release() {
	if s.synth != nil {
		s.synth.release()
		s.synth = nil
	}
}

func (s *Synthesizer) Stop() error {
	if s.synth == nil {
		return fmt.Errorf("SAPI synthesizer is not initialized")
	}
	isSpeaking, err := s.IsSpeaking()
	if err != nil {
		return err
	}
	if !isSpeaking {
		return nil
	}
	if err := s.synth.stop(); err != nil {
		return err
	}
	s.isSpeaking = false
	return nil
}

func (s *Synthesizer) IsSpeaking() (bool, error) {
	if s.synth == nil {
		return false, fmt.Errorf("SAPI synthesizer is not initialized")
	}
	return s.isSpeaking, nil
}

func (s *Synthesizer) Speak(phrase *tts.Phrase) error {
	if s.synth == nil {
		return fmt.Errorf("SAPI synthesizer is not initialized")
	}
	if phrase.Text == "" {
		return fmt.Errorf("no text to speak has been specified")
	}
	isSpeaking, err := s.IsSpeaking()
	if err != nil {
		return err
	}
	if isSpeaking {
		if err := s.Stop(); err != nil {
			return err
		}
	}
	ssml := s.handlePhrase(phrase)
	if err := s.synth.speak(ssml); err != nil {
		return err
	}
	s.isSpeaking = true
	go func() {
		if err := s.synth.wait(); err != nil {
			s.LogError(err)
		}
		s.isSpeaking = false
	}()
	return nil
}

func (s *Synthesizer) handlePhrase(p *tts.Phrase) string {
	rate := p.Rate
	if rate == 0 {
		rate = s.speechRate
	}
	prosodyRate := s.rateToSapiRate(rate)
	ssml := `<?xml version="1.0" encoding="UTF-8"?><speak version="1.0" xmlns="http://www.w3.org/2001/10/synthesis" xml:lang="en-US">`
	if p.Silence > 0 {
		ssml += fmt.Sprintf(`<break time="%dms"/>`, p.Silence)
	}
	if prosodyRate != 0 {
		ssml += fmt.Sprintf(`<prosody rate="%s">%s</prosody>`, string(prosodyRate), p.Text)
	} else {
		ssml += p.Text
	}
	ssml += `</speak>`
	return ssml
}

func (s *Synthesizer) rateToSapiRate(rate int) int {
	if rate <= 0 {
		return 0
	}
	// Map rate (0-600) to SAPI's -10 to +10
	sapiRate := (float64(rate) - 250) / 25
	if sapiRate < -10 {
		sapiRate = -10
	} else if sapiRate > 10 {
		sapiRate = 10
	}
	return int(sapiRate)
}

func (s *Synthesizer) sapiRateToRate(rate int) int {
	return (rate * 25) + 250
}

func (s *Synthesizer) SetSpeechRate(rate int) error {
	if s.synth == nil {
		return fmt.Errorf("SAPI synthesizer is not initialized")
	}
	sapiRate := s.rateToSapiRate(rate)
	if err := s.synth.setSpeechRate(sapiRate); err != nil {
		return err
	}
	newRate, err := s.synth.getSpeechRate()
	if err != nil {
		return err
	}
	s.speechRate = s.sapiRateToRate(newRate)
	return nil
}

func (s *Synthesizer) GetSpeechRate() int {
	return s.speechRate
}
