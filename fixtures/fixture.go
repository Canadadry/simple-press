package fixtures

import (
	"app/fixtures/api"
	"app/pkg/clock"
	"app/pkg/environment"
	"app/pkg/http/httpcaller"
	"fmt"
)

type Article struct {
	Title  string
	Author string
}

func Run(client httpcaller.Caller, c clock.Clock, articles []Article) (environment.Environment, error) {
	env := environment.New()
	api := api.New(client)

	for i, a := range articles {
		slug, err := api.AddArticle(a.Title, a.Author)
		if err != nil {
			return env, fmt.Errorf("cannot add article %s : %w", a.Title, err)
		}
		env.Store(fmt.Sprintf("article_%d_slug", i), slug)
	}
	return env, nil
}
