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

func (r *MockWriterRenderer) BeginForm(name, action, method string) {
	r.logf("BeginForm name=%s action=%s method=%s", name, action, method)
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
func (r *MockWriterRenderer) Input(label, name, inputType string) {
	r.logf("Input label=%s name=%s type=%s", label, name, inputType)
}
func (r *MockWriterRenderer) Checkbox(label, name string) {
	r.logf("Checkbox label=%s name=%s", label, name)
}
func (r *MockWriterRenderer) Select(label, name string, options []string) {
	r.logf("Select label=%s name=%s options=%v", label, name, options)
}
func (r *MockWriterRenderer) BeginArray(name, path string) {
	r.logf("BeginArray name=%s path=%s", name, path)
}
func (r *MockWriterRenderer) EndArray() {
	r.logf("EndArray")
}
func (r *MockWriterRenderer) BeginArrayItem(index int) {
	r.logf("BeginArrayItem index=%d", index)
}
func (r *MockWriterRenderer) EndArrayItem() {
	r.logf("EndArrayItem")
}

func TestGenerateFormDynamicHTMLWithName(t *testing.T) {
	tests := map[string]struct {
		Input       map[string]any
		FormName    string
		ExpectedLog []string
	}{
		"simple dynamic array": {
			Input: map[string]any{
				"children": []any{
					map[string]any{
						"firstname": "string",
						"gender":    "enum:Mr;Mme",
					},
				},
			},
			FormName: "family",
			ExpectedLog: []string{
				"BeginForm name=family action=/submit method=POST",
				"BeginFieldset label=children",
				"BeginArray name=children path=children",
				"BeginArrayItem index=0",
				"Input label=firstname name=children.0.firstname type=text",
				"Select label=gender name=children.0.gender options=[Mr Mme]",
				"EndArrayItem",
				"EndArray",
				"EndFieldset",
				"EndForm",
			},
		},
		"nested dynamic object with array": {
			Input: map[string]any{
				"profile": map[string]any{
					"firstname": "string",
					"emails":    []any{"email"},
				},
			},
			FormName: "nested",
			ExpectedLog: []string{
				"BeginForm name=nested action=/submit method=POST",
				"BeginFieldset label=profile",
				"BeginFieldset label=emails",
				"BeginArray name=emails path=profile.emails",
				"BeginArrayItem index=0",
				"Input label= name=profile.emails.0 type=email",
				"EndArrayItem",
				"EndArray",
				"EndFieldset",
				"Input label=firstname name=profile.firstname type=text",
				"EndFieldset",
				"EndForm",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			root := Parse(tt.Input, true)
			var buf strings.Builder
			mock := &MockWriterRenderer{w: &buf}
			GenerateFormDynamicHTMLWithName(root, mock, tt.FormName)
			lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
			compareCallStacks(t, lines, tt.ExpectedLog)
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
