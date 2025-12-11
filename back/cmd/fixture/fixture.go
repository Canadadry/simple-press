package fixture

import (
	"app/admin/form"
	"app/admin/view"
	"app/cmd/migration"
	"app/config"
	"app/fixtures"
	"app/pkg/clock"
	"app/pkg/http/httpcaller"
	"embed"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//go:embed data
var data embed.FS

const (
	Action = "fixture"
)

var FixedNow = "2025-05-02T17:51:53+02:00"

func readLayoutData(layoutFile string) ([]form.LayoutEdit, error) {
	f, err := data.Open(layoutFile)
	if err != nil {
		return nil, fmt.Errorf("cannot read embed file %s : %w", layoutFile, err)
	}
	defer f.Close()
	layoutReader := csv.NewReader(f)
	layout, err := layoutReader.ReadAll()
	layoutData := []form.LayoutEdit{}
	for idx, l := range layout {
		if idx == 0 {
			continue
		}
		layoutData = append(layoutData, form.LayoutEdit{
			Name:    l[1],
			Content: l[2],
		})
	}
	return layoutData, nil
}

func readTemplateData(templateFile string) ([]view.TemplateEditData, error) {
	f, err := data.Open(templateFile)
	if err != nil {
		return nil, fmt.Errorf("cannot read embed file %s : %w", templateFile, err)
	}
	defer f.Close()
	templateReader := csv.NewReader(f)
	template, err := templateReader.ReadAll()
	templateData := []view.TemplateEditData{}
	for idx, l := range template {
		if idx == 0 {
			continue
		}
		templateData = append(templateData, view.TemplateEditData{
			Name:    l[1],
			Content: l[2],
		})
	}
	return templateData, nil
}

func readArticleData(articleFile string) ([]view.ArticleEditData, error) {
	f, err := data.Open(articleFile)
	if err != nil {
		return nil, fmt.Errorf("cannot read embed file %s : %w", articleFile, err)
	}
	defer f.Close()
	articleReader := csv.NewReader(f)
	article, err := articleReader.ReadAll()
	articleData := []view.ArticleEditData{}
	for idx, l := range article {
		if idx == 0 {
			continue
		}
		articleData = append(articleData, view.ArticleEditData{
			Title: l[1],
			// Date:   l[2],
			Author:  l[3],
			Content: l[4],
			// Slug:l[5],
			// Draft: l[6],
			// LayoutId: l[7],
		})
	}
	return articleData, nil
}

func readBlockData(blockFile string) ([]view.BlockEditData, error) {
	f, err := data.Open(blockFile)
	if err != nil {
		return nil, fmt.Errorf("cannot read embed file %s : %w", blockFile, err)
	}
	defer f.Close()
	blockReader := csv.NewReader(f)
	block, err := blockReader.ReadAll()
	blockData := []view.BlockEditData{}
	for idx, l := range block {
		if idx == 0 {
			continue
		}
		def := map[string]any{}
		err := json.Unmarshal([]byte(l[3]), &def)
		if err != nil {
			return nil, fmt.Errorf("cannot parse definition of %s : %w", l[1], err)
		}
		blockData = append(blockData, view.BlockEditData{
			Name:       l[1],
			Content:    l[2],
			Definition: def,
		})
	}
	return blockData, nil
}

func readFileData(fileFile string) ([]fixtures.File, error) {
	f, err := data.Open(fileFile)
	if err != nil {
		return nil, fmt.Errorf("cannot read embed file %s : %w", fileFile, err)
	}
	defer f.Close()
	fileReader := csv.NewReader(f)
	file, err := fileReader.ReadAll()
	fileData := []fixtures.File{}
	for idx, l := range file {
		if idx == 0 {
			continue
		}
		var r io.ReadCloser
		switch l[3] {
		case "plain":
			r = io.NopCloser(strings.NewReader(l[2]))
		case "rstd":
			r = io.NopCloser(base64.NewDecoder(base64.RawStdEncoding, strings.NewReader(l[2])))
		case "std":
			r = io.NopCloser(base64.NewDecoder(base64.StdEncoding, strings.NewReader(l[2])))
		case "rurl":
			r = io.NopCloser(base64.NewDecoder(base64.RawURLEncoding, strings.NewReader(l[2])))
		case "url":
			r = io.NopCloser(base64.NewDecoder(base64.URLEncoding, strings.NewReader(l[2])))
		case "embed":
			r, err = data.Open(l[2])
			if err != nil {
				return nil, fmt.Errorf("cannot read embed file %s : %w", l[2], err)
			}
		}
		f := fixtures.File{
			Filename: l[1],
			Content:  io.NopCloser(r),
		}
		fileData = append(fileData, f)
	}
	return fileData, nil
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
	templateData, err := readTemplateData("data/templates.csv")
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
	fileData, err := readFileData("data/files.csv")
	if err != nil {
		return fmt.Errorf("cannot read block data : %w", err)
	}
	env, err := fixtures.Run(
		httpcaller.New(fmt.Sprintf("http://localhost:%d", c.Port), http.DefaultClient),
		&clock.Fixed{At: now},
		fixtures.FixtureData{
			Layouts:   layoutData,
			Templates: templateData,
			Articles:  articleData,
			Blocks:    blockData,
			Files:     fileData,
		},
	)
	if err != nil {
		return fmt.Errorf("fixture failed : %w", err)
	}
	env.SaveAt("tests/fixture.json")
	return os.Rename(c.DatabaseUrl, "tests/"+filepath.Base(c.DatabaseUrl))
}
