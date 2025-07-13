//go:build darwin

package NsSpeech

/*
#cgo CFLAGS: -x objective-c -Wno-deprecated-declarations

#cgo LDFLAGS: -framework AppKit
#import <Foundation/Foundation.h>
#import <AppKit/AppKit.h>

typedef void* NSSpeechSynthesizerPtr;

NSSpeechSynthesizerPtr NsSpeechInit(){
NSSpeechSynthesizer *synth=[[NSSpeechSynthesizer alloc] init];
return synth;
}

int NsSpeechFree(NSSpeechSynthesizerPtr synth){
if(!synth) return -1;
[(NSSpeechSynthesizer*)synth stopSpeaking];
[(NSSpeechSynthesizer*)synth release];
return 1;
}

int NsSpeechSpeak(NSSpeechSynthesizerPtr synth,char *str){
if(!synth) return -1;
NSString *nsstr = [NSString stringWithCString: str encoding:NSUTF8StringEncoding];
[(NSSpeechSynthesizer*)synth startSpeakingString:nsstr];
return 1;
}

int NsSpeechStop(NSSpeechSynthesizerPtr synth){
if(!synth) return -1;
[(NSSpeechSynthesizer*)synth stopSpeaking];
    return 1;
}

int NsSpeechSetRate(NSSpeechSynthesizerPtr synth, float rate){
if(!synth) return -1;
if(rate<=0) return -2;
[(NSSpeechSynthesizer*)synth setRate:rate];
return 1;
}

float NsSpeechGetRate(NSSpeechSynthesizerPtr synth){
if(!synth) return -1;
return [(NSSpeechSynthesizer*)synth rate];
}

int NsSpeechIsSpeaking(NSSpeechSynthesizerPtr synth){
if(!synth) return -1;
return [(NSSpeechSynthesizer*)synth isSpeaking] == YES ? 1 : 0;
}


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
