package fixture

import (
	"app/admin/form"
	"app/cmd/migration"
	"app/config"
	"app/fixtures"
	"app/pkg/clock"
	"app/pkg/http/httpcaller"
	"embed"
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

//go:embed data
var data embed.FS

const (
	Action = "fixture"
)

var FixedNow = "2025-05-02T17:51:53+02:00"

const layoutFile = "data/layouts.csv"

func Run(c config.Parameters) error {
	_ = os.Remove(c.DatabaseUrl)
	err := migration.Run(c)
	if err != nil {
		return fmt.Errorf("migration failed : %w", err)
	}
	now, err := time.Parse(config.SerializerDateTimeFormat, FixedNow)
	if err != nil {
		return fmt.Errorf("invalid fixed clock date '%s': %w", FixedNow, err)
	}
	f, err := data.Open(layoutFile)
	if err != nil {
		return fmt.Errorf("cannot read embed file %s : %w", layoutFile, err)
	}
	defer f.Close()
	layoutReader := csv.NewReader(f)
	layout, err := layoutReader.ReadAll()
	layoutData := []form.Layout{}
	for _, l := range layout {
		layoutData = append(layoutData, form.Layout{
			Name: l[1],
			// Content: l[2],
		})
	}
	env, err := fixtures.Run(
		httpcaller.New(fmt.Sprintf("http://localhost:%d", c.Port), http.DefaultClient),
		&clock.Fixed{At: now},
		fixtures.FixtureData{
			Layouts: []form.Layout{
				{Name: "first"},
			},
			Articles: []form.Article{},
		},
	)
	if err != nil {
		return fmt.Errorf("fixture failed : %w", err)
	}
	env.SaveAt("tests/fixture.json")
	return os.Rename(c.DatabaseUrl, "tests/"+filepath.Base(c.DatabaseUrl))
}
