package null_test

import (
	"app/pkg/null"
	"encoding/json"
	"testing"
)

func TestNullableString(t *testing.T) {
	tests := map[string]struct {
		variable null.Nullable[any]
		result   string
	}{
		"valid float":    {variable: null.New[any](0., true), result: "0"},
		"invalid float":  {variable: null.New[any](0., false), result: ""},
		"valid int":      {variable: null.New[any](0, true), result: "0"},
		"invalid int":    {variable: null.New[any](0, false), result: ""},
		"valid string":   {variable: null.New[any]("une chaine de caractere", true), result: "une chaine de caractere"},
		"invalid string": {variable: null.New[any]("", false), result: ""},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := test.variable.String()
			if result != test.result {
				t.Fatalf("failed test %s should have %v but got %v", name, test.result, result)
			}
		})
	}
}

func TestNullableJson(t *testing.T) {
	tests := map[string]null.Nullable[any]{
		"valid float":    null.New[any](1., true),
		"invalid float":  null.New[any](0., false),
		"valid int":      null.New[any](1, true),
		"invalid int":    null.New[any](0, false),
		"valid string":   null.New[any]("une chaine de caractere", true),
		"invalid string": null.New[any]("", false),
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			u, _ := test.MarshalJSON()

			var tmp null.Nullable[any]
			if err := json.Unmarshal(u, &tmp); err != nil {
				t.Fatalf("found error : %s", err.Error())
			}
			if tmp.Valid != test.Valid && tmp.V != test.V {
				t.Fatalf("%s have : %v (%t) want : %v(%t)", name, tmp.V, tmp.Valid, test.V, test.Valid)
			}
		})
	}

	//since the use of any imply some pointer magic
	//here is a manual test ...
	tvar, _ := null.New("", false).MarshalJSON()
	var tmp null.Nullable[string]
	json.Unmarshal(tvar, &tmp)
	if tmp.Valid {
		t.Fatalf("nullable string failed")
	}
}

func TestEmbededNullableJson(t *testing.T) {
	type Embeded struct {
		Sub null.Nullable[any] `json:"sub,omitzero"`
	}
	tests := map[string]struct {
		in  Embeded
		out string
	}{
		"valid string": {
			in:  Embeded{Sub: null.New[any]("une chaine de caractere", true)},
			out: `{"sub":"une chaine de caractere"}`,
		},
		"invalid string": {
			in:  Embeded{Sub: null.New[any]("", false)},
			out: "{}",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			u, err := json.Marshal(tt.in)
			if err != nil {
				t.Fatalf("cannot marshal %v", err)
			}
			if string(u) != tt.out {
				t.Fatalf("got %v want : %v", string(u), tt.out)
			}
		})
	}
}
