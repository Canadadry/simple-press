package data

type DynamicFormRenderer interface {
	BeginForm(name, action, method string)
	EndForm()

	BeginFieldset(label string)
	EndFieldset()

	BeginArray(name, path string)
	EndArray()

	BeginArrayItem(index int)
	EndArrayItem()

	Input(label, name, inputType string)
	Select(label, name string, options []string)
	Checkbox(label, name string)
}

func GenerateFormDynamicHTMLWithName(field Field, r DynamicFormRenderer, formName string) {
	r.BeginForm(formName, "/submit", "POST")
	renderDynamicField(field, r)
	r.EndForm()
}

func renderDynamicField(f Field, r DynamicFormRenderer) {
	switch f.Type {
	case "object":
		if !f.IsRoot && f.Key != "" {
			r.BeginFieldset(f.Key)
		}
		for _, child := range f.Children {
			renderDynamicField(child, r)
		}
		if !f.IsRoot && f.Key != "" {
			r.EndFieldset()
		}

	case "array":
		r.BeginFieldset(f.Key)
		r.BeginArray(f.Key, f.Path)

		r.BeginArrayItem(0)
		renderDynamicField(f.Children[0], r)
		r.EndArrayItem()

		r.EndArray()
		r.EndFieldset()

	case "string", "number", "email", "date":
		inputType := map[string]string{
			"string": "text", "number": "number", "email": "email", "date": "date",
		}[f.Type]
		r.Input(f.Key, f.Path, inputType)

	case "bool":
		r.Checkbox(f.Key, f.Path)

	case "enum":
		r.Select(f.Key, f.Path, f.EnumVals)
	}
}
