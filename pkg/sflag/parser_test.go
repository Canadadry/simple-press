package sflag

import (
	"flag"
	"fmt"
	"reflect"
	"testing"
)

type MockedFlagSet []string

func (mfs *MockedFlagSet) StringVar(ptr *string, name string, value string, desc string) {
	*mfs = append(*mfs, fmt.Sprintf("string '%s' default '%v' desc '%s'", name, value, desc))
}
func (mfs *MockedFlagSet) IntVar(ptr *int, name string, value int, desc string) {
	*mfs = append(*mfs, fmt.Sprintf("int '%s' default '%v' desc '%s'", name, value, desc))
}
func (mfs *MockedFlagSet) BoolVar(ptr *bool, name string, value bool, desc string) {
	*mfs = append(*mfs, fmt.Sprintf("bool '%s' default '%v' desc '%s'", name, value, desc))
}

func (mfs *MockedFlagSet) Match(t *testing.T, exp []string) {
	got := []string((*mfs))
	if len(exp) == len(got) && len(exp) == 0 {
		return
	}
	if !reflect.DeepEqual(exp, got) {
		t.Errorf("parse\n\twant: %#v\n\tgot : %#v", exp, got)
	}
}

func TestParse(t *testing.T) {
	value := 0
	var i interface{}
	i = struct{}{}
	tests := map[string]struct {
		data    interface{}
		flag    []string
		invalid bool
	}{
		"int instead of pointer of struct": {
			data:    1,
			invalid: true,
		},
		"pointer to int instead of pointer of struct": {
			data:    &value,
			invalid: true,
		},
		"pointer to interface instead of pointer of struct": {
			data:    &i,
			invalid: true,
		},
		"empty pointer of struct": {
			data: &struct{}{},
		},
		"one line without tag": {
			data: &struct{ Test string }{},
		},
		"one line with wrong tag": {
			data: &struct {
				Test string `toto:"-"`
			}{},
		},
		"one line with no name flag tag": {
			data: &struct {
				Test string `flag:"-"`
			}{},
		},
		"one line with  named flag tag": {
			data: &struct {
				Test string `flag:"test"`
			}{},
			flag: []string{
				"string 'test' default '' desc ''",
			},
		},
		"one line with named flag tag and desc": {
			data: &struct {
				Test string `flag:"test" desc:"this is test"`
			}{},
			flag: []string{
				"string 'test' default '' desc 'this is test'",
			},
		},
		"one line with named flag tag , desc and default value": {
			data: &struct {
				Test string `flag:"test" default:"default test value" desc:"this is test"`
			}{},
			flag: []string{
				"string 'test' default 'default test value' desc 'this is test'",
			},
		},
		"multi string rich lines": {
			data: &struct {
				Test1 string `flag:"test1" default:"default test1 value" desc:"this is test1"`
				Test2 string `flag:"test2" default:"default test2 value" desc:"this is test2"`
				Test3 string `flag:"test3" default:"default test3 value" desc:"this is test3"`
				Test4 string `flag:"test4" default:"default test4 value" desc:"this is test4"`
			}{},
			flag: []string{
				"string 'test1' default 'default test1 value' desc 'this is test1'",
				"string 'test2' default 'default test2 value' desc 'this is test2'",
				"string 'test3' default 'default test3 value' desc 'this is test3'",
				"string 'test4' default 'default test4 value' desc 'this is test4'",
			},
		},
		"multi rich lines with different types": {
			data: &struct {
				Test1 string `flag:"test1" default:"default test1 value" desc:"this is test1"`
				Test2 bool   `flag:"test2" default:"false" desc:"this is test2"`
				Test3 int    `flag:"test3" default:"3" desc:"this is test3"`
			}{},
			flag: []string{
				"string 'test1' default 'default test1 value' desc 'this is test1'",
				"bool 'test2' default 'false' desc 'this is test2'",
				"int 'test3' default '3' desc 'this is test3'",
			},
		},
		"nested struct": {
			data: &struct {
				Test1  string `flag:"test1" default:"default test1 value" desc:"this is test1"`
				Struct struct {
					Test2 bool `flag:"test2" default:"false" desc:"this is test2"`
					Test3 int  `flag:"test3" default:"3" desc:"this is test3"`
				}
			}{},
			flag: []string{
				"string 'test1' default 'default test1 value' desc 'this is test1'",
				"bool 'test2' default 'false' desc 'this is test2'",
				"int 'test3' default '3' desc 'this is test3'",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mfs := MockedFlagSet{}
			err := Parse(&mfs, tt.data)
			invalid := err != nil

			if invalid != tt.invalid {
				t.Errorf("invalid case\n\twant: %+v\n\tgot : %+v\n\terr : %s", tt.invalid, invalid, err)
			}

			if invalid {
				return
			}
			mfs.Match(t, tt.flag)
		})
	}
}

func TestRealFlagParse(t *testing.T) {

	tests := map[string]struct {
		data interface{}
		args []string
		exp  string
	}{
		"nested struct": {
			data: &struct {
				Test1  string `flag:"test1" default:"default test1 value" desc:"this is test1"`
				Struct struct {
					Test2 bool `flag:"test2" default:"false" desc:"this is test2"`
					Test3 int  `flag:"test3" default:"3" desc:"this is test3"`
				}
			}{},
			args: []string{"-test2"},
			exp:  "&{default test1 value {true 3}}",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			fs := flag.NewFlagSet("test", flag.ContinueOnError)
			err := Parse(fs, tt.data)
			if err != nil {
				t.Errorf("failed : %s", err)
			}
			err = fs.Parse(tt.args)
			if err != nil {
				t.Errorf("failed : %s", err)
			}
			got := fmt.Sprintf("%v", tt.data)
			if got != tt.exp {
				t.Fatalf("got %s want %s", got, tt.exp)
			}
		})
	}
}
