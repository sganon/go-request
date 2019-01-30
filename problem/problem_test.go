package problem_test

import (
	"net/http/httptest"
	"testing"

	"github.com/sganon/go-request/problem"
	"github.com/stretchr/testify/assert"
)

func TestInputProblemError(t *testing.T) {
	prob := problem.InputProblem{}
	assert.Equal(t, problem.ErrInvalidParameters.Error(), prob.Error())
}

func TestInputUnexpectedError(t *testing.T) {
	prob := problem.UnexpectedProblem{}
	assert.Equal(t, problem.ErrUnexpected.Error(), prob.Error())
}

func TestInputProblemSend(t *testing.T) {
	w := httptest.NewRecorder()
	prob := problem.InputProblem{
		Payload: &problem.Payload{
			Title:  "Test problem",
			Status: 400,
		},
	}
	prob.Send(w)

	prob.Status = 0
	assert.Panics(t, func() {
		prob.Send(w)
	}, "send should panic if problem is not valid")
}

func TestUnexpectedProblemSend(t *testing.T) {
	w := httptest.NewRecorder()
	prob := problem.UnexpectedProblem{
		Payload: &problem.Payload{
			Title:  "Test problem",
			Status: 500,
		},
	}
	prob.Send(w)

	prob.Status = 0
	assert.Panics(t, func() {
		prob.Send(w)
	}, "send should panic if problem is not valid")
}
