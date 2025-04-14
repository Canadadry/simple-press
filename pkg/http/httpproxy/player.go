package httpproxy

import (
	"app/pkg/urlutil"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Player struct {
	Records    Records
	Read       map[string]int
	TokenAsked int
}

func NewPlayer(reccord string) (*Player, error) {
	c, err := LoadZip(reccord)
	p := &Player{
		Records: c,
		Read:    map[string]int{},
	}
	return p, err
}

func (p *Player) GetNextToken(string, int) (string, error) {
	if p.TokenAsked >= len(p.Records.TokenStack) {
		return "", fmt.Errorf("can't return more token")
	}
	token := p.Records.TokenStack[p.TokenAsked]
	p.TokenAsked++
	return token, nil
}

func (p *Player) Do(r *http.Request) (*http.Response, error) {
	routeName := fmt.Sprintf("%s:%s", r.Method, urlutil.GetFullPath(r.URL))
	routes, ok := p.Records.Responses[routeName]
	if !ok {
		return nil, fmt.Errorf("route '%s' not found in API", routeName)
	}

	step := p.Read[routeName]
	if step >= len(routes) {
		return nil, fmt.Errorf("route '%s' called to many times (%d)", routeName, step)
	}
	p.Read[routeName] = step + 1

	resp := routes[step]
	return &http.Response{
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		Body:       io.NopCloser(strings.NewReader(resp.Body)),
	}, nil
}

func (p *Player) Env() func(string) string {
	return func(k string) string {
		return p.Records.Env[k]
	}
}
