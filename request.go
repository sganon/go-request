package request

import (
	"encoding/json"
	"net/http"

	"github.com/sganon/go-request/problem"
	"github.com/sganon/go-request/query"
)

type Output interface {
	Validate() []problem.ParamError
}

func Decode(r *http.Request, queryOutput Output, bodyOutput Output) problem.Problem {
	var ok bool
	var inputProblem *problem.Input
	var unexpectedProblem *problem.UnexpectedProblem
	if queryOutput != nil {
		queryDecoder := query.NewDecoder(r)
		err := queryDecoder.Decode(queryOutput)
		if prob, ok := err.(*problem.Input); ok && prob != nil {
			inputProblem = prob
		}
		if unexpectedProblem, ok = err.(*problem.UnexpectedProblem); ok && unexpectedProblem != nil {
			return unexpectedProblem
		}
		queryErrors := queryOutput.Validate()
		if len(queryErrors) > 0 {
			if inputProblem == nil {
				prob := problem.DefaultInput
				inputProblem = &prob
			}
			inputProblem.InvalidParams = append(inputProblem.InvalidParams, queryErrors...)
		}
	}
	if bodyOutput != nil {
		bodyDecoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		err := bodyDecoder.Decode(&bodyOutput)
		if err != nil {
			return &problem.DefaultUnexpected
		}
		bodyErrors := bodyOutput.Validate()
		if len(bodyErrors) > 0 {
			if inputProblem == nil {
				prob := problem.DefaultInput
				inputProblem = &prob
			}
			inputProblem.InvalidParams = append(inputProblem.InvalidParams, bodyErrors...)
		}
	}
	if inputProblem == nil {
		return nil
	}
	return inputProblem
}
