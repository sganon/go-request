package common

import (
	"encoding/json"
	"net/http"
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

// InputProblem extends a standard problem payload with invalid parameters
type InputProblem struct {
	*Payload
	InvalidParams []ParamError `json:"invalid_parameters"`
}

// Send implements Problem interface
func (i InputProblem) Send(w http.ResponseWriter) {
	err := i.Validate()
	if err != nil {
		panic(err)
	}
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
type UnexpectedProblem struct {
	*Payload
}

// Error implements error interface
func (u UnexpectedProblem) Error() string {
	return ErrUnexpected.Error()
}

// Send implements Problem interface
func (u UnexpectedProblem) Send(w http.ResponseWriter) {
	err := u.Validate()
	if err != nil {
		panic(err)
	}
	baseSend(w, http.StatusInternalServerError, u)
}

var DefaultUnexpectedProblem = UnexpectedProblem{
	Payload: &Payload{
		Type:   "about:blank",
		Title:  "An unexpected error occured decoding request",
		Status: http.StatusInternalServerError,
	},
}

var DefaultInputProblem = InputProblem{
	Payload: &Payload{
		Type:   "about:blank",
		Title:  "Your parameters didn't validate",
		Status: http.StatusBadRequest,
	},
}
