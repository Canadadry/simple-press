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
	baseTemplate      = "base.tmpl"
	baseTemplatesPath = "template/base/*"
)

type BasePage struct {
	Version  string
	Flash    flash.Message
	PageData interface{}
}

func render(w io.Writer, tr func(string) string, pageTemplatePath string, pageData interface{}) error {
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

func TemplateData(msg flash.Message, pageData interface{}) BasePage {
	bp := BasePage{
		Version:  version,
		PageData: pageData,
		Flash:    msg,
	}
	return bp
}
