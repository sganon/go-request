package common_test

import (
	"net/http/httptest"
	"testing"

	"github.com/sganon/go-request/common"
	"github.com/stretchr/testify/assert"
)

func TestInputProblemError(t *testing.T) {
	prob := common.InputProblem{}
	assert.Equal(t, common.ErrInvalidParameters.Error(), prob.Error())
}

func TestInputUnexpectedError(t *testing.T) {
	prob := common.UnexpectedProblem{}
	assert.Equal(t, common.ErrUnexpected.Error(), prob.Error())
}

func TestInputProblemSend(t *testing.T) {
	w := httptest.NewRecorder()
	prob := common.InputProblem{
		Payload: common.Payload{
			Title:  "Test problem",
			Status: 400,
		},
	}
	prob.Send(w)
}

func TestUnexpectedProblemSend(t *testing.T) {
	w := httptest.NewRecorder()
	prob := common.UnexpectedProblem{
		Title:  "Test problem",
		Status: 500,
	}
	prob.Send(w)
}
