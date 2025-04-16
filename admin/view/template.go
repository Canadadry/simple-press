package view

import (
	"app/pkg/eval"
	"embed"
	"html/template"
	"io"
)

//go:embed template
var templates embed.FS

const (
	baseTemplate      = "base.html"
	baseTemplatesPath = "template/base/*"
)

type MenuItem struct {
	Name string
	Path string
	Icon string
}

type BasePage[T any] struct {
	PageTitle string
	Menu      []MenuItem
	PageData  T
}

func renderStatic(w io.Writer, tr func(string) string, pageTemplatePath string) error {
	return render[any](w, tr, pageTemplatePath, nil)
}

func render[T any](w io.Writer, tr func(string) string, pageTemplatePath string, pageData T) error {
	allFiles := []string{
		baseTemplatesPath,
		pageTemplatePath,
	}
	funcMap := template.FuncMap{
		"Trans":        trans(tr),
		"Decr":         decr,
		"Incr":         incr,
		"MailTo":       mailTo,
		"Eval":         mergeArgAndApply(eval.Eval),
		"Replace":      replace,
		"Safe":         safe,
		"SafeUrl":      safeUrl,
		"EscapeJS":     escapeJs,
		"NumberFormat": formatNumber,
		"DateFormat":   formatDateTemplate(tr),
	}
	templates, err := template.New(baseTemplate).Funcs(funcMap).ParseFS(templates, allFiles...)
	if err != nil {
		return err
	}

	return templates.ExecuteTemplate(w, baseTemplate, pageData)
}

func TemplateData[T any](pageTitle string, pageData T) BasePage[T] {
	bp := BasePage[T]{
		PageTitle: pageTitle,
		PageData:  pageData,
		Menu: []MenuItem{
			{Name: "MENU.articles", Path: "/admin/articles", Icon: "bi bi-body-text"},
			{Name: "MENU.layouts", Path: "/admin/layouts", Icon: "bi bi-grid-1x2"},
			{Name: "MENU.files", Path: "/admin/files", Icon: "bi bi-file-image"},
		},
	}
	return bp
}
