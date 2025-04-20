package main

import (
	"app/pkg/data"
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

var exampleDef = map[string]any{
	"children": []any{
		map[string]any{
			"name":   "Jean",
			"gender": "Mr",
		},
	},
	"emails": map[string]any{
		"1": "jean@paul.com",
		"2": "paul@jean.com",
	},
}

func main() {
	http.HandleFunc("/", handleForm)
	http.HandleFunc("/submit", handleSubmit)

	log.Println("Listening on http://localhost:8080 ...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleForm(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	renderer := data.NewBootstrapRenderer(&buf, data.ThemeBootstrap)
	data.Render(exampleDef, renderer)
	formHTML := buf.String()

	tmpl := `
<!DOCTYPE html>
<html>
<head>
  <title>Bootstrap Form Demo</title>
  <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
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
		"Form": template.HTML(formHTML),
	})
}

func handleSubmit(w http.ResponseWriter, r *http.Request) {
	var err error
	exampleDef, err = data.ParseFormData(r, exampleDef)
	if err != nil {
		http.Error(w, "ParseFormData: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(exampleDef)
}
