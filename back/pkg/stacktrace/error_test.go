package stacktrace

import (
	"fmt"
	"strings"
	"testing"
)

func addFrameA(err error) error {
	return addFrameB(err)
}

func addFrameB(err error) error {
	return addFrameC(err)
}

func addFrameC(err error) error {
	return From(err)
}

func deepErrorA(err error) error {
	return deepErrorB(err)
}

func deepErrorB(err error) error {
	return deepErrorC(err)
}

func deepErrorC(err error) error {
	return err
}

func TestError(t *testing.T) {
	tests := map[string]struct {
		Error           *Error
		ExpectedMessage string
		ExpectedTrace   string
	}{
		"basic": {
			Error:           addFrameA(fmt.Errorf("error with stack trace")).(*Error),
			ExpectedMessage: "error with stack trace",
			ExpectedTrace: `error with stack trace
pkg/stacktrace/error_test.go:18 app/pkg/stacktrace.addFrameC()
pkg/stacktrace/error_test.go:14 app/pkg/stacktrace.addFrameB()
pkg/stacktrace/error_test.go:10 app/pkg/stacktrace.addFrameA()
pkg/stacktrace/error_test.go:40 app/pkg/stacktrace.TestError()
testing.tRunner()
runtime.goexit()`,
		},
		"basic traced twice": {
			Error:           From(addFrameA(fmt.Errorf("error traced twice"))),
			ExpectedMessage: "error traced twice",
			ExpectedTrace: `error traced twice
pkg/stacktrace/error_test.go:18 app/pkg/stacktrace.addFrameC()
pkg/stacktrace/error_test.go:14 app/pkg/stacktrace.addFrameB()
pkg/stacktrace/error_test.go:10 app/pkg/stacktrace.addFrameA()
pkg/stacktrace/error_test.go:51 app/pkg/stacktrace.TestError()
testing.tRunner()
runtime.goexit()`,
		},
		"basic traced lately": {
			Error:           From(deepErrorA(fmt.Errorf("error traced lately"))),
			ExpectedMessage: "error traced lately",
			ExpectedTrace: `error traced lately
pkg/stacktrace/error_test.go:62 app/pkg/stacktrace.TestError()
testing.tRunner()
runtime.goexit()`,
		},
		"wrapped after traced": {
			Error:           Errorf("wrap : %w", addFrameA(fmt.Errorf("error wrapped after traced"))),
			ExpectedMessage: "wrap : error wrapped after traced",
			ExpectedTrace: `wrap : error wrapped after traced
pkg/stacktrace/error_test.go:18 app/pkg/stacktrace.addFrameC()
pkg/stacktrace/error_test.go:14 app/pkg/stacktrace.addFrameB()
pkg/stacktrace/error_test.go:10 app/pkg/stacktrace.addFrameA()
pkg/stacktrace/error_test.go:70 app/pkg/stacktrace.TestError()
testing.tRunner()
runtime.goexit()`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if tt.Error.Error() != tt.ExpectedMessage {
				t.Fatalf("error exp '%s' got '%s'", tt.ExpectedMessage, tt.Error.Error())
			}
			got := strings.Split(tt.Error.Trace(), "\n")
			exp := strings.Split(tt.ExpectedTrace, "\n")
			for i := 1; i < len(exp); i++ {
				if !strings.HasSuffix(got[i], exp[i]) {
					t.Fatalf("at [%d] trace exp \n%s\n got \n%s\n", i, got[i], exp[i])
				}
			}
		})
	}
}
