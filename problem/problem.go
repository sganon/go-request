package problem

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

// Input extends a standard problem problem with invalid parameters
type Input struct {
	*Payload
	InvalidParams []ParamError `json:"invalid_parameters"`
}

// Send implements Problem interface
func (i Input) Send(w http.ResponseWriter) {
	err := i.Validate()
	if err != nil {
		panic(err)
	}
	baseSend(w, http.StatusBadRequest, i)
}

// Error implement error interface
func (i Input) Error() string {
	return ErrInvalidParameters.Error()
}

// ParamError describe an error on a specific parameter
type ParamError struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

// UnexpectedProblem problem
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

// ForbiddenProblem problem
type ForbiddenProblem struct {
	*Payload
}

// Error implements error interface
func (f ForbiddenProblem) Error() string {
	return ErrForbidden.Error()
}

// Send implements Problem interface
func (f ForbiddenProblem) Send(w http.ResponseWriter) {
	err := f.Validate()
	if err != nil {
		panic(err)
	}
	baseSend(w, http.StatusForbidden, f)
}

var DefaultUnexpected = UnexpectedProblem{
	Payload: &Payload{
		Type:   "about:blank",
		Title:  "An unexpected error occured decoding request",
		Status: http.StatusInternalServerError,
	},
}

var DefaultInput = Input{
	Payload: &Payload{
		Type:   "about:blank",
		Title:  "Your parameters didn't validate",
		Status: http.StatusBadRequest,
	},
}

var DefaultForbidden = Input{
	Payload: &Payload{
		Type:   "about:blank",
		Title:  "Your are missing credentials or have insufficient rights",
		Status: http.StatusForbidden,
	},
}
