package fixtures

import (
	"app/admin/form"
	"app/admin/view"
	"app/fixtures/api"
	"app/pkg/clock"
	"app/pkg/environment"
	"app/pkg/http/httpcaller"
	"fmt"
	"io"
	"path/filepath"
)

type File struct {
	Filename string
	Content  io.ReadCloser
}

type FixtureData struct {
	Layouts   []form.LayoutEdit
	Templates []view.TemplateEditData
	Articles  []view.ArticleEditData
	Blocks    []view.BlockEditData
	Files     []File
}

func Run(client httpcaller.Caller, c clock.Clock, fd FixtureData) (environment.Environment, error) {
	env := environment.New()
	api := api.New(client)

	for i, l := range fd.Layouts {
		_, err := api.AddLayout(l.Name)
		if err != nil {
			return env, fmt.Errorf("cannot add layout %s : %w", l.Name, err)
		}
		env.Store(fmt.Sprintf("layout_%d_name", i), l.Name)
		err = api.EditLayout(l.Name, l.Content)
		if err != nil {
			return env, fmt.Errorf("cannot edit layout %s : %w", l.Name, err)
		}
	}

	for i, t := range fd.Templates {
		_, err := api.AddTemplate(t.Name)
		if err != nil {
			return env, fmt.Errorf("cannot add template %s : %w", t.Name, err)
		}
		env.Store(fmt.Sprintf("template_%d_name", i), t.Name)
		err = api.EditTemplate(t.Name, t.Content)
		if err != nil {
			return env, fmt.Errorf("cannot edit template %s : %w", t.Name, err)
		}
	}

	for i, b := range fd.Blocks {
		name, err := api.AddBlock(b.Name)
		if err != nil {
			return env, fmt.Errorf("cannot add block %s : %w", b.Name, err)
		}
		env.Store(fmt.Sprintf("block_%d_name", i), name)
		err = api.EditBlock(b)
		if err != nil {
			return env, fmt.Errorf("cannot edit block %s : %w", b.Name, err)
		}

	}

	for i, a := range fd.Articles {
		slug, err := api.AddArticle(a.Title, a.Author, filepath.Dir(a.Slug))
		if err != nil {
			return env, fmt.Errorf("cannot add article %s : %w", a.Title, err)
		}
		env.Store(fmt.Sprintf("article_%d_slug", i), slug)
		if len(a.Content) > 0 {
			err = api.EditArticleContent(slug, a.Content)
			if err != nil {
				return env, fmt.Errorf("cannot edit article %s : %w", slug, err)
			}
		}
		firstBlockID := 1
		id, err := api.EditArticleBlockAdd(slug, firstBlockID, 2)
		if err != nil {
			return env, fmt.Errorf("cannot edit article %s : %w", slug, err)
		}
		env.Store(fmt.Sprintf("block_data_%d_id", i), fmt.Sprintf("%d", id))

	}
	for i, f := range fd.Files {
		id, err := api.AddFile(f.Filename, f.Content)
		if err != nil {
			return env, fmt.Errorf("cannot add article %s : %w", f.Filename, err)
		}
		env.Store(fmt.Sprintf("file_%d_id", i), fmt.Sprintf("%d", id))
	}
	return env, nil
}
