package request_test

import (
	"net/http"
	"testing"

	request "github.com/sganon/go-request"
	"github.com/stretchr/testify/assert"
)

func TestPayloadValidate(t *testing.T) {
	payload := request.Payload{
		Title:  "Test Payload",
		Status: http.StatusInternalServerError,
	}
	err := payload.Validate()
	assert.NoError(t, err)
	assert.Equal(t, "about:blank", payload.Type)

	payload.Title = ""
	err = payload.Validate()
	assert.Equal(t, request.ErrInvalidPayload, err)
	err = nil

	payload.Title = "Test Payload"
	payload.Status = http.StatusOK
	err = payload.Validate()
	assert.Equal(t, request.ErrInvalidPayload, err)
}

func TestInputProblemError(t *testing.T) {
	prob := request.InputProblem{}
	assert.Equal(t, request.ErrInvalidParameters.Error(), prob.Error())
}

func TestInputUnexpectedError(t *testing.T) {
	prob := request.UnexpectedProblem{}
	assert.Equal(t, request.ErrUnexpected.Error(), prob.Error())
}
