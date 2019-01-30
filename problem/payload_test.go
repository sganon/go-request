package problem_test

import (
	"net/http"
	"testing"

	"github.com/sganon/go-request/problem"
	"github.com/stretchr/testify/assert"
)

func TestPayloadValidate(t *testing.T) {
	payload := problem.Payload{
		Title:  "Test Payload",
		Status: http.StatusInternalServerError,
	}
	err := payload.Validate()
	assert.NoError(t, err)
	assert.Equal(t, "about:blank", payload.Type)

	payload.Title = ""
	err = payload.Validate()
	assert.Equal(t, problem.ErrInvalidPayload, err)
	err = nil

	payload.Title = "Test Payload"
	payload.Status = http.StatusOK
	err = payload.Validate()
	assert.Equal(t, problem.ErrInvalidPayload, err)
}
