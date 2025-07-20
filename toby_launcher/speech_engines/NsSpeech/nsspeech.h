// nsspeech.h
#ifndef NSSPEECH_H
#define NSSPEECH_H

typedef void* NSSpeechSynthesizerPtr;

NSSpeechSynthesizerPtr NsSpeechInit();
int NsSpeechFree(NSSpeechSynthesizerPtr synth);
int NsSpeechSpeak(NSSpeechSynthesizerPtr synth, char *str);
int NsSpeechStop(NSSpeechSynthesizerPtr synth);
int NsSpeechSetRate(NSSpeechSynthesizerPtr synth, float rate);
float NsSpeechGetRate(NSSpeechSynthesizerPtr synth);
int NsSpeechIsSpeaking(NSSpeechSynthesizerPtr synth);

#endif