package form

import (
	"app/pkg/router"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

const (
	blockEditName       = "name"
	blockEditContent    = "content"
	blockEditDefinition = "definition"
)

type BlockEdit struct {
	Name       string
	Content    string
	Definition map[string]any
}

type BlockEditError struct {
	Name       string
	Content    string
	Definition string
}

func (le BlockEditError) HasError() bool {
	if le.Name != "" {
		return true
	}
	if le.Content != "" {
		return true
	}
	if le.Definition != "" {
		return true
	}
	return false
}

func ParseBlockEdit(r *http.Request) (BlockEdit, BlockEditError, error) {
	err := r.ParseForm()
	if err != nil {
		return BlockEdit{}, BlockEditError{}, fmt.Errorf("cannot parse form : %w", err)
	}
	a := BlockEdit{
		Name:    r.PostForm.Get(blockEditName),
		Content: r.PostForm.Get(blockEditContent),
	}
	errors := BlockEditError{}
	if err := json.Unmarshal([]byte(r.PostForm.Get(blockEditContent)), &a.Definition); err != nil {
		errors.Definition = errorInvalidJson
	}
	if a.Name == "" {
		errors.Name = errorCannotBeEmpty
	}
	if a.Content == "" {
		errors.Content = errorCannotBeEmpty
	}
	if len(a.Name) > maxTitleLen {
		errors.Name = errorTagetToBig
	}
	if len(a.Content) > maxContentLen {
		errors.Content = errorTagetToBig
	}
	re := regexp.MustCompile("^" + router.PathRegexp + "$")
	if !re.Match([]byte(a.Name)) {
		errors.Name = errorNotAPath
	}
	return a, errors, nil
}
