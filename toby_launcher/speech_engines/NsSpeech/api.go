//go:build darwin

package NsSpeech

/*
#cgo CFLAGS: -x objective-c -Wno-deprecated-declarations

#cgo LDFLAGS: -framework AppKit
#include "nsspeech.h"

*/
import "C"

import (
	"fmt"
	"unsafe"
)

type nsSpeechSynthesizer struct {
	ptr unsafe.Pointer
}

func newNsSpeechSynthesizer() (*nsSpeechSynthesizer, error) {
	ptr := C.NsSpeechInit()
	if ptr == nil {
		return nil, fmt.Errorf("cannot initialize NsSpeechSynthesizer interface")
	}
	return &nsSpeechSynthesizer{ptr: unsafe.Pointer(ptr)}, nil
}

func (s *nsSpeechSynthesizer) speak(text string) error {
	ret := C.NsSpeechSpeak(C.NSSpeechSynthesizerPtr(s.ptr), C.CString(text))
	if ret == -1 {
		return fmt.Errorf("nsSpeechSynthesizer interface is not initialized")
	}
	if text == "" {
		return fmt.Errorf("no text to speak has been specified")
	}
	return nil
}

func (s *nsSpeechSynthesizer) stop() error {
	ret := C.NsSpeechStop(C.NSSpeechSynthesizerPtr(s.ptr))
	if ret == -1 {
		return fmt.Errorf("nsSpeechSynthesizer interface is not initialized")
	}
	return nil
}

func (s *nsSpeechSynthesizer) setRate(rate int) error {
	ret := C.NsSpeechSetRate(C.NSSpeechSynthesizerPtr(s.ptr), C.float(rate))
	if ret == -1 {
		return fmt.Errorf("nsSpeechSynthesizer interface is not initialized")
	}
	if ret == -2 {
		return fmt.Errorf("rate value is out of range")
	}
	return nil
}

func (s *nsSpeechSynthesizer) getRate() (float64, error) {
	ret := C.NsSpeechGetRate(C.NSSpeechSynthesizerPtr(s.ptr))
	if ret == -1 {
		return 0, fmt.Errorf("nsSpeechSynthesizer interface is not initialized")
	}
	return float64(ret), nil
}

func (s *nsSpeechSynthesizer) isSpeaking() (bool, error) {
	switch C.NsSpeechIsSpeaking(C.NSSpeechSynthesizerPtr(s.ptr)) {
	case -1:
		return false, fmt.Errorf("nsSpeechSynthesizer interface is not initialized")
	case 0:
		return false, nil
	case 1:
		return true, nil
	}
	return false, fmt.Errorf("unrecognized error")
}

func (s *nsSpeechSynthesizer) free() {
	C.NsSpeechFree(C.NSSpeechSynthesizerPtr(s.ptr))
}
