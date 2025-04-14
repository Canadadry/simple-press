package repository

import (
	"app/model/adminmodel"
	"app/pkg/sqlutil"
	"app/pkg/stacktrace"
	"context"
)

type File struct {
	Name    string
	Content []byte
}

func (r *Repository) UploadFile(ctx context.Context, f File) error {
	_, err := adminmodel.New(r.db).UploadFile(ctx, adminmodel.UploadFileParams{
		Name:    f.Name,
		Content: f.Content,
	})
	if err != nil {
		return stacktrace.From(err)
	}
	return nil
}

func (r *Repository) DeleteFile(ctx context.Context, name string) error {
	err := adminmodel.New(r.db).DeleteFile(ctx, name)
	if err != nil {
		return stacktrace.From(err)
	}
	return nil
}

func (r *Repository) GetFileList(ctx context.Context, limit, offset int) ([]File, error) {
	list, err := adminmodel.New(r.db).GetFileList(ctx, adminmodel.GetFileListParams{
		Limit:  int64(limit),
		Offset: int64(offset),
	})
	if err != nil {
		return nil, stacktrace.From(err)
	}
	return sqlutil.Map(list, func(name string) File {
		return File{
			Name: name,
		}
	}), nil
}

func (r *Repository) DownloadFile(ctx context.Context, name string) (File, bool, error) {
	list, err := adminmodel.New(r.db).DownloadFile(ctx, name)
	if err != nil {
		return File{}, false, stacktrace.From(err)
	}
	if len(list) == 0 {
		return File{}, false, nil
	}
	return File{Name: name, Content: list[0]}, true, nil
}
