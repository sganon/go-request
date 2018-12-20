package request_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	request "github.com/sganon/go-request"
	"github.com/stretchr/testify/assert"
)

type test struct {
	Method         string
	Query          string
	ExpectedStatus int
}

type inputQuery struct {
	Foo string `request:"foo,required"`
}

var tests = []test{
	{
		Method:         "GET",
		Query:          "?",
		ExpectedStatus: http.StatusBadRequest,
	},
	{
		Method:         "GET",
		Query:          "?foo=bar",
		ExpectedStatus: http.StatusOK,
	},
}

func TestDecode(t *testing.T) {
	ts := httptest.NewServer(handlerFunc)
	defer ts.Close()
	for _, te := range tests {
		var body io.Reader
		req, err := http.NewRequest(te.Method, fmt.Sprintf("%s/%s", ts.URL, te.Query), body)
		assert.NoError(t, err)
		client := http.DefaultClient
		res, err := client.Do(req)
		assert.NoError(t, err)

		assert.Equal(t, te.ExpectedStatus, res.StatusCode)
		if res.StatusCode != http.StatusOK {
			assert.Equal(t, "application/problem+json", res.Header.Get("Content-Type"))
		}
	}
}

var handlerFunc = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var query inputQuery
	problem := request.Decode(r, &query, nil)
	if problem != nil {
		problem.Send(w)
		return
	}
})
