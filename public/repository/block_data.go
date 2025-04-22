package repository

import (
	"app/model/publicmodel"
	"app/pkg/sqlutil"
	"app/pkg/stacktrace"
	"context"
	"encoding/json"
)

type BlockData struct {
	ID        int64
	Position  int64
	Data      map[string]any
	ArticleID int64
	BlockID   int64
	BlockName string
}

func blockDataFromModel(model publicmodel.SelectBlockDataByArticleRow) (BlockData, error) {
	out := BlockData{
		ID:        model.ID,
		Position:  model.Position,
		BlockID:   model.BlockID,
		BlockName: model.Name.String,
	}
	err := json.Unmarshal([]byte(model.Data), &out.Data)
	return out, err
}

func (r *Repository) SelectBlockDataByArticle(ctx context.Context, articleID int64) ([]BlockData, error) {
	list, err := publicmodel.New(r.Db).SelectBlockDataByArticle(ctx, articleID)
	if err != nil {
		return nil, stacktrace.From(err)
	}
	return sqlutil.MapWithError(list, blockDataFromModel)
}
