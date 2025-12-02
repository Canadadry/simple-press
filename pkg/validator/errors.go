package validator

import (
	"app/pkg/maputil"
	"encoding/json"
	"fmt"
	"io"
)

type Errors struct {
	Errors   map[string][]string
	HasError bool
}

func (es *Errors) Accumulate(name string, err error) {
	if es.Errors == nil {
		es.Errors = map[string][]string{}
	}
	_, ok := es.Errors[name]
	if !ok {
		es.Errors[name] = []string{}
	}
	if err != nil {
		es.Errors[name] = append(es.Errors[name], err.Error())
		es.HasError = true
	}
}

func (es Errors) Error() string {
	return fmt.Sprintf("%v", es.Errors)
}

func (es Errors) ContentType() string {
	return "application/json"
}

func (es Errors) Write(w io.Writer) {
	_ = json.NewEncoder(w).Encode(maputil.Expand(es.Errors, "."))
}
