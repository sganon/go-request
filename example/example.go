package example

import (
	"fmt"
	"net/http"

	"github.com/sganon/go-request"
	"github.com/sganon/go-request/problem"
	"github.com/sganon/go-request/query"
)

// Define a the model of an expected query
type Query struct {
	Foo   string           `request:"foo"`
	IDs   query.IntList    `request:"ids,required"`
	Names query.StringList `request:"names"`
}

// Validate implements request.Output
// This is executed after request decoding
// You can implement specific logic in here
func (q Query) Validate() (problems []problem.ParamError) {
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

// QueryHandler shows how to handles query decoding and error response handling
// This handler is used in example_test.go.
func QueryHandler(w http.ResponseWriter, r *http.Request) {
	var query Query
	problem := request.Decode(r, &query, nil)
	if problem != nil {
		problem.Send(w)
		return
	}
	// From here the query is decoded, validated
}
