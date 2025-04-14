package trycatch

import (
	"errors"
	"testing"
)

func TestCatch(t *testing.T) {
	err1 := errors.New("fake")
	tests := map[string]struct {
		Func     func() error
		Expected error
	}{
		"basic return nil": {
			Func:     func() error { return nil },
			Expected: nil,
		},
		"basic return err": {
			Func:     func() error { return err1 },
			Expected: err1,
		},
		"throw": {
			Func: func() error {
				var mySlice []int
				_ = mySlice[0]
				return nil
			},
			Expected: ThrowErr,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := Catch(tt.Func)
			if !errors.Is(result, tt.Expected) {
				t.Fatalf("exp %v got %v", tt.Expected, result)
			}
		})
	}
}
