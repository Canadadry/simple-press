package main

import (
	"app/pkg/data"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

var (
	form_data = map[string]any{}
)

func main() {

	defaultValue := map[string]any{
		"children": map[string]any{
			"name":   "Jean",
			"gender": "Mr",
		},
		"emails": map[string]any{
			"1": "jean@paul.com",
			"2": "paul@jean.com",
		},
	}
	defaultStr, _ := json.Marshal(defaultValue)

	rawJSON := ""
	flag.StringVar(&rawJSON, "data", string(defaultStr), "JSON input defining the form structure")
	flag.Parse()

	if rawJSON != "" {
		if err := json.Unmarshal([]byte(rawJSON), &form_data); err != nil {
			fmt.Fprintf(os.Stderr, "Invalid JSON passed to -data: %v\n", err)
			os.Exit(1)
		}
	} else {
		form_data["value"] = "none"
	}

	http.HandleFunc("/", handleForm)
	http.HandleFunc("/submit", handleSubmit)

	log.Println("Listening on http://localhost:8080 ...")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func handleForm(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	renderer := data.NewBootstrapRenderer(&buf, data.ThemeBootstrap)
	err := data.Render(form_data, renderer)
	if err != nil {
		fmt.Fprintln(w, "invalid data", err)
		return
	}
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
	form_data, err = data.ParseFormData(r, form_data)
	if err != nil {
		http.Error(w, "ParseFormData: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(form_data)
}
