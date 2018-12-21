package common

import (
	"encoding"
	"errors"
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
	if p.Title == "" || p.Status < 400 {
		return ErrInvalidPayload
	}
	return nil
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
