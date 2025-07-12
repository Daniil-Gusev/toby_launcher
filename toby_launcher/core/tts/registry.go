package tts

import (
	"sort"
	"toby_launcher/core/logger"
)

type factoryFunc func() (SpeechSynthesizer, error)

type factory struct {
	create   factoryFunc
	priority int
}

var synthesizerFactories []factory

func RegisterSynthesizer(f factoryFunc, priority int) {
	fact := factory{
		create:   f,
		priority: priority,
	}
	synthesizerFactories = append(synthesizerFactories, fact)
}

func GetAvailableSynthesizers(logger logger.Logger) []SpeechSynthesizer {
	sort.Slice(synthesizerFactories, func(i, j int) bool {
		return synthesizerFactories[i].priority < synthesizerFactories[j].priority
	})
	syns := make([]SpeechSynthesizer, 0, len(synthesizerFactories))
	for _, factory := range synthesizerFactories {
		s, err := factory.create()
		if err == nil {
			syns = append(syns, s)
		} else {
			logger.DebugPrintf("Failed to initialize factory speech synthesizer: %v\r\n", err)
		}
	}
	for _, s := range syns {
		s.Release()
	}
	return syns
}
