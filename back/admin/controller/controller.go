package controller

import (
	"app/admin/repository"
	"app/admin/translation"
	"app/admin/view"
	"app/pkg/clock"
	"app/pkg/i18n"
	"fmt"
	"net/http"
	"strings"
)

type Controller struct {
	Repository repository.Repository
	tr         i18n.Translator
	Clock      clock.Clock
}

func New(repo repository.Repository, c clock.Clock) (*Controller, error) {
	tr, err := translation.GetTranslator()
	if err != nil {
		return nil, fmt.Errorf("can't load translations : %w", err)
	}
	return &Controller{Repository: repo, Clock: c, tr: tr}, nil
}

func (c *Controller) render(w http.ResponseWriter, r *http.Request, fn view.ViewFunc) error {
	return c.renderWithStatus(w, r, http.StatusOK, fn)
}

func (c *Controller) renderWithStatus(w http.ResponseWriter, r *http.Request, st int, fn view.ViewFunc) error {
	lang := translation.GetLocal(w, r)
	tr := func(key string) string { return c.tr.Trans(key, lang) }
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(st)
	return fn(w, tr)
}

func (c *Controller) redirect(w http.ResponseWriter, r *http.Request, url string) error {
	_ = translation.GetLocal(w, r)
	http.Redirect(w, r, url, http.StatusSeeOther)
	return nil
}

func IsJsonRequest(r *http.Request) bool {
	json := "application/json"
	ct := "Content-Type"
	accept := "Accept"
	return r.Header.Get(ct) == json ||
		strings.Contains(r.Header.Get(accept), json)
}
