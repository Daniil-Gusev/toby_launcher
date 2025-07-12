package core

import (
	"fmt"
	"github.com/chzyer/readline"
	"io"
	"os"
	"strings"
	"toby_launcher/apperrors"
)

type Console interface {
	Read() (string, error)
	Write(string) error
	Close() error
}

type ReadlineConsole struct {
	rl *readline.Instance
}

func NewReadlineConsole() (*ReadlineConsole, error) {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:       "> ",
		HistoryLimit: 200,
	})
	if err != nil {
		return nil, err
	}
	return &ReadlineConsole{rl: rl}, nil
}

func (c *ReadlineConsole) Read() (string, error) {
	line, err := c.rl.Readline()
	if err == readline.ErrInterrupt {
		return "", apperrors.New(apperrors.ErrEOF, "interrupt", nil)
	}
	if err == io.EOF {
		return "", apperrors.New(apperrors.ErrEOF, "interrupt", nil)
	}
	if err != nil {
		return "", apperrors.New(apperrors.ErrInternal, "Read error: $error", map[string]any{
			"error": fmt.Sprintf("%v", err),
		})
	}
	return strings.TrimSpace(line), nil
}

func (c *ReadlineConsole) Write(s string) error {
	if _, err := os.Stdout.WriteString(s); err != nil {
		return err
	}
	return nil
}

func (c *ReadlineConsole) Close() error {
	return c.rl.Close()
}
