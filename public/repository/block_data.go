package repository

import (
	"app/model/publicmodel"
	"app/pkg/sqlutil"
	"app/pkg/stacktrace"
	"context"
)

type BlockData struct {
	ID          int64
	Position    int64
	Data        string
	ArticleID   int64
	BlockDataID int64
}

func blockDataFromModel(model publicmodel.BlockDatum) (BlockData, error) {
	return BlockData{}, nil
}

func (r *Repository) SelectBlockDataByArticle(ctx context.Context, articleID int64) ([]BlockData, error) {
	list, err := publicmodel.New(r.Db).SelectBlockDataByArticle(ctx, articleID)
	if err != nil {
		return nil, stacktrace.From(err)
	}
	return sqlutil.MapWithError(list, blockDataFromModel)
}
