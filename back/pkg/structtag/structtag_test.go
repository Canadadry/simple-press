package structtag

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {

	tests := map[string]struct {
		tag     string
		exp     []Tag
		invalid bool
	}{
		"empty tag": {
			exp: []Tag{},
		},
		"tag with one key (invalid)": {
			tag:     "json",
			invalid: true,
			exp:     []Tag{},
		},
		"tag with one key (valid)": {
			tag: `json:""`,
			exp: []Tag{
				{Key: "json"},
			},
		},
		"tag with one key and dash name": {
			tag: `json:"-"`,
			exp: []Tag{
				{Key: "json", Name: "-"},
			},
		},
		"tag with key and name": {
			tag: `json:"foo"`,
			exp: []Tag{
				{Key: "json", Name: "foo"},
			},
		},
		"tag with key, name and option": {
			tag: `json:"foo,omitempty"`,
			exp: []Tag{
				{Key: "json", Name: "foo", Options: []string{"omitempty"}},
			},
		},
		"tag with multiple keys": {
			tag: `json:"" hcl:""`,
			exp: []Tag{
				{Key: "json"},
				{Key: "hcl"},
			},
		},
		"tag with multiple keys and names": {
			tag: `json:"foo" hcl:"foo"`,
			exp: []Tag{
				{Key: "json", Name: "foo"},
				{Key: "hcl", Name: "foo"},
			},
		},
		"tag with multiple keys and different names": {
			tag: `json:"foo" hcl:"bar"`,
			exp: []Tag{
				{Key: "json", Name: "foo"},
				{Key: "hcl", Name: "bar"},
			},
		},
		"tag with multiple keys, different names and options": {
			tag: `json:"foo,omitempty" structs:"bar,omitnested" hcl:"-"`,
			exp: []Tag{
				{Key: "json", Name: "foo", Options: []string{"omitempty"}},
				{Key: "structs", Name: "bar", Options: []string{"omitnested"}},
				{Key: "hcl", Name: "-"},
			},
		},
		"tag with quoted name": {
			tag: `json:"foo,bar:\"baz\""`,
			exp: []Tag{
				{Key: "json", Name: "foo", Options: []string{`bar:"baz"`}},
			},
		},
		"tag with trailing space": {
			tag: `json:"foo" `,
			exp: []Tag{
				{Key: "json", Name: "foo"},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := Parse(tt.tag)
			invalid := err != nil

			if invalid != tt.invalid {
				t.Errorf("invalid case\n\twant: %+v\n\tgot : %+v\n\terr : %s", tt.invalid, invalid, err)
			}

			if invalid {
				return
			}

			if !reflect.DeepEqual(tt.exp, result) {
				t.Errorf("parse\n\twant: %#v\n\tgot : %#v", tt.exp, result)
			}
		})
	}
}
