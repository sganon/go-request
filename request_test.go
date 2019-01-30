package request_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	request "github.com/sganon/go-request"
	"github.com/sganon/go-request/problem"
	"github.com/stretchr/testify/assert"
)

type inputQuery struct {
	Foo string `request:"foo,required"`
}

func (q inputQuery) Validate() (errs []problem.ParamError) {
	if q.Foo != "bar" {
		errs = append(errs, problem.ParamError{
			Field:  "foo",
			Reason: "in query foo key should have `bar` value",
		})
	}
	return errs
}

type inputBody struct {
	Foo string `json:"foo"`
}

func (b inputBody) Validate() (errs []problem.ParamError) {
	if b.Foo != "baz" {
		errs = append(errs, problem.ParamError{
			Field:  "foo",
			Reason: "in body foo key should have `baz` value",
		})
	}
	return errs
}

type test struct {
	Method         string
	Query          string
	ExpectedStatus int
	ErrorsLen      int
	Body           inputBody
}

var tests = []test{
	{
		Method:         "GET",
		Query:          "?",
		ExpectedStatus: http.StatusBadRequest,
		ErrorsLen:      2,
	},
	{
		Method:         "GET",
		Query:          "?foo=bar",
		ExpectedStatus: http.StatusOK,
	},
	{
		Method:         "POST",
		Query:          "?foo=bar",
		Body:           inputBody{Foo: "baz"},
		ExpectedStatus: http.StatusOK,
	},
	{
		Method:         "POST",
		Query:          "?foo=bar",
		Body:           inputBody{Foo: "bur"},
		ExpectedStatus: http.StatusBadRequest,
		ErrorsLen:      1,
	},
	{
		Method:         "POST",
		Query:          "?foo=nope",
		Body:           inputBody{Foo: "bur"},
		ExpectedStatus: http.StatusBadRequest,
		ErrorsLen:      2,
	},
}

func TestDecode(t *testing.T) {
	ts := httptest.NewServer(handlerFunc)
	defer ts.Close()
	for _, te := range tests {
		body := new(bytes.Buffer)
		if te.Body != (inputBody{}) {
			encoder := json.NewEncoder(body)
			err := encoder.Encode(te.Body)
			assert.NoError(t, err)
		}
		req, err := http.NewRequest(te.Method, fmt.Sprintf("%s/%s", ts.URL, te.Query), body)
		assert.NoError(t, err)
		client := http.DefaultClient
		res, err := client.Do(req)
		assert.NoError(t, err)

		assert.Equal(t, te.ExpectedStatus, res.StatusCode)
		if res.StatusCode != http.StatusOK {
			assert.Equal(t, "application/problem+json", res.Header.Get("Content-Type"))
		}
		if res.StatusCode == http.StatusBadRequest {
			var resBody problem.Input
			decoder := json.NewDecoder(res.Body)
			defer res.Body.Close()
			err := decoder.Decode(&resBody)
			assert.NoError(t, err)
			assert.Equal(t, te.ErrorsLen, len(resBody.InvalidParams))
		}
	}
}

var handlerFunc = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var query inputQuery
	var body inputBody
	var problem problem.Problem
	if r.Method == "POST" {
		problem = request.Decode(r, &query, &body)
	} else {
		problem = request.Decode(r, &query, nil)
	}
	if problem != nil {
		problem.Send(w)
		return
	}
})
