package form

import (
	"app/pkg/data"
	"fmt"
	"net/http"
	"strconv"
)

const (
	articleEditEditedBlockID = "edited_block_id"
)

type ParsedArticleEditBlockEdit struct {
	EditedBlockID       int64
	EditedBlockData     map[string]any
	EditedBlockPosition int
}

type ParsedArticleEditErrorBlockEdit struct {
	EditedBlockID       string
	EditedBlockData     string
	EditedBlockPosition string
}

func (pe ParsedArticleEditErrorBlockEdit) HasError() bool {
	if pe.EditedBlockID != "" {
		return true
	}
	if pe.EditedBlockData != "" {
		return true
	}
	if pe.EditedBlockPosition != "" {
		return true
	}
	return false
}

func ParseArticleEditBlockEdit(r *http.Request, get_previous_data func(int64) (map[string]any, bool)) (ParsedArticleEditBlockEdit, ParsedArticleEditErrorBlockEdit, error) {
	err := r.ParseForm()
	if err != nil {
		return ParsedArticleEditBlockEdit{}, ParsedArticleEditErrorBlockEdit{}, fmt.Errorf("cannot parse form : %w", err)
	}
	id, _ := strconv.ParseInt(r.PostForm.Get(articleEditEditedBlockID), 10, 64)
	a := ParsedArticleEditBlockEdit{
		EditedBlockID: id,
	}
	errors := ParsedArticleEditErrorBlockEdit{}
	form_data, ok := get_previous_data(id)
	if !ok {
		errors.EditedBlockID = errorInvalidId
		return a, errors, nil
	}
	form_data, err = data.ParseFormData(r, form_data)
	if err != nil {
		errors.EditedBlockData = errorInvalidJson
		return a, errors, nil
	}
	a.EditedBlockID = id
	a.EditedBlockData = form_data
	return a, errors, nil
}
