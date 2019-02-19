package request

import (
	"encoding/json"
	"net/http"

	"github.com/sganon/go-request/problem"
	"github.com/sganon/go-request/form"
)

type Output interface {
	Validate() []problem.ParamError
}

func Decode(r *http.Request, formOutput Output, bodyOutput Output) problem.Problem {
	var ok bool
	var inputProblem *problem.Input
	var unexpectedProblem *problem.UnexpectedProblem
	if formOutput != nil {
		formDecoder := form.NewDecoder(r)
		err := formDecoder.Decode(formOutput)
		if prob, ok := err.(*problem.Input); ok && prob != nil {
			inputProblem = prob
		}
		if unexpectedProblem, ok = err.(*problem.UnexpectedProblem); ok && unexpectedProblem != nil {
			return unexpectedProblem
		}
		formErrors := formOutput.Validate()
		if len(formErrors) > 0 {
			if inputProblem == nil {
				prob := problem.DefaultInput
				inputProblem = &prob
			}
			inputProblem.InvalidParams = append(inputProblem.InvalidParams, formErrors...)
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
