package apperrors

import (
	"fmt"
	"strings"
	"toby_launcher/utils"
)

type ErrorCode string

const (
	Err           ErrorCode = "ERROR"
	ErrUnknown    ErrorCode = "UNKNOWN_ERROR"
	ErrEOF        ErrorCode = "END_OF_INPUT"
	ErrStateStack           = "STATE_STACK_ERROR"
	ErrInternal   ErrorCode = "INTERNAL_ERROR"
	ErrSpeech     ErrorCode = "SPEECH_ERROR"
)

type AppError struct {
	Code    ErrorCode
	Message string
	Details map[string]any
}

func (e *AppError) Error() string {
	msg := e.Message
	if len(e.Details) == 0 {
		return msg
	}
	msg += "\r\nDetails:"
	for key, value := range e.Details {
		msg += fmt.Sprintf("\r\n%s: %v", key, value)
	}
	return msg
}

func New(code ErrorCode, message string, details map[string]any) *AppError {
	return &AppError{Code: code, Message: message, Details: details}
}

type AppErrors struct {
	errors []error
}

func NewErrors(errs []error) *AppErrors {
	appErrors := &AppErrors{
		errors: make([]error, 0, 10),
	}
	for _, err := range errs {
		if err != nil {
			appErrors.Add(err)
		}
	}
	return appErrors
}

func (e *AppErrors) Add(err error) {
	e.errors = append(e.errors, err)
}

func (e *AppErrors) Error() string {
	var text string
	for _, err := range e.errors {
		text += (err.Error() + "\r\n")
	}
	return text
}

func (e *AppErrors) Count() int {
	return len(e.errors)
}

func (e *AppErrors) Errors() []error {
	return e.errors
}

type ErrorHandler interface {
	Handle(error) string
}

type StdErrorHandler struct{}

func (h *StdErrorHandler) Handle(err error) string {
	if err == nil {
		return ""
	}
	if appErrs, ok := err.(*AppErrors); ok {
		if appErrs.Count() == 0 {
			return ""
		}
		var text string
		for i, e := range appErrs.Errors() {
			text += (h.Handle(e))
			if i != (appErrs.Count() - 1) {
				text += "\r\n"
			}
		}
		return text
	}
	if appErr, ok := err.(*AppError); ok {
		errMsg := appErr.Message
		var information string
		switch appErr.Code {
		case ErrInternal:
			information = "Internal error"
		case ErrSpeech:
			information = "Speech error"
		case "ErrUnknown":
			information = "Error:"
		}
		for key, value := range appErr.Details {
			if appErrValue, ok := value.(*AppError); ok {
				errMsg = strings.ReplaceAll(errMsg, "$"+key, h.Handle(appErrValue))
			}
			if appErrsValue, ok := value.(*AppErrors); ok {
				errMsg = strings.ReplaceAll(errMsg, "$"+key, h.Handle(appErrsValue))
			}
		}
		errMsg = utils.SubstituteParams(errMsg, appErr.Details)
		if information != "" {
			return fmt.Sprintf("%s: %s", information, errMsg)
		}
		return errMsg
	}
	appErr := New(ErrUnknown, fmt.Sprintf("%v", err), nil)
	return h.Handle(appErr)
}
