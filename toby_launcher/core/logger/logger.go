package logger

import (
	"io"
	"log"
	"os"
	"toby_launcher/apperrors"
)

type Logger interface {
	Printf(format string, v ...any)
	Error(err error)
	DebugPrintf(format string, v ...any)
	DebugError(err error)
	Release()
}

type StdLogger struct {
	logger       *log.Logger
	fileLogger   *log.Logger
	file         *os.File
	errorHandler apperrors.ErrorHandler
}

func NewStdLogger(output io.Writer, logFilePath string, errorHandler apperrors.ErrorHandler) (*StdLogger, error) {
	var file *os.File
	var fileLogger *log.Logger

	if logFilePath != "" {
		var err error
		file, err = os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		fileLogger = log.New(file, "", log.LstdFlags)
	} else {
		fileLogger = log.New(io.Discard, "", log.LstdFlags)
	}

	return &StdLogger{
		logger:       log.New(output, "", 0),
		fileLogger:   fileLogger,
		file:         file,
		errorHandler: errorHandler,
	}, nil
}

func (l *StdLogger) Printf(format string, v ...any) {
	l.logger.Printf(format, v...)
	l.fileLogger.Printf(format, v...)
}

func (l *StdLogger) Error(err error) {
	if err != nil {
		text := l.errorHandler.Handle(err)
		if text != "" {
			l.logger.Println(text)
			l.fileLogger.Println(text)
		}
	}
}

func (l *StdLogger) DebugPrintf(format string, v ...any) {
	l.fileLogger.Printf(format, v...)
}

func (l *StdLogger) DebugError(err error) {
	if err != nil {
		text := l.errorHandler.Handle(err)
		if text != "" {
			l.fileLogger.Println(text)
		}
	}
}

func (l *StdLogger) Release() {
	if l.file != nil {
		l.file.Close()
	}
}
