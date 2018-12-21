package request

import (
	"encoding/json"
	"net/http"

	"github.com/sganon/go-request/common"
	"github.com/sganon/go-request/query"
)

type Output interface {
	Validate() []common.ParamError
}

func Decode(r *http.Request, queryOutput Output, bodyOutput Output) common.Problem {
	var ok bool
	var inputProblem *common.InputProblem
	var unexpectedProblem *common.UnexpectedProblem
	if queryOutput != nil {
		queryDecoder := query.NewDecoder(r)
		err := queryDecoder.Decode(queryOutput)
		if prob, ok := err.(*common.InputProblem); ok && prob != nil {
			inputProblem = prob
		}
		if unexpectedProblem, ok = err.(*common.UnexpectedProblem); ok && unexpectedProblem != nil {
			return unexpectedProblem
		}
		queryErrors := queryOutput.Validate()
		if len(queryErrors) > 0 {
			if inputProblem == nil {
				prob := common.DefaultInputProblem
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
			return &common.DefaultUnexpectedProblem
		}
		bodyErrors := bodyOutput.Validate()
		if len(bodyErrors) > 0 {
			if inputProblem == nil {
				prob := common.DefaultInputProblem
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
