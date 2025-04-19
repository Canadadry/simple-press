package main

import (
	"app/pkg/data"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

var exampleDef = map[string]any{
	"children": []any{
		map[string]any{
			"name":   "string",
			"gender": "enum:Mr;Mme",
		},
	},
	"emails": []any{"email"},
}

var bootstrapTheme = data.FormTheme{
	FormClass:         "needs-validation",
	LabelClass:        "form-label",
	InputClass:        "form-control",
	SelectClass:       "form-select",
	CheckboxClass:     "form-check-input",
	FieldWrapper:      "mb-3",
	RowWrapper:        "row",
	FieldsetClass:     "mb-4",
	AddButtonClass:    "btn btn-secondary",
	SubmitButtonClass: "btn btn-primary",
	LegendClass:       "fw-bold",
}

func main() {
	http.HandleFunc("/", handleForm)
	http.HandleFunc("/submit", handleSubmit)

	log.Println("Listening on http://localhost:8080 ...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleForm(w http.ResponseWriter, r *http.Request) {
	field := data.Parse(exampleDef, true)

	html := data.GenerateFormDynamicHTMLWithName(field, bootstrapTheme, "demo")
	js := data.GenerateDynamicJS(field)

	tmpl := `
<!DOCTYPE html>
<html>
<head>
  <title>Bootstrap Form Demo</title>
  <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
  <script>
  {{ .JS }}
  </script>
</head>
<body class="p-4">
  <div class="container">
    <h1 class="mb-4">Dynamic Form (Bootstrap)</h1>
    {{ .Form }}
  </div>
</body>
</html>`

	t := template.Must(template.New("page").Parse(tmpl))
	_ = t.Execute(w, map[string]any{
		"Form": template.HTML(html),
		"JS":   template.JS(js),
	})
}

func handleSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "ParseForm: "+err.Error(), http.StatusBadRequest)
		return
	}

	field := data.Parse(exampleDef, true)
	result, err := data.ParseFormData(r, field)
	if err != nil {
		http.Error(w, "ParseFormData: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}
