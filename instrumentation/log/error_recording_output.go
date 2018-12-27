package log

import (
	"fmt"
	"regexp"
	"sync"
)

type ErrorRecordingOutput struct {
	allowedErrors []string
	errorRecorder *errorRecorder
}

type errorRecorder struct {
	sync.Mutex
	unexpectedErrors []*unexpectedError
}

type unexpectedError struct {
	message string
	err     error
}

func (e *unexpectedError) String() string {
	return fmt.Sprintf("%s (passed Error object: %v)", e.message, e.err)
}

func NewErrorRecordingOutput(allowedErrors []string) *ErrorRecordingOutput {
	return &ErrorRecordingOutput{allowedErrors: allowedErrors, errorRecorder: &errorRecorder{}}
}

func (o *ErrorRecordingOutput) Append(level string, message string, params ...*Field) {
	if level == "error" {
		for _, allowedMessage := range o.allowedErrors {
			if matched, _ := regexp.MatchString(allowedMessage, message); matched {
				return
			}
			for _, f := range params {
				if f.Key == "error" {
					if matched, _ := regexp.MatchString(allowedMessage, f.String()); matched {
						return
					}
				}
			}
		}
		o.recordUnexpectedError(message, params)
	}
}

func (o *ErrorRecordingOutput) recordUnexpectedError(message string, fields []*Field) {
	o.errorRecorder.Lock()
	defer o.errorRecorder.Unlock()

	var err error
	for _, f := range fields {
		if f.Error != nil {
			err = f.Error
		}
	}

	o.errorRecorder.unexpectedErrors = append(o.errorRecorder.unexpectedErrors, &unexpectedError{message: message, err: err})

}

func (o *ErrorRecordingOutput) GetUnexpectedErrors() (errors []string) {
	o.errorRecorder.Lock()
	defer o.errorRecorder.Unlock()

	for _, e := range o.errorRecorder.unexpectedErrors {
		errors = append(errors, e.String())
	}

	return
}

func (o *ErrorRecordingOutput) HasErrors() bool {
	o.errorRecorder.Lock()
	defer o.errorRecorder.Unlock()
	return len(o.errorRecorder.unexpectedErrors) > 0
}
