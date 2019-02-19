# go-request
[![Build](https://travis-ci.org/sganon/go-request.svg?branch=master)](https://travis-ci.org/sganon/go-request)
[![codecov](https://codecov.io/gh/sganon/go-request/branch/master/graph/badge.svg)](https://codecov.io/gh/sganon/go-request)
[![godoc](https://godoc.org/github.com/sganon/go-request?status.svg)](http://godoc.org/github.com/sganon/go-request)

## About

This package was develop to ease request parameters decoding, both in form (i.e. query, url-encoded, multipart etc.) and in JSON body. It also provide
an error handling based on [RFC-7807](https://tools.ietf.org/html/rfc7807), the errors are those from
decoding (i.e a type error), and those you defines.

**NB**: This package is still a work in progress. The API is not definitive at all and has it flaws.
Also the rfc implementation is just the bare minimum.

## Get Started

First of all you need to define models for your parameters, one for your form parameters, and one for your body (or just one of those if you only need one ofc). 

A model contains field with some tags and a `Validate` methods where you can implement your custom
error handling.

```go
// Form represents the expected form parameters of your request
type Form struct {
  Foo string        `form:"foo"` // the 'foo' form parameter will be stored into this field
  Bar int           `form:"bar,required"` // this parameter is required
  IDs form.IntList  `form:"ids"` // this type is an helper to store `ids=1,2,3....` into an int slice
  // You can set custom types but they need to implement encoding.TextUnmarshaler or query.StringSetter
}

// Validate implements request.Output
func (q Form) Validate() (problems []problem.ParamError) {
  // Here you can set up your custom validation
  if q.Foo != "allowed_value" {
    problems = append(problems, problem.ParamError{
      Field: "foo",
      Reason: "the value "+q.Foo+" is not allowed",
    })
  }
  // ...
  return problems
}
```
The body model is simply a struct containing correct `json` tags for its fields and implementing request.Output.

Once your models are defined you can now decode your requests:
```go
import (
  "http"
  
  request "github.com/sganon/go-request"
)

func YourHandler(w http.ResponseWriter, r *http.Request) {
  // see explanations above for these types
  var form Form
  var body Body

  // You can pass nil for form or body if you don't want it to be decoded
  problem := request.Decode(r, &form, &body)
  if problem != nil {
    // An input problem has been detected, send it and stop
    problem.Send(w)
    return
  }
  // From here form, and body has been populated and their values validated
}
```

## TODO
- [ ] Make sure API is coherent in most use cases
- [ ] Implement path parameters decoding for package like [httprouter](https://github.com/julienschmidt/httprouter)
- [ ] Be more strict on the RFC implementation (e.g Type should'nt be about:blank)