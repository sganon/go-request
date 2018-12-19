package query

import (
	"encoding"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/fatih/structtag"
	request "github.com/sganon/go-request"
)

// Decoder handles unmarshalling and validation of its request's query
type Decoder struct {
	r              *http.Request
	BoolStrictMode bool
	InputProblem   *request.InputProblem
}

// NewDecoder return a pointer to a new decoder
func NewDecoder(r *http.Request) *Decoder {
	return &Decoder{
		r:              r,
		BoolStrictMode: true,
	}
}

var defaultErr = request.UnexpectedProblem{
	Type:   "about:blank",
	Title:  "An unexpected error occured decoding request",
	Status: http.StatusInternalServerError,
}

// Decode input data from its request and stores it onto i.
func (d *Decoder) Decode(v interface{}) error {
	// call ParseForm to prepare query extraction
	if err := d.r.ParseForm(); err != nil {
		return defaultErr
	}

	elem := reflect.ValueOf(v).Elem()
	// Loop through each fields of v checking for correct input rules.
	for i := 0; i < elem.NumField(); i++ {
		fieldTag := elem.Type().Field(i).Tag

		tags, err := structtag.Parse(string(fieldTag))
		inputTag, err := tags.Get("request")
		if err != nil {
			// Skip if field has no input tag
			continue
		}

		d.extractQuery(inputTag.Name, inputTag.HasOption("required"), elem.Field(i))
	}
	return d.InputProblem
}

func (d *Decoder) extractQuery(key string, required bool, e reflect.Value) {
	val := d.r.FormValue(key)
	// If bool strict mode is deactivated having ?<key> will be evaluated as true
	if val == "" && required && e.Type().Name() == "bool" && !d.BoolStrictMode {
		e.SetBool(true)
		return
	}
	if val == "" && required {
		d.addParamsError(request.ParamError{
			Field:  key,
			Reason: "parameter is required",
		})
		return
	} else if val == "" && !required {
		return
	}
	d.setFromType(e, key, val)
}

func (d *Decoder) addParamsError(e request.ParamError) {
	if d.InputProblem == nil {
		d.initInputProblem()
	}
	d.InputProblem.InvalidParams = append(d.InputProblem.InvalidParams, e)
}

func (d *Decoder) initInputProblem() {
	d.InputProblem = &request.InputProblem{
		Payload: request.Payload{
			Type:   "about:blank",
			Title:  "Your request parameters didn't validate.",
			Status: http.StatusBadRequest,
		},
	}
}

func (d *Decoder) setFromType(e reflect.Value, key, val string) {
	switch e.Type().Name() {
	case "string":
		e.SetString(val)
		break
	case "int64":
	case "int":
		v, err := strconv.Atoi(val)
		if err != nil {
			d.addParamsError(request.ParamError{
				Field:  key,
				Reason: "syntax error: unable to convert to integer",
			})
			break
		}
		e.SetInt(int64(v))
		break
	case "float32":
	case "float64":
		v, err := strconv.ParseFloat(val, 32)
		if err != nil {
			d.addParamsError(request.ParamError{
				Field:  key,
				Reason: "syntax error: unable to convert to float",
			})
			break
		}
		e.SetFloat(float64(v))
		break
	case "bool":
		v, err := strconv.ParseBool(val)
		if err != nil {
			d.addParamsError(request.ParamError{
				Field:  key,
				Reason: "syntax error: unable to convert to bool",
			})
			break
		}
		e.SetBool(v)
		break
	default:
		var err error
		targetType := reflect.PtrTo(e.Type())
		if targetType.Implements(request.TextUnmarshalerType) {
			err = e.Addr().Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(val))
		} else if targetType.Implements(request.StringSetterType) {
			err = e.Addr().Interface().(request.StringSetter).Set(val)
		}
		if err != nil {
			d.addParamsError(request.ParamError{
				Field:  key,
				Reason: fmt.Sprintf("an error occured via UnmarshalText: %v", err),
			})
		}
	}
}
