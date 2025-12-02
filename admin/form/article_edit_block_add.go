package form

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

const (
	articleEditNewBlock = "new_block"
)

type ParsedArticleEditBlockAdd struct {
	AddedBlockID int64
}

type ParsedArticleEditErrorBlockAdd struct {
	AddedBlockID string
}

func (pe ParsedArticleEditErrorBlockAdd) HasError() bool {
	if pe.AddedBlockID != "" {
		return true
	}
	return false
}

func ParseArticleEditBlockAdd(r *http.Request, check_id func(context.Context, int64) (int, error)) (ParsedArticleEditBlockAdd, ParsedArticleEditErrorBlockAdd, error) {
	err := r.ParseForm()
	if err != nil {
		return ParsedArticleEditBlockAdd{}, ParsedArticleEditErrorBlockAdd{}, fmt.Errorf("cannot parse form : %w", err)
	}
	id, _ := strconv.ParseInt(r.PostForm.Get(articleEditNewBlock), 10, 64)
	a := ParsedArticleEditBlockAdd{
		AddedBlockID: id,
	}
	errors := ParsedArticleEditErrorBlockAdd{}
	c, err := check_id(r.Context(), id)
	if err != nil {
		return a, errors, err
	}
	if c == 0 {
		errors.AddedBlockID = errorInvalidId
	}
	return a, errors, nil
}
