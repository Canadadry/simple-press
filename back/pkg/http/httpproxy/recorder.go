package httpproxy

import (
	"app/pkg/http/httpcaller"
	"app/pkg/urlutil"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Recorder struct {
	Captured      map[string][]Response
	Server        httpcaller.Client
	TokenStack    []string
	ResponseStack []Response
	Env           map[string]string
}

type Response struct {
	StatusCode int                 `json:"status_code"`
	Body       string              `json:"body"`
	Header     map[string][]string `json:"-"`
}

func NewRecorder(c httpcaller.Client) *Recorder {
	return &Recorder{
		Server:   c,
		Captured: map[string][]Response{},
	}
}

func (p *Recorder) SaveAndReturnToken(token string, err error) (string, error) {
	if err == nil {
		p.TokenStack = append(p.TokenStack, token)
	}
	return token, err
}

func (p *Recorder) Do(r *http.Request) (*http.Response, error) {
	resp, err := p.Server.Do(r)
	if err != nil {
		return resp, err
	}
	defer resp.Body.Close()

	routeName := fmt.Sprintf("%s:%s", r.Method, urlutil.GetFullPath(r.URL))

	buf := &bytes.Buffer{}
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("can't read server response : %w", err)
	}

	savedResponse := Response{
		StatusCode: resp.StatusCode,
		Body:       buf.String(),
		Header:     resp.Header,
	}

	p.Captured[routeName] = append(p.Captured[routeName], savedResponse)
	p.ResponseStack = append(p.ResponseStack, savedResponse)

	return &http.Response{
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		Body:       io.NopCloser(strings.NewReader(savedResponse.Body)),
	}, nil
}
