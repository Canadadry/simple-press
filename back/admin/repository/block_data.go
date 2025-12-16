package repository

import (
	"app/model/adminmodel"
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

func blockDataFromModel(model adminmodel.SelectBlockDataByArticleRow) (BlockData, error) {
	out := BlockData{
		ID:        model.ID,
		Position:  model.Position,
		BlockID:   model.BlockID,
		BlockName: model.Name.String,
	}
	err := json.Unmarshal([]byte(model.Data), &out.Data)
	return out, err
}

func (r *Repository) CountBlockDataByID(ctx context.Context, id int64) (int, error) {
	c, err := adminmodel.New(r.Db).CountBlockDataByID(ctx, id)
	return int(c), err
}

type CreateBlockDataParams struct {
	Position  int64
	Block     Block
	ArticleID int64
}

func (r *Repository) CreateBlockData(ctx context.Context, l CreateBlockDataParams) (int64, error) {
	if len(l.Block.Definition) == 0 {
		return 0, stacktrace.Errorf("empty block definition")
	}
	data, err := json.Marshal(l.Block.Definition)
	if err != nil {
		return 0, stacktrace.From(err)
	}
	id, err := adminmodel.New(r.Db).CreateBlockData(ctx, adminmodel.CreateBlockDataParams{
		ArticleID: l.ArticleID,
		Position:  l.Position,
		BlockID:   l.Block.ID,
		Data:      string(data),
	})
	if err != nil {
		return 0, stacktrace.From(err)
	}
	return id, nil
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

func (r *Repository) SelectBlockDataByID(ctx context.Context, ID int64) (BlockData, bool, error) {
	b, err := adminmodel.New(r.Db).SelectBlockDataByID(ctx, ID)
	if err != nil {
		return BlockData{}, false, stacktrace.From(err)
	}
	if len(b) == 0 {
		return BlockData{}, false, nil
	}
	out := BlockData{
		ID:        b[0].ID,
		Position:  b[0].Position,
		ArticleID: b[0].ArticleID,
		BlockID:   b[0].BlockID,
	}
	err = json.Unmarshal([]byte(b[0].Data), &out.Data)
	return out, true, err
}

func (r *Repository) UpdateBlockData(ctx context.Context, l BlockData) error {
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
