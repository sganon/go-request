package problem

import (
	"errors"
	"net/http"
)

// Differents problem issues
var (
	ErrInvalidPayload    = errors.New("invalid problem")
	ErrInvalidParameters = errors.New("invalid parameters")
	ErrUnexpected        = errors.New("unexpected error")
	ErrForbidden         = errors.New("forbidden")
)

// Payload represents most basic problem of an `application/problem+json` response.
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

// Send implements Problem
func (p Payload) Send(w http.ResponseWriter) {
	baseSend(w, p.Status, p)
}
