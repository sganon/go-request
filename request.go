package request

import (
	"encoding/json"
	"net/http"

	"github.com/sganon/go-request/form"
	"github.com/sganon/go-request/problem"
)

type Output interface {
	Validate() []problem.ParamError
}

func Decode(r *http.Request, formOutput interface{}, bodyOutput interface{}) problem.Problem {
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
	}
	if bodyOutput != nil {
		bodyDecoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		err := bodyDecoder.Decode(&bodyOutput)
		if err != nil {
			return &problem.DefaultUnexpected
		}
	}
	if inputProblem == nil {
		return nil
	}
	return inputProblem
}

func DecodeAndValidate(r *http.Request, formOutput Output, bodyOutput Output) problem.Problem {
	var inputProblem *problem.Input
	p := Decode(r, formOutput, bodyOutput)
	if prob, ok := p.(*problem.Input); ok {
		inputProblem = prob
	} else {
		return prob
	}
	if formOutput != nil {
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
