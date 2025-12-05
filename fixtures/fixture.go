package fixtures

import (
	"app/admin/form"
	"app/fixtures/api"
	"app/pkg/clock"
	"app/pkg/environment"
	"app/pkg/http/httpcaller"
	"fmt"
)

type FixtureData struct {
	Layouts   []form.Layout
	Templates []form.Template
	Articles  []form.Article
	Blocks    []form.Block
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
	}

	for i, t := range fd.Templates {
		_, err := api.AddTemplate(t.Name)
		if err != nil {
			return env, fmt.Errorf("cannot add template %s : %w", t.Name, err)
		}
		env.Store(fmt.Sprintf("template_%d_name", i), t.Name)
	}
	for i, a := range fd.Articles {
		slug, err := api.AddArticle(a.Title, a.Author)
		if err != nil {
			return env, fmt.Errorf("cannot add article %s : %w", a.Title, err)
		}
		env.Store(fmt.Sprintf("article_%d_slug", i), slug)
	}
	for i, b := range fd.Blocks {
		name, err := api.AddBlock(b.Name)
		if err != nil {
			return env, fmt.Errorf("cannot add article %s : %w", b.Name, err)
		}
		env.Store(fmt.Sprintf("block_%d_name", i), name)
	}
	return env, nil
}
