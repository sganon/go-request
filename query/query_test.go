package query_test

import (
	"io"
	"net/http"
	"testing"

	request "github.com/sganon/go-request"
	"github.com/sganon/go-request/query"
	"github.com/stretchr/testify/assert"
)

type test struct {
	Query     string
	Input     input
	Output    input
	ShouldErr bool
}

type input struct {
	Foo            string `request:"foo"`
	ID             int    `request:"id"`
	IsPresent      bool   `request:"is_present"`
	ImportantField string `request:"imp,required"`
}

var tests = []test{
	{
		Query: "?imp=here",
		Input: input{},
		Output: input{
			ImportantField: "here",
		},
	},
	{
		Query: "?imp=here&foo=bar",
		Input: input{},
		Output: input{
			ImportantField: "here",
			Foo:            "bar",
		},
	},
	{
		Query: "?imp=here&id=42",
		Input: input{},
		Output: input{
			ImportantField: "here",
			ID:             42,
		},
	},
	{
		Query: "?imp=here&id=fourtytwo",
		Input: input{},
		Output: input{
			ImportantField: "here",
		},
		ShouldErr: true,
	},
	{
		Query: "?imp=here&is_present=true",
		Input: input{},
		Output: input{
			ImportantField: "here",
			IsPresent:      true,
		},
	},
	{
		Query: "?imp=here&is_present=false",
		Input: input{},
		Output: input{
			ImportantField: "here",
		},
	},
	{
		Query:     "",
		Input:     input{},
		Output:    input{},
		ShouldErr: true,
	},
}

var suite []test

func TestQueryDecoder(t *testing.T) {
	var body io.Reader
	for _, te := range tests {
		req, err := http.NewRequest("GET", "http://localhost:80/"+te.Query, body)
		assert.Nil(t, err, "request should have been created")
		decoder := query.NewDecoder(req)
		err = decoder.Decode(&te.Input)
		if te.ShouldErr {
			assert.Equal(t, request.ErrInvalidParameters.Error(), err.Error(), "error should be equal to predefined one")
			assert.NotNil(t, err, "decode should have returned an error")
		}
		assert.Equal(t, te.Output, te.Input, "input should have been correctly decoded")
	}
}
