package repository

import (
	"app/model/publicmodel"
	"app/pkg/stacktrace"
	"context"
)

type File struct {
	Name    string
	Content []byte
}

func (r *Repository) DownloadFile(ctx context.Context, name string) (File, bool, error) {
	list, err := publicmodel.New(r.Db).DownloadFile(ctx, name)
	if err != nil {
		return File{}, false, stacktrace.From(err)
	}
	if len(list) == 0 {
		return File{}, false, nil
	}
	return File{Name: name, Content: list[0]}, true, nil
}
