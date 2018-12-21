package common

import (
	"errors"
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
