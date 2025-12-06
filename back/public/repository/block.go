package repository

import (
	"app/model/publicmodel"
	"app/pkg/sqlutil"
	"app/pkg/stacktrace"
	"context"
	"encoding/json"
)

type Block struct {
	ID         int64
	Name       string
	Content    string
	Definition map[string]any
}

func blockFromModel(model publicmodel.Block) (Block, error) {
	out := Block{
		ID:      model.ID,
		Name:    model.Name,
		Content: model.Content,
	}
	if model.Definition == "" {
		out.Definition = map[string]any{}
		return out, nil
	}
	err := json.Unmarshal([]byte(model.Definition), &out.Definition)
	if err != nil {
		return out, stacktrace.From(err)
	}
	return out, nil
}

func (r *Repository) SelectAllBlock(ctx context.Context) ([]Block, error) {
	list, err := publicmodel.New(r.Db).SelectAllBlock(ctx)
	if err != nil {
		return nil, stacktrace.From(err)
	}
	return sqlutil.MapWithError(list, blockFromModel)
}
