package agentsdk

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"net/http"
	"reflect"
	"strings"
)

var validate *validator.Validate

// This init is used to create a validator and register validation-specific
// functionality for the HTTP API.
//
// A single validator instance is used, because it caches struct parsing.
func init() {
	validate = validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// Write outputs a standardized format to an HTTP response body. ctx is used for
// tracing and can be nil for tracing to be disabled. Tracing this function is
// helpful because JSON marshaling can sometimes take a non-insignificant amount
// of time, and could help us catch outliers. Additionally, we can enrich span
// data a bit more since we have access to the actual interface{} we're
// marshaling, such as the number of elements in an array, which could help us
// spot routes that need to be paginated.
func Write(ctx context.Context, rw http.ResponseWriter, status int, response interface{}) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	err := enc.Encode(response)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(status)
	_, err = rw.Write(buf.Bytes())
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}
