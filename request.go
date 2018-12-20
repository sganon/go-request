package request

import (
	"net/http"

	"github.com/sganon/go-request/common"
	"github.com/sganon/go-request/query"
)

func Decode(r *http.Request, inputQuery interface{}, inputBody interface{}) common.Problem {
	queryDecoder := query.NewDecoder(r)
	err := queryDecoder.Decode(inputQuery)
	if inputProblem, ok := err.(*common.InputProblem); ok && inputProblem != nil {
		return inputProblem
	}
	if unexpectedProblem, ok := err.(*common.UnexpectedProblem); ok && unexpectedProblem != nil {
		return unexpectedProblem
	}
	return nil
}
