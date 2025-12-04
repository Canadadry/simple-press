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
	Layouts  []form.Layout
	Articles []form.Article
}

func Run(client httpcaller.Caller, c clock.Clock, fd FixtureData) (environment.Environment, error) {
	env := environment.New()
	api := api.New(client)

	for i, l := range fd.Layouts {
		id, err := api.AddLayout(l.Name)
		if err != nil {
			return env, fmt.Errorf("cannot add article %s : %w", l.Name, err)
		}
		env.Store(fmt.Sprintf("layout_%d_id", i), fmt.Sprintf("%v", id))
	}
	for i, a := range fd.Articles {
		slug, err := api.AddArticle(a.Title, a.Author)
		if err != nil {
			return env, fmt.Errorf("cannot add article %s : %w", a.Title, err)
		}
		env.Store(fmt.Sprintf("article_%d_slug", i), slug)
	}
	return env, nil
}
