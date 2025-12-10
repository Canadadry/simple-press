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

func (r *Repository) CountFiles(ctx context.Context) (int, error) {
	c, err := adminmodel.New(r.Db).CountFile(ctx)
	return int(c), err
}

func (r *Repository) CountFileByName(ctx context.Context, name string) (int, error) {
	c, err := adminmodel.New(r.Db).CountFileByName(ctx, name)
	return int(c), err
}

func (r *Repository) UploadFile(ctx context.Context, f File) (int64, error) {
	id, err := adminmodel.New(r.Db).UploadFile(ctx, adminmodel.UploadFileParams{
		Name:    slugify(f.Name),
		Content: f.Content,
	})
	if err != nil {
		return 0, stacktrace.From(err)
	}
	return id, nil
}

func (r *Repository) DeleteFile(ctx context.Context, name string) error {
	err := adminmodel.New(r.Db).DeleteFile(ctx, name)
	if err != nil {
		return stacktrace.From(err)
	}
	return nil
}

func (r *Repository) GetFileList(ctx context.Context, limit, offset int) ([]File, error) {
	list, err := adminmodel.New(r.Db).GetFileList(ctx, adminmodel.GetFileListParams{
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
	list, err := adminmodel.New(r.Db).DownloadFile(ctx, name)
	if err != nil {
		return File{}, false, stacktrace.From(err)
	}
	if len(list) == 0 {
		return File{}, false, nil
	}
	return File{Name: name, Content: list[0]}, true, nil
}
