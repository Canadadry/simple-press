package httpproxy

import (
	"context"
	"net/http"
)

type Client interface {
	Do(r *http.Request) (*http.Response, error)
}

type Proxy struct {
	url     string
	client  Client
	headers map[string]string
}

func New(base string, client Client) Proxy {
	if client == nil {
		client = &http.Client{}
	}
	return Proxy{url: base, client: client}
}

func (p Proxy) WithToken(name, value string) Proxy {
	return p.WithHeader(name, value)
}

func (p Proxy) WithHeader(name, value string) Proxy {
	caller := Proxy{
		url:     p.url,
		client:  p.client,
		headers: map[string]string{},
	}
	for k, v := range p.headers {
		caller.headers[k] = v
	}
	caller.headers[name] = value
	return caller
}

func (p Proxy) Get(ctx context.Context, url string) (*http.Response, error) {
	requestOptions := &requestOption{
		URL:         p.url + url,
		Method:      http.MethodGet,
		RequestBody: nil,
		Client:      p.client,
		Headers:     p.headers,
	}

	return doRequest(ctx, requestOptions)
}

func (p Proxy) Post(ctx context.Context, url string, payload interface{}) (*http.Response, error) {
	requestOptions := &requestOption{
		URL:         p.url + url,
		Method:      http.MethodPost,
		RequestBody: payload,
		Client:      p.client,
		Headers:     p.headers,
	}

	return doRequest(ctx, requestOptions)
}

func (p Proxy) Patch(ctx context.Context, url string, payload interface{}) (*http.Response, error) {
	requestOptions := &requestOption{
		URL:         p.url + url,
		Method:      http.MethodPatch,
		RequestBody: payload,
		Client:      p.client,
		Headers:     p.headers,
	}

	return doRequest(ctx, requestOptions)
}

func (p Proxy) Delete(ctx context.Context, url string) (*http.Response, error) {
	requestOptions := &requestOption{
		URL:         p.url + url,
		Method:      http.MethodDelete,
		RequestBody: struct{}{},
		Client:      p.client,
		Headers:     p.headers,
	}

	return doRequest(ctx, requestOptions)
}

func (p Proxy) DeleteWithBody(ctx context.Context, url string, payload interface{}) (*http.Response, error) {
	requestOptions := &requestOption{
		URL:         p.url + url,
		Method:      http.MethodDelete,
		RequestBody: payload,
		Client:      p.client,
		Headers:     p.headers,
	}

	return doRequest(ctx, requestOptions)
}
