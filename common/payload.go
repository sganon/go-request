package common

import (
	"encoding"
	"encoding/json"
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

type Problem interface {
	Send(http.ResponseWriter)
}

func baseSend(w http.ResponseWriter, status int, v interface{}) {
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)
	err := encoder.Encode(v)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "an unexpected error occured"}`))
	}
}

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

// InputProblem extends a standard problem payload with invalid parameters
type InputProblem struct {
	Payload
	InvalidParams []ParamError `json:"invalid_parameters"`
}

// Send implements Problem interface
func (i InputProblem) Send(w http.ResponseWriter) {
	baseSend(w, http.StatusBadRequest, i)
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

// Send implements Problem interface
func (u UnexpectedProblem) Send(w http.ResponseWriter) {
	baseSend(w, http.StatusInternalServerError, u)
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
