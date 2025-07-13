package espeak

import (
	"fmt"
	"os/exec"
	"time"
	"toby_launcher/core/tts"
)

func init() {
	tts.RegisterSynthesizer(NewSynthesizer, 1)
}

type Synthesizer struct {
	tts.BaseSynthesizer
	speechRate int
	cmd        *exec.Cmd
	cmdPath    string
	isSpeaking bool
}

func NewSynthesizer() (tts.SpeechSynthesizer, error) {
	var cmdPath string
	if path, err := exec.LookPath(espeakNgExecutableName); err == nil {
		cmdPath = path
	} else {
		path, err := exec.LookPath(espeakExecutableName)
		if err != nil {
			return nil, err
		}
		cmdPath = path
	}
	synth := &Synthesizer{
		cmdPath:    cmdPath,
		speechRate: 180,
		isSpeaking: false,
	}
	return synth, nil
}

func (s *Synthesizer) Name() string {
	return "espeak"
}

func (s *Synthesizer) CreateNew() (tts.SpeechSynthesizer, error) {
	return NewSynthesizer()
}

func (s *Synthesizer) Release() {
	if err := s.Stop(); err != nil {
		s.LogError(err)
	}
}

func (s *Synthesizer) Stop() error {
	if s.cmd == nil || s.cmd.Process == nil {
		return fmt.Errorf("espeak process is not initialized")
	}
	isSpeaking, err := s.IsSpeaking()
	if err != nil {
		return err
	}
	if isSpeaking {
		if err := s.cmd.Process.Kill(); err != nil {
			return err
		}
		s.isSpeaking = false
		s.cmd = nil
	}
	return nil
}

func (s *Synthesizer) IsSpeaking() (bool, error) {
	if s.cmd != nil && s.cmd.Process == nil {
		return false, fmt.Errorf("espeak process is not initialized")
	}
	return s.isSpeaking, nil
}

func (s *Synthesizer) Speak(phrase *tts.Phrase) error {
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
	isSpeaking, err := s.IsSpeaking()
	if err != nil {
		return err
	}
	if isSpeaking {
		if err := s.Stop(); err != nil {
			return err
		}
	}
	args := make([]string, 0, 2)
	rate := phrase.Rate
	if rate == 0 {
		rate = s.speechRate
	}
	if rate > 0 {
		args = append(args, fmt.Sprintf("-s%d", s.speechRate))
	}
	args = append(args, phrase.Text)
	s.cmd = exec.Command(s.cmdPath, args...)
	if err := s.cmd.Start(); err != nil {
		return err
	}
	s.isSpeaking = true
	go func() {
		if err := s.cmd.Wait(); err != nil {
			if s.isSpeaking {
				s.LogError(err)
			}
		}
		s.isSpeaking = false
	}()
	return nil
}

func (s *Synthesizer) SetSpeechRate(rate int) error {
	s.speechRate = rate
	return nil
}

func (s *Synthesizer) GetSpeechRate() int {
	return s.speechRate
}
