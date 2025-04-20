package repository

import (
	"app/model/adminmodel"
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

func blockFromModel(model adminmodel.Block) (Block, error) {
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

func (r *Repository) CountBlocks(ctx context.Context) (int, error) {
	c, err := adminmodel.New(r.Db).CountBlock(ctx)
	return int(c), err
}

func (r *Repository) CountBlockByName(ctx context.Context, name string) (int, error) {
	c, err := adminmodel.New(r.Db).CountBlockByName(ctx, name)
	return int(c), err
}

type CreateBlockParams struct {
	Name string
}

func (r *Repository) CreateBlock(ctx context.Context, l CreateBlockParams) error {
	_, err := adminmodel.New(r.Db).CreateBlock(ctx, adminmodel.CreateBlockParams{
		Name: l.Name,
	})
	if err != nil {
		return stacktrace.From(err)
	}
	return nil
}

func (r *Repository) DeleteBlock(ctx context.Context, name string) error {
	err := adminmodel.New(r.Db).DeleteBlock(ctx, name)
	if err != nil {
		return stacktrace.From(err)
	}
	return nil
}

func (r *Repository) GetBlockList(ctx context.Context, limit, offset int) ([]Block, error) {
	list, err := adminmodel.New(r.Db).GetBlockList(ctx, adminmodel.GetBlockListParams{
		Limit:  int64(limit),
		Offset: int64(offset),
	})
	if err != nil {
		return nil, stacktrace.From(err)
	}
	return sqlutil.Map(list, func(name string) Block {
		return Block{
			Name: name,
		}
	}), nil
}

func (r *Repository) SelectAllBlock(ctx context.Context) ([]Block, error) {
	list, err := adminmodel.New(r.Db).SelectAllBlock(ctx)
	if err != nil {
		return nil, stacktrace.From(err)
	}
	return sqlutil.MapWithError(list, blockFromModel)
}

func (r *Repository) SelectBlock(ctx context.Context, name string) (Block, bool, error) {
	list, err := adminmodel.New(r.Db).SelectBlock(ctx, name)
	if err != nil {
		return Block{}, false, stacktrace.From(err)
	}
	if len(list) == 0 {
		return Block{}, false, nil
	}
	out, err := blockFromModel(list[0])
	return out, true, err
}

func (r *Repository) UpdateBlock(ctx context.Context, name string, l Block) error {
	def, err := json.Marshal(l.Definition)
	if err != nil {
		return stacktrace.From(err)
	}
	err = adminmodel.New(r.Db).UpdateBlock(ctx, adminmodel.UpdateBlockParams{
		Name:       l.Name,
		Content:    l.Content,
		Definition: string(def),
		Name_2:     name,
	})
	if err != nil {
		return stacktrace.From(err)
	}
	return nil
}
