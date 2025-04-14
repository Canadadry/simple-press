package sflag

import (
	"app/pkg/structtag"
	"fmt"
	"reflect"
	"strconv"
)

const (
	tagFlagName         = "flag"
	tagFlagDefaultValue = "default"
	tagFlagDesc         = "desc"
)

type FlagSet interface {
	StringVar(ptr *string, name string, value string, desc string)
	IntVar(ptr *int, name string, value int, desc string)
	BoolVar(ptr *bool, name string, value bool, desc string)
}

func Parse(fs FlagSet, s interface{}) error {
	val := reflect.ValueOf(s)

	if val.Kind() != reflect.Ptr {
		return fmt.Errorf("should pass pointer of struct")
	}
	val = val.Elem()

	if val.Kind() != reflect.Struct {
		return fmt.Errorf("should pass pointer of struct")
	}
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		err := parseField(fs, valueField, typeField)
		if err != nil {
			return fmt.Errorf("while parsing %s : %w", typeField.Name, err)
		}
	}
	return nil
}

func parseField(fs FlagSet, valueField reflect.Value, typeField reflect.StructField) error {
	extractedTag := string(typeField.Tag)
	tag, err := structtag.Parse(extractedTag)
	if err != nil {
		return fmt.Errorf("cannot parse tag : %w", err)
	}
	flagName, fOk := structtag.Get(tag, tagFlagName)

	if (!fOk || flagName == "-") && valueField.Kind() != reflect.Struct {
		return nil
	}

	flagDefaultValue, flagHasDefaultValue := structtag.Get(tag, tagFlagDefaultValue)
	flagDesc, _ := structtag.Get(tag, tagFlagDesc)

	switch valueField.Kind() {
	case reflect.String:
		fs.StringVar(valueField.Addr().Interface().(*string), flagName, flagDefaultValue, flagDesc)
	case reflect.Int:
		defaultInt := 0
		if flagHasDefaultValue {
			i, err := strconv.ParseInt(flagDefaultValue, 10, 64)
			if err != nil {
				return fmt.Errorf("cannot parse default int value : %w", err)
			}
			defaultInt = int(i)
		}
		fs.IntVar(valueField.Addr().Interface().(*int), flagName, defaultInt, flagDesc)
	case reflect.Bool:
		defaultBool := false
		if flagHasDefaultValue {
			b, err := strconv.ParseBool(flagDefaultValue)
			if err != nil {
				return fmt.Errorf("cannot parse default int value : %w", err)
			}
			defaultBool = b
		}
		fs.BoolVar(valueField.Addr().Interface().(*bool), flagName, defaultBool, flagDesc)
	case reflect.Struct:
		err := Parse(fs, valueField.Addr().Interface())
		if err != nil {
			return fmt.Errorf("cannot parse field : %w", err)
		}
	default:
		return fmt.Errorf("cannont parse type %v", valueField.Kind())
	}

	return nil

}
