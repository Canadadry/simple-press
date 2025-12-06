package validator

import (
	"app/pkg/null"
	"fmt"
	"strings"
	"testing"
)

func TestValidateRequieredField(t *testing.T) {
	tests := map[string]struct {
		in    map[string]interface{}
		funcs []func(val string) error
		exp   string
	}{
		"basic error": {
			in:    map[string]interface{}{"test": "test"},
			funcs: []func(val string) error{Length(1, 3)},
			exp:   `"test":higher than max length of 3`,
		},
		"basic succes": {
			in:    map[string]interface{}{"test": "tes"},
			funcs: []func(val string) error{Length(1, 3)},
			exp:   `"tes":`,
		},
		"key not found": {
			in:    map[string]interface{}{"tet": "tes"},
			funcs: []func(val string) error{Length(1, 3)},
			exp:   `"":this value should not be blank`,
		},
		"basic valid email": {
			in:    map[string]interface{}{"test": "test@example.com"},
			funcs: []func(val string) error{Regexp(EmailRe)},
			exp:   `"test@example.com":`,
		},
		"basic invalid email": {
			in:    map[string]interface{}{"test": "example.com"},
			funcs: []func(val string) error{Regexp(EmailRe)},
			exp:   `"example.com":value dont match format`,
		},
		"basic invalid email mx": {
			in:    map[string]interface{}{"test": "everycheck.io"},
			funcs: []func(val string) error{EmailMX},
			exp:   `"everycheck.io":mx record dont exist on everycheck.io`,
		},
		"basic valid email mx": {
			in:    map[string]interface{}{"test": "test@everycheck.com"},
			funcs: []func(val string) error{EmailMX},
			exp:   `"test@everycheck.com":`,
		},
	}

	for title, tt := range tests {
		t.Run(title, func(t *testing.T) {
			out := ""
			errs := BindWithMap(tt.in, func(b Binder) {
				b.RequiredStringVar("test", &out, tt.funcs...)
			})

			result := fmt.Sprintf("%#v:%v", out, strings.Join(errs.Errors["test"], ","))
			if result != tt.exp {
				t.Fatalf("\nexp '%s'\ngot '%s'\n", tt.exp, result)
			}

		})
	}
}

func TestValidateOptionalField(t *testing.T) {
	tests := map[string]struct {
		in    map[string]interface{}
		funcs []func(val string) error
		exp   string
	}{
		"basic error": {
			in:    map[string]interface{}{"test": "test"},
			funcs: []func(val string) error{Length(1, 3)},
			exp:   `"test":higher than max length of 3`,
		},
		"basic succes": {
			in:    map[string]interface{}{"test": "tes"},
			funcs: []func(val string) error{Length(1, 3)},
			exp:   `"tes":`,
		},
		"key not found but not required": {
			in:    map[string]interface{}{"tet": "tes"},
			funcs: []func(val string) error{Length(1, 3)},
			exp:   `"emp":`,
		},
		"key found but nil as value": {
			in:    map[string]interface{}{"test": nil},
			funcs: []func(val string) error{},
			exp:   `"emp":this is not a valid value. Expect a string`,
		},
		"key not found and default value dont match constraint": {
			in: map[string]interface{}{},
			funcs: []func(val string) error{
				func(val string) error {
					return fmt.Errorf("invalid value")
				},
			},
			exp: `"emp":`,
		},
	}

	for title, tt := range tests {
		t.Run(title, func(t *testing.T) {
			out := null.New[string]("emp", true)
			errs := BindWithMap(tt.in, func(b Binder) {
				b.StringVar("test", &out, tt.funcs...)
			})

			result := fmt.Sprintf("%#v:%v", out.V, strings.Join(errs.Errors["test"], ","))
			if result != tt.exp {
				t.Fatalf("\nexp '%s'\ngot '%s'\n", tt.exp, result)
			}
		})
	}
}
