package validator

import (
	"app/pkg/maputil"
	"app/pkg/null"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"time"
)

var (
	ErrNotBlank     = fmt.Errorf("this value should not be blank")
	ErrInvalidValue = fmt.Errorf("this is not a valid value")
)

const (
	defaultMaxMemory                    = 10 << 20 // 10 MB
	ApplicationFormContentType          = "application/x-www-form-urlencoded"
	ApplicationMultiPartFormContentType = "multipart/form-data"
	ApplicationJsonContentType          = "application/json"
	ContentTypeHeaderName               = "Content-Type"
)

func BindWithJson(r io.ReadCloser, bind func(Binder)) (Errors, error) {
	defer r.Close()
	body := map[string]any{}
	err := json.NewDecoder(r).Decode(&body)
	if err != nil && !errors.Is(err, io.EOF) {
		return Errors{}, fmt.Errorf("while decoding body as json : %w", err)
	}
	return BindWithMap(body, bind), nil
}

func BindWithForm(r *http.Request, bind func(Binder)) (Errors, error) {
	defer r.Body.Close()

	body := map[string]any{}

	ct := r.Header.Get(ContentTypeHeaderName)
	switch ct {
	case ApplicationJsonContentType:
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil && !errors.Is(err, io.EOF) {
			return Errors{}, fmt.Errorf("while decoding body as json : %w", err)
		}
	case ApplicationFormContentType:
		err := r.ParseForm()
		if err != nil {
			return Errors{}, err
		}
		for name, value := range r.Form {
			if len(value) > 0 {
				body[name] = value[0]
			} else {
				body[name] = ""
			}
		}
	case ApplicationMultiPartFormContentType:
		err := r.ParseMultipartForm(defaultMaxMemory)
		if err != nil {
			return Errors{}, err
		}
		for name, value := range r.MultipartForm.Value {
			if len(value) > 0 {
				body[name] = value[0]
			} else {
				body[name] = ""
			}
		}
	default:
		return Errors{}, fmt.Errorf("invalid content type '%s'", ct)
	}
	return BindWithMap(body, bind), nil
}

func BindWithMap(body map[string]any, bind func(Binder)) Errors {
	b := &BinderMap{
		Data: maputil.Flattern(body, "."),
	}
	bind(b)
	return b.Errors
}

type BinderMap struct {
	Data   map[string]any
	prefix []string
	Errors Errors
}

func (b *BinderMap) pushPrefix(p string) {
	b.prefix = append(b.prefix, p)
}

func (b *BinderMap) popPrefix() {
	b.prefix = b.prefix[:len(b.prefix)-1]
}

func (b *BinderMap) fullKey(name string) string {
	if len(b.prefix) == 0 {
		return name
	}
	return strings.Join(append(b.prefix, name), ".")
}

func (bm *BinderMap) RequiredStringVar(name string, ptr *string, fns ...func(string) error) {
	bm.stringVar(name, true, ptr, fns...)
}

func (bm *BinderMap) RequiredFloat64Var(name string, ptr *float64, fns ...func(float64) error) {
	bm.float64Var(name, true, ptr, fns...)
}

func (bm *BinderMap) RequiredIntVar(name string, ptr *int, fns ...func(int) error) {
	bm.intVar(name, true, ptr, fns...)
}

func (bm *BinderMap) RequiredTimeVar(name string, ptr *time.Time, format string, fns ...func(time.Time) error) {
	bm.timeVar(name, true, ptr, format, fns...)
}

func (bm *BinderMap) RequiredBoolVar(name string, ptr *bool, trueChoice []string, falseChoice []string) {
	bm.boolVar(name, true, ptr, trueChoice, falseChoice)
}

func (bm *BinderMap) StringVar(name string, ptr *null.Nullable[string], fns ...func(string) error) {
	ptr.Valid = bm.stringVar(name, false, &ptr.V, fns...)
}

func (bm *BinderMap) Float64Var(name string, ptr *null.Nullable[float64], fns ...func(float64) error) {
	ptr.Valid = bm.float64Var(name, false, &ptr.V, fns...)
}

func (bm *BinderMap) IntVar(name string, ptr *null.Nullable[int], fns ...func(int) error) {
	ptr.Valid = bm.intVar(name, false, &ptr.V, fns...)
}

func (bm *BinderMap) TimeVar(name string, ptr *null.Nullable[time.Time], format string, fns ...func(time.Time) error) {
	ptr.Valid = bm.timeVar(name, false, &ptr.V, format, fns...)
	if ptr.Valid && ptr.V.IsZero() {
		ptr.Valid = false
	}
}

func (bm *BinderMap) BoolVar(name string, ptr *null.Nullable[bool], trueChoice []string, falseChoice []string) {
	ptr.Valid = bm.boolVar(name, false, &ptr.V, trueChoice, falseChoice)
}

func (bm *BinderMap) stringVar(name string, requiered bool, ptr *string, fns ...func(string) error) bool {
	return bm.anyFunc(name, requiered, ParseAndValidate[any, string](ptr, ParseToString, fns...))
}

func (bm *BinderMap) float64Var(name string, requiered bool, ptr *float64, fns ...func(float64) error) bool {
	return bm.anyFunc(name, requiered, ParseAndValidate[any, float64](ptr, ParseToFloat64, fns...))
}

func (bm *BinderMap) intVar(name string, requiered bool, ptr *int, fns ...func(int) error) bool {
	return bm.anyFunc(name, requiered, ParseAndValidate[any, int](ptr, ParseToInt, fns...))
}

func (bm *BinderMap) timeVar(name string, requiered bool, ptr *time.Time, format string, fns ...func(time.Time) error) bool {
	return bm.anyFunc(name, requiered, ParseAndValidate[any, time.Time](ptr, ParseToTime(format), fns...))
}

func (bm *BinderMap) boolVar(name string, requiered bool, ptr *bool, trueChoice, falseChoice []string) bool {
	return bm.anyFunc(name, requiered, ParseAndValidate[any, bool](ptr, ParseToBool(trueChoice, falseChoice)))
}

func (bm *BinderMap) RequiredAnyFunc(name string, fn func(string, bool, any, *Errors)) {
	bm.anyFunc(name, true, fn)
}

func (bm *BinderMap) AnyFunc(name string, fn func(string, bool, any, *Errors)) {
	bm.anyFunc(name, false, fn)
}

func (bm *BinderMap) anyFunc(name string, requiered bool, fn func(string, bool, any, *Errors)) bool {
	name = bm.fullKey(name)
	bm.Errors.Accumulate(name, nil)
	val, ok := bm.Data[name]
	if ok {
		fn(name, ok, val, &bm.Errors)
		return true
	}
	if bm.HasPrefix(name) {
		bm.Errors.Accumulate(name, ErrInvalidValue)
	}
	if requiered {
		bm.Errors.Accumulate(name, ErrNotBlank)
	}
	return false
}

func (bm *BinderMap) RequiredArrayVar(name string, ptr any, create func() Bindable) {
	bm.ArrayVar(name, ptr, create)
	if reflect.ValueOf(ptr).Elem().Len() == 0 {
		bm.Errors.Accumulate(name, ErrNotBlank)
	}
}

func (bm *BinderMap) ArrayVar(name string, ptr any, create func() Bindable) {
	slicePtrVal := reflect.ValueOf(ptr)
	if slicePtrVal.Kind() != reflect.Ptr || slicePtrVal.Elem().Kind() != reflect.Slice {
		panic("RequiredArrayVar: ptr must be a pointer to a slice")
	}

	sliceVal := slicePtrVal.Elem()
	elemType := sliceVal.Type().Elem()

	for i := 0; bm.HasPrefix(fmt.Sprintf("%s.%d", name, i)); i++ {
		bindable := create()

		bm.pushPrefix(fmt.Sprintf("%s.%d", name, i))
		bindable.Bind(bm)
		bm.popPrefix()

		bindableVal := reflect.ValueOf(bindable)

		if bindableVal.Type().AssignableTo(elemType) {
		} else if bindableVal.Kind() == reflect.Ptr && bindableVal.Elem().Type().AssignableTo(elemType) {
			bindableVal = bindableVal.Elem()
		} else {
			panic(fmt.Sprintf(
				"RequiredArrayVar: bindable of type %T is not assignable to slice element type %s",
				bindable, elemType,
			))
		}

		sliceVal.Set(reflect.Append(sliceVal, bindableVal))
	}
}

func (bm *BinderMap) RequiredMapVar(name string, ptr *map[string]any) {
	name = bm.fullKey(name)
	bm.MapVar(name, ptr)
	if len(*ptr) == 0 {
		bm.Errors.Accumulate(name, ErrNotBlank)
	}
}

func (bm *BinderMap) MapVar(name string, ptr *map[string]any) {
	prefix := name + "."
	data := map[string]any{}

	for k, v := range bm.Data {
		if strings.HasPrefix(k, prefix) {
			trimmedKey := strings.TrimPrefix(k, prefix)
			data[trimmedKey] = v
		}
	}
	*ptr = maputil.Expand(data, ".")

	bm.Errors.Accumulate(name, nil)
}

func (bm *BinderMap) HasPrefix(name string) bool {
	for k := range bm.Data {
		if strings.HasPrefix(k, name) {
			return true
		}
	}
	return false
}
