package form

import (
	"app/pkg/router"
	"fmt"
	"net/http"
	"regexp"
)

const (
	blockAddName = "name"
)

type Block struct {
	Name string
}

type BlockError struct {
	Name string
}

func (le BlockError) HasError() bool {
	if le.Name != "" {
		return true
	}
	if le.Name != "" {
		return true
	}
	return false
}

func ParseBlockAdd(r *http.Request) (Block, BlockError, error) {
	err := r.ParseForm()
	if err != nil {
		return Block{}, BlockError{}, fmt.Errorf("cannot parse form : %w", err)
	}
	l := Block{
		Name: r.PostForm.Get(blockAddName),
	}
	errors := BlockError{}
	if l.Name == "" {
		errors.Name = errorCannotBeEmpty
	}
	if len(l.Name) > maxTitleLen {
		errors.Name = errorTagetToBig
	}
	re := regexp.MustCompile("^" + router.PathRegexp + "$")
	if !re.Match([]byte(l.Name)) {
		errors.Name = errorNotAPath
	}
	return l, errors, nil
}
