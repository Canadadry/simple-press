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

func readLayoutData(layoutFile string) ([]form.Layout, error) {
	f, err := data.Open(layoutFile)
	if err != nil {
		return nil, fmt.Errorf("cannot read embed file %s : %w", layoutFile, err)
	}
	defer f.Close()
	layoutReader := csv.NewReader(f)
	layout, err := layoutReader.ReadAll()
	layoutData := []form.Layout{}
	for idx, l := range layout {
		if idx == 0 {
			continue
		}
		layoutData = append(layoutData, form.Layout{
			Name: l[1],
			// Content: l[2],
		})
	}
	return layoutData, nil
}

func readArticleData(articleFile string) ([]form.Article, error) {
	f, err := data.Open(articleFile)
	if err != nil {
		return nil, fmt.Errorf("cannot read embed file %s : %w", articleFile, err)
	}
	defer f.Close()
	articleReader := csv.NewReader(f)
	article, err := articleReader.ReadAll()
	articleData := []form.Article{}
	for idx, l := range article {
		if idx == 0 {
			continue
		}
		articleData = append(articleData, form.Article{
			Title: l[1],
			// Date:   l[2],
			Author: l[3],
			// Content: l[4],
			// Slug:l[5],
			// Draft: l[6],
			// LayoutId: l[7],
		})
	}
	return articleData, nil
}

func readBlockData(blockFile string) ([]form.Block, error) {
	f, err := data.Open(blockFile)
	if err != nil {
		return nil, fmt.Errorf("cannot read embed file %s : %w", blockFile, err)
	}
	defer f.Close()
	blockReader := csv.NewReader(f)
	block, err := blockReader.ReadAll()
	blockData := []form.Block{}
	for idx, l := range block {
		if idx == 0 {
			continue
		}
		blockData = append(blockData, form.Block{
			Name: l[1],
			// Content: l[2],
		})
	}
	return blockData, nil
}

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
	layoutData, err := readLayoutData("data/layouts.csv")
	if err != nil {
		return fmt.Errorf("cannot read layout data : %w", err)
	}
	articleData, err := readArticleData("data/articles.csv")
	if err != nil {
		return fmt.Errorf("cannot read article data : %w", err)
	}
	blockData, err := readBlockData("data/blocks.csv")
	if err != nil {
		return fmt.Errorf("cannot read block data : %w", err)
	}
	env, err := fixtures.Run(
		httpcaller.New(fmt.Sprintf("http://localhost:%d", c.Port), http.DefaultClient),
		&clock.Fixed{At: now},
		fixtures.FixtureData{
			Layouts:  layoutData,
			Articles: articleData,
			Blocks:   blockData,
		},
	)
	if err != nil {
		return fmt.Errorf("fixture failed : %w", err)
	}
	env.SaveAt("tests/fixture.json")
	return os.Rename(c.DatabaseUrl, "tests/"+filepath.Base(c.DatabaseUrl))
}
