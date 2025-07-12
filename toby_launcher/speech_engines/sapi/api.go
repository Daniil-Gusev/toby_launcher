//go:build windows

package sapi

import (
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

type sapiSynthesizer struct {
	voice *ole.IDispatch
}

func newSapiSynthesizer() (*sapiSynthesizer, error) {
	err := ole.CoInitialize(0)
	if err != nil {
		return nil, err
	}

	unknown, err := oleutil.CreateObject("SAPI.SpVoice")
	if err != nil {
		ole.CoUninitialize()
		return nil, err
	}

	voice, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil || voice == nil {
		unknown.Release()
		ole.CoUninitialize()
		return nil, ole.NewError(ole.E_NOINTERFACE)
	}

	return &sapiSynthesizer{voice: voice}, nil
}

func (s *sapiSynthesizer) release() {
	if s.voice != nil {
		s.voice.Release()
		s.voice = nil
	}
	ole.CoUninitialize()
}

func (s *sapiSynthesizer) speak(text string) error {
	_, err := oleutil.CallMethod(s.voice, "Speak", text, 9)
	return err
}

func (s *sapiSynthesizer) stop() error {
	_, err := oleutil.CallMethod(s.voice, "Pause")
	return err
}

func (s *sapiSynthesizer) wait() error {
	_, err := oleutil.CallMethod(s.voice, "WaitUntilDone", -1)
	return err
}

func (s *sapiSynthesizer) isSpeaking() (bool, error) {
	status, err := oleutil.GetProperty(s.voice, "Status")
	if err != nil {
		return false, err
	}
	defer status.Clear()

	runningState, err := oleutil.GetProperty(status.ToIDispatch(), "RunningState")
	if err != nil {
		return false, err
	}
	return runningState.Value().(int32) == 2, nil
}

func (s *sapiSynthesizer) setSpeechRate(rate int) error {
	_, err := oleutil.PutProperty(s.voice, "Rate", rate)
	return err
}

func (s *sapiSynthesizer) getSpeechRate() (int, error) {
	rate, err := oleutil.GetProperty(s.voice, "Rate")
	if err != nil {
		return 0, err
	}
	return int(rate.Value().(int32)), nil
}
