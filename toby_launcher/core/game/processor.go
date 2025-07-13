package game

import (
	"bufio"
	"bytes"
	"regexp"
	"toby_launcher/apperrors"
	"toby_launcher/config"
	"toby_launcher/core/logger"
	"toby_launcher/core/tts"
	"toby_launcher/utils/file_utils"
)

type SubstitutionData struct {
	Pattern     string `json:"pattern"`
	Replacement string `json:"replacement"`
}

type TextRulesData struct {
	Separator     string             `json:"separator"`
	Exclusions    []string           `json:"exclusions"`
	Substitutions []SubstitutionData `json:"substitutions"`
}

type Substitution struct {
	pattern     *regexp.Regexp
	replacement string
}

// TextProcessor processes game output based on rules from tts_lines.json.
type TextProcessor struct {
	logger          logger.Logger
	tts             *tts.TtsManager
	config          *config.Config
	separator       *regexp.Regexp
	exclusions      []*regexp.Regexp
	substitutions   []Substitution
	startProcessing bool
}

// NewTextProcessor creates a new TextProcessor instance.
func NewTextProcessor(cfg *config.Config, logger logger.Logger, tts *tts.TtsManager) *TextProcessor {
	processor := &TextProcessor{
		logger:          logger,
		tts:             tts,
		config:          cfg,
		exclusions:      make([]*regexp.Regexp, 0, 20),
		substitutions:   make([]Substitution, 0, 20),
		startProcessing: false,
	}
	if err := processor.loadRules(); err != nil {
		logger.Error(err)
	}
	return processor
}

func (p *TextProcessor) loadRules() error {
	path := p.config.Paths.TextRulesPath()
	var rules TextRulesData
	if err := file_utils.LoadData(path, &rules); err != nil {
		return apperrors.New(apperrors.Err, "Failed to load text handling rules in file $file: $error", map[string]any{"error": err, "file": path})
	}
	re, err := regexp.Compile(rules.Separator)
	if err != nil {
		p.logger.Error(apperrors.New(apperrors.Err, "Invalid separate regex pattern \"$pattern\" in file $file: $error", map[string]any{
			"pattern": rules.Separator,
			"file":    path,
			"error":   err,
		}))
	}
	p.separator = re
	for _, pattern := range rules.Exclusions {
		re, err := regexp.Compile(pattern)
		if err != nil {
			p.logger.Error(apperrors.New(apperrors.Err, "Invalid exclude regex pattern \"$pattern\" in file $file: $error", map[string]any{
				"pattern": pattern,
				"file":    path,
				"error":   err,
			}))
			continue
		}
		p.exclusions = append(p.exclusions, re)
	}
	for _, rule := range rules.Substitutions {
		re, err := regexp.Compile(rule.Pattern)
		if err != nil {
			p.logger.Error(apperrors.New(apperrors.Err, "Invalid substitute regex pattern \"$pattern\" in file $file: $error", map[string]any{
				"pattern": rule.Pattern,
				"file":    path,
				"error":   err,
			}))
			continue
		}
		subst := Substitution{
			pattern:     re,
			replacement: rule.Replacement,
		}
		p.substitutions = append(p.substitutions, subst)
	}
	return nil
}

func (p *TextProcessor) Write(data []byte) (int, error) {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()
		if p.config.Gzdoom.DebugOutput {
			p.logger.Printf(line)
		}
		if p.separator != nil && p.separator.MatchString(line) {
			p.startProcessing = true
			continue
		}
		if !p.startProcessing {
			continue
		}
		for _, re := range p.exclusions {
			if re.MatchString(line) {
				continue
			}
		}
		processedLine := line
		for _, rule := range p.substitutions {
			processedLine = rule.pattern.ReplaceAllString(processedLine, rule.replacement)
		}
		if processedLine != "" {
			p.tts.Speak(processedLine)
			p.logger.DebugPrintf("speaking: %s\r\n", processedLine)
		}
	}
	if err := scanner.Err(); err != nil {
		p.logger.Error(apperrors.New(apperrors.Err, "Error scanning output: $error", map[string]any{"error": err}))
	}
	return len(data), nil
}
