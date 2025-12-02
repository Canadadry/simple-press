package validator_test

import (
	"app/pkg/maputil"
	"app/pkg/null"
	"app/pkg/validator"
	"fmt"
	"reflect"
	"testing"
)

type FakeDataStruct struct {
	StringField null.Nullable[string]
	FloatField  null.Nullable[float64]
	IntField    null.Nullable[int]
}

func (f *FakeDataStruct) Bind(bm validator.Binder) {
	bm.StringVar("string_field", &f.StringField, validator.Length(0, 15))
	bm.Float64Var("float_field", &f.FloatField, validator.Min[float64](10))
	bm.IntVar("int_field", &f.IntField, validator.Min[int](0))
}

type FakeRequiredDataStruct struct {
	StringField string
	FloatField  float64
	IntField    int
}

func (f *FakeRequiredDataStruct) Bind(bm validator.Binder) {
	bm.RequiredStringVar("string_field", &f.StringField, validator.Length(0, 15))
	bm.RequiredFloat64Var("float_field", &f.FloatField, validator.Min[float64](10))
	bm.RequiredIntVar("int_field", &f.IntField, validator.Min[int](0))
}

type FakeFuncDataStruct struct {
	Field3 string
	Field4 string
}

func (f *FakeFuncDataStruct) Bind(bm validator.Binder) {
	bm.RequiredAnyFunc("field3", func(name string, ok bool, value any, errs *validator.Errors) {
		if ok {
			f.Field3 = fmt.Sprintf("%v", value)
		}
		validator.Validate[string](errs, name, f.Field3, validator.Length(0, 15))
	})
	bm.AnyFunc("field4", func(name string, ok bool, value any, errs *validator.Errors) {
		if ok {
			f.Field4 = fmt.Sprintf("%v", value)
		}
		validator.Validate[string](errs, name, f.Field4, validator.Length(0, 15))
	})
}

type SubDataStruct struct {
	SubStringField null.Nullable[string]
}

func (f *SubDataStruct) Bind(bm validator.Binder) {
	bm.StringVar("sub_field", &f.SubStringField)
}

type FakeDataWithArrayStruct struct {
	Array []SubDataStruct
}

func (f *FakeDataWithArrayStruct) Bind(bm validator.Binder) {
	bm.RequiredArrayVar("array_field", &f.Array, func() validator.Bindable { return &SubDataStruct{} })
}

type FakeMapStruct struct {
	Data map[string]interface{}
}

func (f *FakeMapStruct) Bind(bm validator.Binder) {
	bm.RequiredMapVar("data", &f.Data)
}

type FakeMapWithSubMapStruct struct {
	Subdata FakeMapStruct
}

func (f *FakeMapWithSubMapStruct) Bind(bm validator.Binder) {
	bm.MapVar("data.subdata", &f.Subdata.Data)
}

type FakeMapInArrayStruct struct {
	Array []FakeMapStruct
}

func (f *FakeMapInArrayStruct) Bind(bm validator.Binder) {
	bm.RequiredArrayVar("array_field", &f.Array, func() validator.Bindable { return &FakeMapStruct{} })
}

func TestBinderData(t *testing.T) {
	tests := map[string]struct {
		In             map[string]interface{}
		BindTo         validator.Bindable
		ExpectedData   string
		ExpectedErrors map[string][]string
	}{
		"Valid input data": {
			In: map[string]interface{}{
				"string_field": "valid",
				"float_field":  12.5,
				"int_field":    42,
			},
			BindTo:       &FakeDataStruct{},
			ExpectedData: `&validator_test.FakeDataStruct{StringField:null.Nullable[string]{V:"valid", Valid:true}, FloatField:null.Nullable[float64]{V:12.5, Valid:true}, IntField:null.Nullable[int]{V:42, Valid:true}}`,
			ExpectedErrors: map[string][]string{
				"string_field": {},
				"float_field":  {},
				"int_field":    {},
			},
		},
		"Empty struct": {
			In:           map[string]interface{}{},
			BindTo:       &FakeDataStruct{},
			ExpectedData: `&validator_test.FakeDataStruct{StringField:null.Nullable[string]{V:"", Valid:false}, FloatField:null.Nullable[float64]{V:0, Valid:false}, IntField:null.Nullable[int]{V:0, Valid:false}}`,
			ExpectedErrors: map[string][]string{
				"string_field": {},
				"float_field":  {},
				"int_field":    {},
			},
		},
		"Empty requiered struct": {
			In:           map[string]interface{}{},
			BindTo:       &FakeRequiredDataStruct{},
			ExpectedData: `&validator_test.FakeRequiredDataStruct{StringField:"", FloatField:0, IntField:0}`,
			ExpectedErrors: map[string][]string{
				"string_field": {"this value should not be blank"},
				"float_field":  {"this value should not be blank"},
				"int_field":    {"this value should not be blank"},
			},
		},
		"Zero value required struct": {
			In: map[string]interface{}{
				"string_field": "",
				"float_field":  0.0,
				"int_field":    0,
			},
			BindTo:       &FakeRequiredDataStruct{},
			ExpectedData: `&validator_test.FakeRequiredDataStruct{StringField:"", FloatField:0, IntField:0}`,
			ExpectedErrors: map[string][]string{
				"string_field": {},
				"float_field":  {"this value should be higher or equal to 10"},
				"int_field":    {},
			},
		},
		"String too long field": {
			In: map[string]interface{}{
				"string_field": "12345678901234567890",
				"float_field":  12.5,
				"int_field":    42,
			},
			BindTo:       &FakeDataStruct{},
			ExpectedData: `&validator_test.FakeDataStruct{StringField:null.Nullable[string]{V:"12345678901234567890", Valid:true}, FloatField:null.Nullable[float64]{V:12.5, Valid:true}, IntField:null.Nullable[int]{V:42, Valid:true}}`,
			ExpectedErrors: map[string][]string{
				"string_field": {"higher than max length of 15"},
				"float_field":  {},
				"int_field":    {},
			},
		},
		"Negative float field": {
			In: map[string]interface{}{
				"string_field": "valid",
				"float_field":  -5.0,
				"int_field":    42,
			},
			BindTo:       &FakeDataStruct{},
			ExpectedData: `&validator_test.FakeDataStruct{StringField:null.Nullable[string]{V:"valid", Valid:true}, FloatField:null.Nullable[float64]{V:-5, Valid:true}, IntField:null.Nullable[int]{V:42, Valid:true}}`,
			ExpectedErrors: map[string][]string{
				"string_field": {},
				"float_field":  {"this value should be higher or equal to 10"},
				"int_field":    {},
			},
		},
		"Negative int field": {
			In: map[string]interface{}{
				"string_field": "valid",
				"float_field":  12.5,
				"int_field":    -1,
			},
			BindTo:       &FakeDataStruct{},
			ExpectedData: `&validator_test.FakeDataStruct{StringField:null.Nullable[string]{V:"valid", Valid:true}, FloatField:null.Nullable[float64]{V:12.5, Valid:true}, IntField:null.Nullable[int]{V:-1, Valid:true}}`,
			ExpectedErrors: map[string][]string{
				"string_field": {},
				"float_field":  {},
				"int_field":    {"this value should be higher or equal to 0"},
			},
		},
		"Int instead of a string": {
			In: map[string]interface{}{
				"string_field": 10,
				"float_field":  12.5,
				"int_field":    1,
			},
			BindTo:       &FakeDataStruct{},
			ExpectedData: `&validator_test.FakeDataStruct{StringField:null.Nullable[string]{V:"", Valid:true}, FloatField:null.Nullable[float64]{V:12.5, Valid:true}, IntField:null.Nullable[int]{V:1, Valid:true}}`,
			ExpectedErrors: map[string][]string{
				"string_field": {"this is not a valid value. Expect a string"},
				"int_field":    {},
				"float_field":  {},
			},
		},
		"bool instead of a string": {
			In: map[string]interface{}{
				"string_field": true,
				"float_field":  12.5,
				"int_field":    1,
			},
			BindTo:       &FakeDataStruct{},
			ExpectedData: `&validator_test.FakeDataStruct{StringField:null.Nullable[string]{V:"", Valid:true}, FloatField:null.Nullable[float64]{V:12.5, Valid:true}, IntField:null.Nullable[int]{V:1, Valid:true}}`,
			ExpectedErrors: map[string][]string{
				"string_field": {"this is not a valid value. Expect a string"},
				"int_field":    {},
				"float_field":  {},
			},
		},
		"slice instead of a string": {
			In: map[string]interface{}{
				"string_field": []any{"this", "is", 1},
				"float_field":  12.5,
				"int_field":    1,
			},
			BindTo:       &FakeDataStruct{},
			ExpectedData: `&validator_test.FakeDataStruct{StringField:null.Nullable[string]{V:"", Valid:false}, FloatField:null.Nullable[float64]{V:12.5, Valid:true}, IntField:null.Nullable[int]{V:1, Valid:true}}`,
			ExpectedErrors: map[string][]string{
				"string_field": {"this is not a valid value"},
				"int_field":    {},
				"float_field":  {},
			},
		},
		"valid float String instead of a float": {
			In: map[string]interface{}{
				"string_field": "",
				"float_field":  "12.5",
				"int_field":    1,
			},
			BindTo:       &FakeDataStruct{},
			ExpectedData: `&validator_test.FakeDataStruct{StringField:null.Nullable[string]{V:"", Valid:true}, FloatField:null.Nullable[float64]{V:12.5, Valid:true}, IntField:null.Nullable[int]{V:1, Valid:true}}`,
			ExpectedErrors: map[string][]string{
				"string_field": {},
				"int_field":    {},
				"float_field":  {},
			},
		},
		"invalid float string instead of a float": {
			In: map[string]interface{}{
				"string_field": "",
				"float_field":  "azerty",
				"int_field":    1,
			},
			BindTo:       &FakeDataStruct{},
			ExpectedData: `&validator_test.FakeDataStruct{StringField:null.Nullable[string]{V:"", Valid:true}, FloatField:null.Nullable[float64]{V:0, Valid:true}, IntField:null.Nullable[int]{V:1, Valid:true}}`,
			ExpectedErrors: map[string][]string{
				"string_field": {},
				"int_field":    {},
				"float_field":  {"cannot cast 'azerty':string to float"},
			},
		},
		"valid int String instead of a int": {
			In: map[string]interface{}{
				"string_field": "",
				"float_field":  12.5,
				"int_field":    "1",
			},
			BindTo:       &FakeDataStruct{},
			ExpectedData: `&validator_test.FakeDataStruct{StringField:null.Nullable[string]{V:"", Valid:true}, FloatField:null.Nullable[float64]{V:12.5, Valid:true}, IntField:null.Nullable[int]{V:1, Valid:true}}`,
			ExpectedErrors: map[string][]string{
				"string_field": {},
				"int_field":    {},
				"float_field":  {},
			},
		},
		"invalid int string instead of a int": {
			In: map[string]interface{}{
				"string_field": "",
				"float_field":  "12.5",
				"int_field":    "azerty",
			},
			BindTo:       &FakeDataStruct{},
			ExpectedData: `&validator_test.FakeDataStruct{StringField:null.Nullable[string]{V:"", Valid:true}, FloatField:null.Nullable[float64]{V:12.5, Valid:true}, IntField:null.Nullable[int]{V:0, Valid:true}}`,
			ExpectedErrors: map[string][]string{
				"string_field": {},
				"int_field":    {"cannot cast 'azerty':string to int"},
				"float_field":  {},
			},
		},
		"Empty func struct": {
			In:           map[string]interface{}{},
			BindTo:       &FakeFuncDataStruct{},
			ExpectedData: `&validator_test.FakeFuncDataStruct{Field3:"", Field4:""}`,
			ExpectedErrors: map[string][]string{
				"field3": {"this value should not be blank"},
				"field4": {},
			},
		},
		"invalid func struct": {
			In: map[string]interface{}{
				"field3": "12345678901234567890",
				"field4": "12345678901234567890",
			},
			BindTo:       &FakeFuncDataStruct{},
			ExpectedData: `&validator_test.FakeFuncDataStruct{Field3:"12345678901234567890", Field4:"12345678901234567890"}`,
			ExpectedErrors: map[string][]string{
				"field3": {"higher than max length of 15"},
				"field4": {"higher than max length of 15"},
			},
		},
		"Empty func struct does not erase previous value, but the value is still required een if a value already exist": {
			In:           map[string]interface{}{},
			BindTo:       &FakeFuncDataStruct{},
			ExpectedData: `&validator_test.FakeFuncDataStruct{Field3:"", Field4:""}`,
			ExpectedErrors: map[string][]string{
				"field3": {"this value should not be blank"},
				"field4": {},
			},
		},
		"can send float or int to float or int": {
			In: map[string]interface{}{
				"string_field": "",
				"float_field":  10,
				"int_field":    10.0,
			},
			BindTo:       &FakeDataStruct{},
			ExpectedData: `&validator_test.FakeDataStruct{StringField:null.Nullable[string]{V:"", Valid:true}, FloatField:null.Nullable[float64]{V:10, Valid:true}, IntField:null.Nullable[int]{V:10, Valid:true}}`,
			ExpectedErrors: map[string][]string{
				"string_field": {},
				"int_field":    {},
				"float_field":  {},
			},
		},
		"can fetch data from array": {
			In: map[string]interface{}{
				"array_field": []map[string]interface{}{
					{"sub_field": "field_0"},
					{"sub_field": "field_1"},
					{"sub_field": "field_2"},
				},
			},
			BindTo:       &FakeDataWithArrayStruct{},
			ExpectedData: `&validator_test.FakeDataWithArrayStruct{Array:[]validator_test.SubDataStruct{validator_test.SubDataStruct{SubStringField:null.Nullable[string]{V:"field_0", Valid:true}}, validator_test.SubDataStruct{SubStringField:null.Nullable[string]{V:"field_1", Valid:true}}, validator_test.SubDataStruct{SubStringField:null.Nullable[string]{V:"field_2", Valid:true}}}}`,
			ExpectedErrors: map[string][]string{
				"array_field.0.sub_field": {},
				"array_field.1.sub_field": {},
				"array_field.2.sub_field": {},
			},
		},
		"Valid nested map": {
			In: map[string]interface{}{
				"data": map[string]any{
					"string_field": "ok",
					"float_field":  12.34,
					"int_field":    99,
				},
			},
			BindTo:       &FakeMapStruct{},
			ExpectedData: `&validator_test.FakeMapStruct{Data:map[string]interface {}{"float_field":12.34, "int_field":99, "string_field":"ok"}}`,
			ExpectedErrors: map[string][]string{
				"data": {},
			},
		},
		"Map inside a map (extract only subdata)": {
			In: map[string]interface{}{

				"data": map[string]any{
					"subdata": map[string]any{
						"subfield1": "value1",
						"subfield2": 42,
					},
				},
			},
			BindTo:       &FakeMapWithSubMapStruct{},
			ExpectedData: `&validator_test.FakeMapWithSubMapStruct{Subdata:validator_test.FakeMapStruct{Data:map[string]interface {}{"subfield1":"value1", "subfield2":42}}}`,
			ExpectedErrors: map[string][]string{
				"data.subdata": {},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			bm := validator.BinderMap{Data: maputil.Flattern(tt.In, ".")}
			tt.BindTo.Bind(&bm)

			resultData := fmt.Sprintf("%#v", tt.BindTo)

			if resultData != tt.ExpectedData {
				t.Errorf("got data \n%v\nwant \n%v", resultData, tt.ExpectedData)
			}

			if !reflect.DeepEqual(bm.Errors.Errors, tt.ExpectedErrors) {
				t.Logf("data : %v", resultData)
				t.Fatalf("got errors \n%#v\nwant \n%#v", bm.Errors.Errors, tt.ExpectedErrors)
			}
		})
	}
}
