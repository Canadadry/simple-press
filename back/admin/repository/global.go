package repository

import (
	"app/model/adminmodel"
	"app/pkg/stacktrace"
	"context"
	"encoding/json"
)

func (r *Repository) GetGlobalDefinition(ctx context.Context) (map[string]any, error) {
	def, err := adminmodel.New(r.Db).GetGlobalDefinition(ctx)
	if err != nil {
		return nil, stacktrace.From(err)
	}
	out := map[string]any{}
	err = json.Unmarshal([]byte(def), &out)
	if err != nil {
		return nil, stacktrace.From(err)
	}
	return out, nil
}

func (r *Repository) UpdateGlobalDefinition(ctx context.Context, def map[string]any) error {
	marshaledDef, err := json.Marshal(def)
	if err != nil {
		return stacktrace.From(err)
	}
	err = adminmodel.New(r.Db).UpdateGlobalDefinition(ctx, string(marshaledDef))
	if err != nil {
		return stacktrace.From(err)
	}
	return nil
}

func (r *Repository) GetGlobalData(ctx context.Context) (map[string]any, error) {
	data, err := adminmodel.New(r.Db).GetGlobalDefinition(ctx)
	if err != nil {
		return nil, stacktrace.From(err)
	}
	out := map[string]any{}
	err = json.Unmarshal([]byte(data), &out)
	if err != nil {
		return nil, stacktrace.From(err)
	}
	return out, nil
}

func (r *Repository) UpdateGlobalData(ctx context.Context, data map[string]any) error {
	marshaledData, err := json.Marshal(data)
	if err != nil {
		return stacktrace.From(err)
	}
	err = adminmodel.New(r.Db).UpdateGlobalData(ctx, string(marshaledData))
	if err != nil {
		return stacktrace.From(err)
	}
	return nil
}
