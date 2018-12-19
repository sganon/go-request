package request

import (
	"encoding"
	"errors"
	"net/http"
	"reflect"
)

// Differents payload issues
var (
	ErrInvalidPayload    = errors.New("invalid payload")
	ErrInvalidParameters = errors.New("invalid parameters")
	ErrUnexpected        = errors.New("unexpected error")
)

// Payload represents most basic payload of an `application/problem+json` response.
type Payload struct {
	Type   string `json:"type"`
	Title  string `json:"title"`
	Status int    `json:"status"`
}

// Validate provides validation and sets default values if needed
func (p *Payload) Validate() error {
	if p.Type == "" {
		p.Type = "about:blank"
	}
	if p.Title == "" || p.Status < http.StatusContinue {
		return ErrInvalidPayload
	}
	return nil
}

// InputProblem extends a standard problem payload with invalid parameters
type InputProblem struct {
	Payload
	InvalidParams []ParamError `json:"invalid_parameters"`
}

// Error implement error interface
func (i InputProblem) Error() string {
	return ErrInvalidParameters.Error()
}

// ParamError describe an error on a specific parameter
type ParamError struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

// UnexpectedProblem payload
type UnexpectedProblem Payload

// Error implements error interface
func (u UnexpectedProblem) Error() string {
	return ErrUnexpected.Error()
}

// StringSetter see flag.Value
type StringSetter interface {
	Set(string) error
}

// Types implementing TexTUnmarshaler and StringSetter used on unmarshalling
var (
	TextUnmarshalerType = reflect.TypeOf(new(encoding.TextUnmarshaler)).Elem()
	StringSetterType    = reflect.TypeOf(new(StringSetter)).Elem()
)
