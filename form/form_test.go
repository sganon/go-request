package form_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/sganon/go-request/form"
	"github.com/sganon/go-request/problem"
	"github.com/stretchr/testify/assert"
)

type test struct {
	Query     string
	Input     input
	Output    input
	ShouldErr bool
}

type input struct {
	Foo            string          `form:"foo"`
	ID             int             `form:"id"`
	IsPresent      bool            `form:"is_present"`
	ImportantField string          `form:"imp,required"`
	Temperature    float32         `form:"temperature"`
	IDList         form.IntList    `form:"id_list"`
	Authors        form.StringList `form:"authors"`
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
		Query: "?imp=here&is_present=vrai",
		Input: input{},
		Output: input{
			ImportantField: "here",
		},
		ShouldErr: true,
	},
	{
		Query: "?imp=here&is_present=false",
		Input: input{},
		Output: input{
			ImportantField: "here",
		},
	},
	{
		Query: "?imp=here&temperature=42.21",
		Input: input{},
		Output: input{
			ImportantField: "here",
			Temperature:    42.21,
		},
	},
	{
		Query: "?imp=here&temperature=not_a_float",
		Input: input{},
		Output: input{
			ImportantField: "here",
		},
		ShouldErr: true,
	},
	{
		Query: "?imp=here&id_list=21,42,84",
		Input: input{},
		Output: input{
			ImportantField: "here",
			IDList:         form.IntList{21, 42, 84},
		},
	},
	{
		Query: "?imp=here&id_list=21,AB,84",
		Input: input{},
		Output: input{
			ImportantField: "here",
		},
		ShouldErr: true,
	},
	{
		Query: "?imp=here&authors=sganon,ganondorf",
		Input: input{},
		Output: input{
			ImportantField: "here",
			Authors:        form.StringList{"sganon", "ganondorf"},
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

func TestFormDecoder(t *testing.T) {
	var body io.Reader
	for i, te := range tests {
		req, err := http.NewRequest("GET", "http://localhost:80/"+te.Query, body)
		assert.Nil(t, err, "request should have been created", i)
		decoder := form.NewDecoder(req)
		err = decoder.Decode(&te.Input)
		if te.ShouldErr {
			assert.Equal(t, problem.ErrInvalidParameters.Error(), err.Error(), "error should be equal to predefined one", i)
			assert.NotNil(t, err, "decode should have returned an error", i)
		}
		assert.Equal(t, te.Output, te.Input, "input should have been correctly decoded", i)
	}
}
