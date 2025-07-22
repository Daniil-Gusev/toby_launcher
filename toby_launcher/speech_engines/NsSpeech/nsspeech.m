//go:build darwin

#import <Foundation/Foundation.h>
#import <AppKit/AppKit.h>

typedef void* NSSpeechSynthesizerPtr;

NSSpeechSynthesizerPtr NsSpeechInit() {
    NSSpeechSynthesizer *synth = [[NSSpeechSynthesizer alloc] init];
    return synth;
}

int NsSpeechFree(NSSpeechSynthesizerPtr synth) {
    if (!synth) return -1;
    [(NSSpeechSynthesizer*)synth stopSpeaking];
    [(NSSpeechSynthesizer*)synth release];
    return 1;
}

int NsSpeechSpeak(NSSpeechSynthesizerPtr synth, char *str) {
    if (!synth) return -1;
    NSString *nsstr = [NSString stringWithCString:str encoding:NSUTF8StringEncoding];
    [(NSSpeechSynthesizer*)synth startSpeakingString:nsstr];
    return 1;
}

int NsSpeechStop(NSSpeechSynthesizerPtr synth) {
    if (!synth) return -1;
    [(NSSpeechSynthesizer*)synth stopSpeaking];
    return 1;
}

int NsSpeechSetRate(NSSpeechSynthesizerPtr synth, float rate) {
    if (!synth) return -1;
    if (rate <= 0) return -2;
    [(NSSpeechSynthesizer*)synth setRate:rate];
    return 1;
}

float NsSpeechGetRate(NSSpeechSynthesizerPtr synth) {
    if (!synth) return -1;
    return [(NSSpeechSynthesizer*)synth rate];
}

int NsSpeechIsSpeaking(NSSpeechSynthesizerPtr synth) {
    if (!synth) return -1;
    return [(NSSpeechSynthesizer*)synth isSpeaking] == YES ? 1 : 0;
}
