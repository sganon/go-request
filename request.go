package request

import (
	"encoding/json"
	"net/http"

	"github.com/sganon/go-request/common"
	"github.com/sganon/go-request/query"
)

func Decode(r *http.Request, inputQuery interface{}, inputBody interface{}) common.Problem {
	var ok bool
	var inputProblem *common.InputProblem
	var unexpectedProblem *common.UnexpectedProblem
	if inputQuery != nil {
		queryDecoder := query.NewDecoder(r)
		err := queryDecoder.Decode(inputQuery)
		if inputProblem, ok = err.(*common.InputProblem); ok && inputProblem != nil {
			return inputProblem
		}
		if unexpectedProblem, ok = err.(*common.UnexpectedProblem); ok && unexpectedProblem != nil {
			return unexpectedProblem
		}
	}
	if inputBody != nil {
		bodyDecoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		err := bodyDecoder.Decode(&inputBody)
		if err != nil {
			return &common.UnexpectedProblem{
				Title:  "Unexpected error",
				Type:   "about:blank",
				Status: http.StatusInternalServerError,
			}
		}
	}
	return nil
}
