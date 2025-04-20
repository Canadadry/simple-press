package data

import (
	"fmt"
	"io"
	"strings"
	"testing"
)

type MockWriterRenderer struct {
	w io.Writer
}

func (r *MockWriterRenderer) logf(format string, args ...any) {
	fmt.Fprintf(r.w, format+"\n", args...)
}

func (r *MockWriterRenderer) BeginForm() {
	r.logf("BeginForm")
}
func (r *MockWriterRenderer) EndForm() {
	r.logf("EndForm")
}
func (r *MockWriterRenderer) BeginFieldset(label string) {
	r.logf("BeginFieldset label=%s", label)
}
func (r *MockWriterRenderer) EndFieldset() {
	r.logf("EndFieldset")
}
func (r *MockWriterRenderer) Input(label, name, inputType, value string) {
	r.logf("Input label=%s name=%s type=%s value=%s", label, name, inputType, value)
}
func (r *MockWriterRenderer) Checkbox(label, name string, checked bool) {
	r.logf("Checkbox label=%s name=%s checked=%t", label, name, checked)
}
func (r *MockWriterRenderer) Select(label, name string, options []string, value string) {
	r.logf("Select label=%s name=%s options=%v value=%s", label, name, options, value)
}

func TestRenderDefinition(t *testing.T) {
	tests := map[string]struct {
		Input    any
		Expected []string
	}{

		"nested object": {
			Input: map[string]any{
				"profile": map[string]any{
					"name": map[string]any{
						"first": "Jane",
						"last":  "Doe",
					},
					"age": 42,
				},
			},
			Expected: []string{
				"BeginForm name=form action=/submit method=POST",
				"BeginFieldset label=profile",
				"Input label=age name=profile.age type=number value=42",
				"BeginFieldset label=name",
				"Input label=first name=profile.name.first type=text value=Jane",
				"Input label=last name=profile.name.last type=text value=Doe",
				"EndFieldset",
				"EndFieldset",
				"EndForm Submit",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var buf strings.Builder
			mock := &MockWriterRenderer{w: &buf}
			Render(tt.Input, mock)
			lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
			compareCallStacks(t, lines, tt.Expected)
		})
	}
}

func compareCallStacks(t *testing.T, got, want []string) {
	t.Helper()

	if len(got) != len(want) {
		t.Logf("Length mismatch: got %d calls, want %d calls", len(got), len(want))
	}

	max := len(got)
	if len(want) > max {
		max = len(want)
	}

	hasDiff := false

	for i := 0; i < max; i++ {
		var g, w string
		if i < len(got) {
			g = strings.TrimSpace(got[i])
		}
		if i < len(want) {
			w = strings.TrimSpace(want[i])
		}

		status := "✓"
		if g != w {
			status = "✗"
			hasDiff = true
		}

		t.Logf("[%s] [%2d] want: %-60s got: %s", status, i, w, g)
	}

	if hasDiff {
		t.FailNow()
	}
}
