package controller

import (
	"app/admin/form"
	"app/admin/view"
	"app/pkg/http/httpresponse"
	"fmt"
	"net/http"
)

func (c *Controller) GetGlobalDefinition(w http.ResponseWriter, r *http.Request) error {
	def, err := c.Repository.GetGlobalDefinition(r.Context())
	if err != nil {
		return fmt.Errorf("cannot update global def : %w", err)
	}
	return view.GlobalDefinitionOk(w, def)
}

func (c *Controller) PatchGlobalDefinition(w http.ResponseWriter, r *http.Request) error {
	def, errs, err := form.ParseGlobalDefinitionEdit(r)
	if err != nil {
		return fmt.Errorf("cannot parse form request : %w", err)
	}
	if errs.HasError {
		return httpresponse.BadRequest(w, errs)
	}
	err = c.Repository.UpdateGlobalDefinition(r.Context(), def.Definition)
	if err != nil {
		return fmt.Errorf("cannot update global def : %w", err)
	}
	return view.GlobalDefinitionOk(w, def.Definition)
}

func (c *Controller) GetGlobalData(w http.ResponseWriter, r *http.Request) error {
	data, err := c.Repository.GetGlobalData(r.Context())
	if err != nil {
		return fmt.Errorf("cannot update global def : %w", err)
	}

	return view.GlobalDataOk(w, data)
}

func (c *Controller) PatchGlobalData(w http.ResponseWriter, r *http.Request) error {
	def, err := c.Repository.GetGlobalDefinition(r.Context())
	if err != nil {
		return fmt.Errorf("cannot update global def : %w", err)
	}
	data, errs, err := form.ParseGlobalDataEdit(r, def)
	if err != nil {
		return fmt.Errorf("cannot parse form request : %w", err)
	}
	if errs.HasError {
		return httpresponse.BadRequest(w, errs)
	}

	err = c.Repository.UpdateGlobalDefinition(r.Context(), data.Data)
	if err != nil {
		return fmt.Errorf("cannot update global def : %w", err)
	}
	return view.GlobalDataOk(w, data.Data)
}
