package common_test

import (
	"net/http"
	"testing"

	"github.com/sganon/go-request/common"
	"github.com/stretchr/testify/assert"
)

func TestPayloadValidate(t *testing.T) {
	payload := common.Payload{
		Title:  "Test Payload",
		Status: http.StatusInternalServerError,
	}
	err := payload.Validate()
	assert.NoError(t, err)
	assert.Equal(t, "about:blank", payload.Type)

	payload.Title = ""
	err = payload.Validate()
	assert.Equal(t, common.ErrInvalidPayload, err)
	err = nil

	payload.Title = "Test Payload"
	payload.Status = http.StatusOK
	err = payload.Validate()
	assert.Equal(t, common.ErrInvalidPayload, err)
}

func TestInputProblemError(t *testing.T) {
	prob := common.InputProblem{}
	assert.Equal(t, common.ErrInvalidParameters.Error(), prob.Error())
}

func TestInputUnexpectedError(t *testing.T) {
	prob := common.UnexpectedProblem{}
	assert.Equal(t, common.ErrUnexpected.Error(), prob.Error())
}
