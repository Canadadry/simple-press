package data

var ThemeNoStyle = FormTheme{
	FormClass:     "",
	LabelClass:    "",
	InputClass:    "",
	SelectClass:   "",
	CheckboxClass: "",
	FieldWrapper:  "",
	RowWrapper:    "",
	FieldsetClass: "",
	LegendClass:   "",
	Repeat:        5,
}

var ThemeBootstrap = FormTheme{
	FormClass:         "needs-validation",
	LabelClass:        "form-label",
	InputClass:        "form-control",
	SelectClass:       "form-select",
	CheckboxClass:     "form-check-input",
	FieldWrapper:      "mb-3",
	RowWrapper:        "row",
	FieldsetClass:     "mb-4",
	AddButtonClass:    "btn btn-secondary",
	DeleteButtonClass: "btn btn-outline-danger",
	SubmitButtonClass: "btn btn-primary",
	LegendClass:       "fw-bold",
}
