package form

import (
	"encoding"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"

	"github.com/fatih/structtag"
	"github.com/sganon/go-request/problem"
)

const (
	mb = 1 << 20
)

// Decoder handles unmarshalling and validation of its request's query
type Decoder struct {
	r              *http.Request
	BoolStrictMode bool
	Input          *problem.Input
}

// NewDecoder return a pointer to a new decoder
func NewDecoder(r *http.Request) *Decoder {
	return &Decoder{
		r:              r,
		BoolStrictMode: true,
	}
}

// Decode input data from its request and stores it onto i.
func (d *Decoder) Decode(v interface{}) error {
	elem := reflect.ValueOf(v).Elem()
	// Loop through each fields of v checking for correct input rules.
	for i := 0; i < elem.NumField(); i++ {
		fieldTag := elem.Type().Field(i).Tag

		tags, err := structtag.Parse(string(fieldTag))
		if err != nil {
			return problem.DefaultUnexpected
		}
		queryTag, noFormErr := tags.Get("form")
		fileTag, noFileErr := tags.Get("file")
		if noFormErr != nil && noFileErr != nil {
			// Skip if field has no input tag
			continue
		}

		if queryTag != nil {
			d.extractForm(queryTag.Name, queryTag.HasOption("required"), elem.Field(i))
		}
		if fileTag != nil {
			if err := d.r.ParseMultipartForm(5 * mb); err != nil {
				return problem.DefaultUnexpected
			}
			e := elem.Field(i)
			file, _, err := d.r.FormFile(fileTag.Name)
			if err != nil && fileTag.HasOption("required") {
				d.addParamsError(problem.ParamError{
					Field:  fileTag.Name,
					Reason: "this file is required",
				})
				continue
			} else if err != nil && !fileTag.HasOption("required") {
				continue
			}

			b, err := ioutil.ReadAll(file)
			if err != nil {
				return problem.DefaultUnexpected
			}
			if err = e.Addr().Interface().(encoding.BinaryUnmarshaler).UnmarshalBinary(b); err != nil {
				d.addParamsError(problem.ParamError{
					Field:  fileTag.Name,
					Reason: fmt.Sprintf("an error occured via UnmarshalBinary: %v", err),
				})
			}
		}
	}
	return d.Input
}

func (d *Decoder) extractForm(key string, required bool, e reflect.Value) {
	val := d.r.FormValue(key)
	fmt.Println(key, val)
	// If bool strict mode is deactivated having ?<key> will be evaluated as true
	if val == "" && required && e.Type().Name() == "bool" && !d.BoolStrictMode {
		e.SetBool(true)
		return
	}
	if val == "" && required {
		d.addParamsError(problem.ParamError{
			Field:  key,
			Reason: "parameter is required",
		})
		return
	} else if val == "" && !required {
		return
	}
	d.setFromType(e, key, val)
}

func (d *Decoder) addParamsError(e problem.ParamError) {
	if d.Input == nil {
		d.initInputProblem()
	}
	d.Input.InvalidParams = append(d.Input.InvalidParams, e)
}

func (d *Decoder) initInputProblem() {
	prob := problem.DefaultInput
	prob.Title = "Your form parameters could not be decoded"
	d.Input = &prob
}

func (d *Decoder) setFromType(e reflect.Value, key, val string) {
	switch e.Type().Name() {
	case "string":
		e.SetString(val)
		break
	case "int64":
		fallthrough
	case "int":
		v, err := strconv.Atoi(val)
		if err != nil {
			d.addParamsError(problem.ParamError{
				Field:  key,
				Reason: "syntax error: unable to convert to integer",
			})
			break
		}
		e.SetInt(int64(v))
		break
	case "float32":
		fallthrough
	case "float64":
		v, err := strconv.ParseFloat(val, 32)
		if err != nil {
			d.addParamsError(problem.ParamError{
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
			d.addParamsError(problem.ParamError{
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
		if targetType.Implements(TextUnmarshalerType) {
			err = e.Addr().Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(val))
		} else if targetType.Implements(StringSetterType) {
			err = e.Addr().Interface().(StringSetter).Set(val)
		}
		if err != nil {
			d.addParamsError(problem.ParamError{
				Field:  key,
				Reason: fmt.Sprintf("an error occured via UnmarshalText: %v", err),
			})
		}
	}
}
