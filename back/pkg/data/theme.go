package data

type FormTheme struct {
	FormName          string
	FormMethod        string
	FormAction        string
	FormClass         string
	LabelClass        string
	InputClass        string
	SelectClass       string
	CheckboxClass     string
	FieldWrapper      string
	RowWrapper        string
	FieldsetClass     string
	LegendClass       string
	SubmitButtonClass string
	SubmitButtonName  string
}

var ThemeNoStyle = FormTheme{
	FormName:         "form",
	FormAction:       "/submit",
	FormMethod:       "POST",
	SubmitButtonName: "Submit",
}

var ThemeBootstrap = FormTheme{
	FormName:          "form",
	FormAction:        "/submit",
	FormMethod:        "POST",
	SubmitButtonName:  "Submit",
	FormClass:         "form-bootstrap",
	LabelClass:        "form-label",
	InputClass:        "form-control",
	SelectClass:       "form-select",
	CheckboxClass:     "form-check-input",
	FieldWrapper:      "mb-3",
	RowWrapper:        "row",
	FieldsetClass:     "mb-4",
	SubmitButtonClass: "btn btn-primary",
	LegendClass:       "fw-bold",
}
