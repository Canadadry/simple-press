package validator

import (
	"app/pkg/null"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	TrueChoice  = []string{"true", "1", "on", "yes"}
	FalseChoice = []string{"false", "0", "off", "no"}
)

type Bindable interface {
	Bind(b Binder)
}

func NewBindable[T any]() Bindable {
	var zero T
	ptr := any(&zero)
	bindable, ok := ptr.(Bindable)
	if !ok {
		panic(fmt.Sprintf("*%T does not implement Bindable", zero))
	}
	return bindable
}

type Binder interface {
	RequiredStringVar(name string, ptr *string, fns ...func(string) error)
	RequiredFloat64Var(name string, ptr *float64, fns ...func(float64) error)
	RequiredIntVar(name string, ptr *int, fns ...func(int) error)
	RequiredTimeVar(name string, ptr *time.Time, format string, fns ...func(time.Time) error)
	RequiredBoolVar(name string, ptr *bool, trueChoice []string, falseChoice []string)
	RequiredAnyFunc(name string, fn func(string, bool, any, *Errors))
	RequiredArrayVar(name string, ptr any, create func() Bindable)
	RequiredMapVar(name string, ptr *map[string]any)
	StringVar(name string, ptr *null.Nullable[string], fns ...func(string) error)
	Float64Var(name string, ptr *null.Nullable[float64], fns ...func(float64) error)
	IntVar(name string, ptr *null.Nullable[int], fns ...func(int) error)
	TimeVar(name string, ptr *null.Nullable[time.Time], format string, fns ...func(time.Time) error)
	BoolVar(name string, ptr *null.Nullable[bool], trueChoice []string, falseChoice []string)
	AnyFunc(name string, fn func(string, bool, any, *Errors))
	ArrayVar(name string, ptr any, create func() Bindable)
	MapVar(name string, ptr *map[string]any)
}

func ParseAndValidate[From, To any](t *To, parse func(From) (To, error), fns ...func(To) error) func(string, bool, From, *Errors) {
	return func(name string, ok bool, val From, errs *Errors) {
		if ok {
			parsed, err := parse(val)
			if err != nil {
				errs.Accumulate(name, err)
				return
			}
			*t = parsed
		}
		Validate[To](errs, name, *t, fns...)
	}
}

func Validate[T any](errs *Errors, name string, value T, fns ...func(T) error) {
	for _, fn := range fns {
		if err := fn(value); err != nil {
			errs.Accumulate(name, err)
		}
	}
}

func ParseToString(val any) (string, error) {
	switch v := val.(type) {
	case string:
		return v, nil
	}
	return "", fmt.Errorf("this is not a valid value. Expect a string")
}

func ParseToTime(format string) func(val any) (time.Time, error) {
	return func(val any) (time.Time, error) {
		if val == nil {
			return time.Time{}, nil
		}
		v, err := ParseToString(val)
		if err != nil {
			return time.Unix(0, 0), err
		}
		return time.Parse(DateTimeFormat, v)
	}
}

func ParseToBool(trueChoice, falseChoice []string) func(val any) (bool, error) {
	return func(val any) (bool, error) {
		v := fmt.Sprintf("%v", val)
		isTrue := Choice(trueChoice)
		isFalse := Choice(falseChoice)
		if err := isTrue(strings.ToLower(v)); err == nil {
			return true, nil
		}
		if err := isFalse(strings.ToLower(v)); err == nil {
			return false, nil
		}
		return false, fmt.Errorf("invalid bool value")
	}
}

func ParseToFloat64(val any) (float64, error) {
	switch v := val.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case string:
		parsed, err := strconv.ParseFloat(v, 64)
		if err == nil {
			return parsed, nil
		}
	}
	return 0.0, fmt.Errorf("cannot cast '%v':%T to float", val, val)
}

func ParseToInt(val any) (int, error) {
	switch v := val.(type) {
	case int:
		return v, nil
	case float64:
		return int(v), nil
	case string:
		parsed, err := strconv.Atoi(v)
		if err == nil {
			return parsed, nil
		}
	}
	return 0, fmt.Errorf("cannot cast '%v':%T to int", val, val)
}

var ParseToTimeDefault = ParseToTime(DateTimeFormat)
