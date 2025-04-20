package repository

import (
	"app/model/adminmodel"
	"app/pkg/sqlutil"
	"app/pkg/stacktrace"
	"context"
	"encoding/json"
)

type BlockData struct {
	ID          int64
	Position    int64
	Data        string
	ArticleID   int64
	BlockDataID int64
}

func blockDataFromModel(model adminmodel.BlockDatum) (BlockData, error) {
	return BlockData{}, nil
}

type CreateBlockDataParams struct {
	Position  int64
	Block     Block
	ArticleID int64
}

func (r *Repository) CreateBlockData(ctx context.Context, l CreateBlockDataParams) error {
	data, err := json.Marshal(l.Block.Definition)
	if err != nil {
		return stacktrace.From(err)
	}
	_, err = adminmodel.New(r.Db).CreateBlockData(ctx, adminmodel.CreateBlockDataParams{
		ArticleID: l.ArticleID,
		Position:  l.Position,
		BlockID:   l.Block.ID,
		Data:      string(data),
	})
	if err != nil {
		return stacktrace.From(err)
	}
	return nil
}

func (r *Repository) DeleteBlockData(ctx context.Context, id int64) error {
	err := adminmodel.New(r.Db).DeleteBlockData(ctx, id)
	if err != nil {
		return stacktrace.From(err)
	}
	return nil
}

func (r *Repository) SelectBlockDataByArticle(ctx context.Context, articleID int64) ([]BlockData, error) {
	list, err := adminmodel.New(r.Db).SelectBlockDataByArticle(ctx, articleID)
	if err != nil {
		return nil, stacktrace.From(err)
	}
	return sqlutil.MapWithError(list, blockDataFromModel)
}

func (r *Repository) UpdateBlockData(ctx context.Context, name string, l BlockData) error {
	data, err := json.Marshal(l.Data)
	if err != nil {
		return stacktrace.From(err)
	}
	err = adminmodel.New(r.Db).UpdateBlockData(ctx, adminmodel.UpdateBlockDataParams{
		ID:       l.ID,
		Position: l.Position,
		Data:     string(data),
	})
	if err != nil {
		return stacktrace.From(err)
	}
	return nil
}
