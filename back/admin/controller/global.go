package controller

import (
	"app/admin/form"
	"app/admin/view"
	"app/pkg/http/httpresponse"
	"fmt"
	"net/http"
)

func (c *Controller) GetGlobalDefinition(w http.ResponseWriter, r *http.Request) error {
	return view.GlobalDefinitionOk(w, map[string]any{})
}

func (c *Controller) PatchGlobalDefinition(w http.ResponseWriter, r *http.Request) error {
	def, errs, err := form.ParseGlobalDefinitionEdit(r)
	if err != nil {
		return fmt.Errorf("cannot parse form request : %w", err)
	}
	if errs.HasError {
		return httpresponse.BadRequest(w, errs)
	}

	return view.GlobalDefinitionOk(w, def.Definition)
}

func (c *Controller) GetGlobalData(w http.ResponseWriter, r *http.Request) error {
	return view.GlobalDataOk(w, map[string]any{})
}
func (c *Controller) PatchGlobalData(w http.ResponseWriter, r *http.Request) error {
	data, errs, err := form.ParseGlobalDataEdit(r)
	if err != nil {
		return fmt.Errorf("cannot parse form request : %w", err)
	}
	if errs.HasError {
		return httpresponse.BadRequest(w, errs)
	}

	return view.GlobalDefinitionOk(w, data.Data)
}
