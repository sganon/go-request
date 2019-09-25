package example_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sganon/go-request/example"
	"github.com/sganon/go-request/problem"
)

var queryTests = []struct {
	QueryString    string
	ExpectedStatus int
	ErrorsLen      int
	Errors         []problem.ParamError
}{
	{
		QueryString:    "?ids=1,2,3,4",
		ExpectedStatus: http.StatusOK,
	},
	{
		QueryString:    "?ids=1,2,3,4&foo=bar",
		ExpectedStatus: http.StatusOK,
	},
	{
		QueryString:    "?ids=1,2,3,4&names=simon,ganon",
		ExpectedStatus: http.StatusOK,
	},
	{
		QueryString:    "?ids=1,2,3,4&names=simon,pierre,jacques,ganon",
		ExpectedStatus: http.StatusOK,
	},
	{
		QueryString:    "?ids=foo,bar,baz",
		ExpectedStatus: http.StatusBadRequest,
		ErrorsLen:      1,
		Errors: []problem.ParamError{{
			Field:  "ids",
			Reason: "an error occured via UnmarshalText: strconv.Atoi: parsing \"foo\": invalid syntax",
		}},
	},
}

func TestQueryHandler(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(example.QueryHandler))
	defer ts.Close()

	for _, test := range queryTests {
		endpoint := ts.URL + test.QueryString
		t.Logf("testing endpoint %s", endpoint)
		res, err := http.Get(endpoint)
		assert.NoError(t, err)
		assert.Equal(t, test.ExpectedStatus, res.StatusCode)

		if test.ExpectedStatus != http.StatusOK {
			var body problem.Input
			decoder := json.NewDecoder(res.Body)
			err := decoder.Decode(&body)
			assert.NoError(t, err)
			assert.Equal(t, test.ErrorsLen, len(body.InvalidParams))
			assert.Equal(t, test.Errors, body.InvalidParams)
			assert.Equal(t, problem.DefaultInput.Title, body.Title)
			assert.Equal(t, problem.DefaultInput.Status, body.Status)
			assert.Equal(t, test.ExpectedStatus, body.Status)
		}
	}
}
