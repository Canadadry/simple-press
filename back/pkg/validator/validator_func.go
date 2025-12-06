package validator

import (
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	DateTimeFormat = "2006-01-02T15:04:05Z07:00"
	EmailRe        = "^[a-zA-Z0-9_!#$%&’*+/=?`{|}~^.-]+@[a-zA-Z0-9.-]+$"
	UuidRe         = "^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$"
	PhoneRe        = "^(?:(?:\\+|00)33[\\s.-]{0,3}(?:\\(0\\)[\\s.-]{0,3})?|0)[1-9](?:(?:[\\s.-]?\\d{2}){4}|\\d{2}(?:[\\s.-]?\\d{3}){2})$"
)

func IsDatetime(layout string) func(val string) error {
	return func(val string) error {
		_, err := time.Parse(layout, val)
		if err != nil {
			return fmt.Errorf("invalid date format want \"%s\" : %w", layout, err)
		}
		return nil
	}
}

func Length(min, max int) func(val string) error {
	return func(val string) error {
		length := utf8.RuneCountInString(val)
		if length > max {
			return fmt.Errorf("higher than max length of %d", max)
		}
		if length < min {
			return fmt.Errorf("lowr than min length of %d", min)
		}
		return nil
	}
}

func Choice[T comparable](choices []T) func(val T) error {
	return func(val T) error {
		for _, c := range choices {
			if c == val {
				return nil
			}
		}
		return fmt.Errorf("the value you selected is not a valid choice")
	}
}

type Number interface {
	int | int64 | float64
}

type NullableNumber[T Number] struct {
	Number T
	Valid  bool
}

func Min[T Number](min T) func(val T) error {
	return inRange[T](validNumber(min), NullableNumber[T]{})
}

func Max[T Number](max T) func(val T) error {
	return inRange[T](NullableNumber[T]{}, validNumber(max))
}

func Range[T Number](min, max T) func(val T) error {
	return inRange[T](validNumber(min), validNumber(max))
}

func validNumber[T Number](val T) NullableNumber[T] {
	return NullableNumber[T]{
		Number: val,
		Valid:  true,
	}
}

func inRange[T Number](min, max NullableNumber[T]) func(val T) error {
	return func(val T) error {
		if val > max.Number && max.Valid {
			return fmt.Errorf("this value should be lower or equal to %v", max.Number)
		}
		if val < min.Number && min.Valid {
			return fmt.Errorf("this value should be higher or equal to %v", min.Number)
		}
		return nil
	}
}

func Integer(val float64) error {
	if val != float64(int64(val)) {
		return fmt.Errorf("this value is not an integer")
	}
	return nil
}

func Exist[T any](exist func(val T) bool) func(val T) error {
	return func(val T) error {
		if !exist(val) {
			return fmt.Errorf("this value does not exist")
		}
		return nil
	}
}

func ExistWithCheck[T, U any](exist func(val T) (U, bool), check func(u U) error) func(val T) error {
	return func(val T) error {
		u, ok := exist(val)
		if !ok {
			return fmt.Errorf("this value does not exist")
		}
		if err := check(u); err != nil {
			return err
		}
		return nil
	}
}

func IsUnique[T any](find func(val T) bool) func(val T) error {
	return func(val T) error {
		if find(val) {
			return fmt.Errorf("this value is already used")
		}
		return nil
	}
}

func EmailMX(val string) error {
	part := strings.Split(val, "@")
	domain := part[len(part)-1]
	mxrecords, err := net.LookupMX(domain)
	if err != nil || len(mxrecords) == 0 {
		return fmt.Errorf("mx record dont exist on %s", domain)
	}
	return nil
}

func Regexp(pattern string) func(val string) error {
	re, err := regexp.Compile(pattern)
	return func(val string) error {
		if err != nil {
			return err
		}
		if !re.MatchString(val) {
			return fmt.Errorf("value dont match format")
		}
		return nil
	}
}

func OnlyOneError[T any](fns ...func(val T) error) func(T) error {
	return func(val T) error {
		for _, fn := range fns {
			if err := fn(val); err != nil {
				return err
			}
		}
		return nil
	}
}

func ReplaceErrorWith[T any](newErr error, fn func(T) error) func(T) error {
	return func(val T) error {
		if err := fn(val); err != nil {
			return newErr
		}
		return nil
	}
}
