package example

import (
	"fmt"
	"net/http"

	"github.com/sganon/go-request"
	"github.com/sganon/go-request/problem"
	"github.com/sganon/go-request/form"
)

// Form Define a the model of an expected form
type Form struct {
	Foo   string           `form:"foo"`
	IDs   form.IntList    `form:"ids,required"`
	Names form.StringList `form:"names"`
}

// Validate implements request.Output
// This is executed after request decoding
// You can implement specific logic in here
func (q Form) Validate() (problems []problem.ParamError) {
	if q.Foo != "" && q.Foo != "bar" {
		problems = append(problems, problem.ParamError{
			Field:  "foo",
			Reason: fmt.Sprintf("value %s is not acceptable", q.Foo),
		})
	}
	if len(q.Names) > 3 {
		problems = append(problems, problem.ParamError{
			Field:  "names",
			Reason: "too much names given",
		})
	}
	return problems
}

// QueryHandler shows how to handles form decoding and error response handling
// This handler is used in example_test.go.
func QueryHandler(w http.ResponseWriter, r *http.Request) {
	var form Form
	problem := request.Decode(r, &form, nil)
	if problem != nil {
		problem.Send(w)
		return
	}
	// From here the form is decoded, validated
}
