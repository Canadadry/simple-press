package view

import (
	"app/pkg/eval"
	"app/pkg/flash"
	"embed"
	"html/template"
	"io"
)

//go:embed template
var templates embed.FS

var version = "dev"

const (
	baseTemplate      = "base.html"
	baseTemplatesPath = "template/base/*"
)

type BasePage[T any] struct {
	Version         string
	BreadcrumbItems any
	Flash           flash.Message
	PageData        T
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

func TemplateData[T any](msg flash.Message, pageData T) BasePage[T] {
	bp := BasePage[T]{
		Version:  version,
		PageData: pageData,
		Flash:    msg,
	}
	return bp
}
