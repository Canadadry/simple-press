package view

import (
	"io"
)

func Dashboard() ViewFunc {
	return func(w io.Writer, tr func(string) string) error {
		return render(w, tr,
			"template/pages/dashboard.html",
			TemplateData("DASHBOARD.page_title", map[string]any{}),
		)
	}
}
