package httpcaller

import (
	"app/pkg/http/httpcache"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Client interface {
	Do(r *http.Request) (*http.Response, error)
}

type Caller struct {
	url     string
	client  Client
	headers map[string]string
	cache   httpcache.Cache
}

func New(base string, client Client) Caller {
	if client == nil {
		client = &http.Client{}
	}
	return Caller{url: base, client: client}
}

func (c Caller) WithToken(name, value string) Caller {
	return c.WithHeader(name, value)
}

func (c Caller) WithHeader(name, value string) Caller {
	caller := Caller{
		url:     c.url,
		client:  c.client,
		headers: map[string]string{},
	}
	for k, v := range c.headers {
		caller.headers[k] = v
	}
	caller.headers[name] = value
	return caller
}

func (c Caller) WithCache(cache httpcache.Cache) Caller {
	caller := Caller{
		url:     c.url,
		client:  c.client,
		headers: map[string]string{},
		cache:   cache,
	}
	for k, v := range c.headers {
		caller.headers[k] = v
	}
	return caller
}

func (c Caller) Get(ctx context.Context, url string, payload interface{}) (int, error) {
	requestOptions := &requestOption{
		URL:         combineUrl(c.url, url),
		Method:      http.MethodGet,
		RequestBody: nil,
		ResponesMap: map[int]interface{}{http.StatusOK: payload},
		Client:      httpcache.NewCachedClient(c.cache, c.client),
		Headers:     c.headers,
	}

	m, ok := payload.(map[int]interface{})
	if ok {
		requestOptions.ResponesMap = m
	}

	st, rsp, err := doRequest(ctx, requestOptions)
	if err != nil {
		return st, err
	}
	rsp.Body.Close()
	return st, nil
}

func (c Caller) GetFile(ctx context.Context, url string) (io.ReadCloser, int, error) {
	requestOptions := &requestOption{
		URL:         combineUrl(c.url, url),
		Method:      http.MethodGet,
		ResponesMap: map[int]interface{}{http.StatusOK: nil, http.StatusNotFound: nil},
		Client:      httpcache.NewCachedClient(c.cache, c.client),
		Headers:     c.headers,
	}

	status, response, err := doRequest(ctx, requestOptions)
	if err != nil {
		return nil, status, fmt.Errorf("%w : while getting file %s", err, url)
	}
	return response.Body, status, nil
}

func (c Caller) request(ctx context.Context, method, url string, payload interface{}, responseObjInstances map[int]interface{}) (int, error) {
	if responseObjInstances == nil {
		responseObjInstances = map[int]interface{}{
			http.StatusCreated:    nil,
			http.StatusBadRequest: nil,
		}
	}

	if payload == nil {
		payload = struct{}{}
	}

	requestOptions := &requestOption{
		URL:         combineUrl(c.url, url),
		Method:      method,
		RequestBody: payload,
		ResponesMap: responseObjInstances,
		Client:      c.client,
		Headers:     c.headers,
	}

	st, rsp, err := doRequest(ctx, requestOptions)
	if err != nil {
		return st, err
	}
	rsp.Body.Close()
	return st, nil
}

func (c Caller) Post(ctx context.Context, url string, payload interface{}, responseObjInstances map[int]interface{}) (int, error) {
	return c.request(ctx, http.MethodPost, url, payload, responseObjInstances)
}

func (c Caller) Put(ctx context.Context, url string, payload interface{}, responseObjInstances map[int]interface{}) (int, error) {
	return c.request(ctx, http.MethodPut, url, payload, responseObjInstances)
}

func (c Caller) PostFile(ctx context.Context, url string, file io.Reader, contentType string, responseObjInstances map[int]interface{}) (int, error) {
	return c.Post(ctx, url, File{Content: file, ContentType: contentType}, responseObjInstances)
}

func (c Caller) Patch(ctx context.Context, url string, payload interface{}, responseObjInstances map[int]interface{}) (int, error) {
	if responseObjInstances == nil {
		responseObjInstances = map[int]interface{}{
			http.StatusOK:         nil,
			http.StatusBadRequest: nil,
		}
	}
	if payload == nil {
		payload = struct{}{}
	}
	requestOptions := &requestOption{
		URL:         combineUrl(c.url, url),
		Method:      http.MethodPatch,
		RequestBody: payload,
		ResponesMap: responseObjInstances,
		Client:      c.client,
		Headers:     c.headers,
	}

	st, rsp, err := doRequest(ctx, requestOptions)
	if err != nil {
		return st, err
	}
	rsp.Body.Close()
	return st, nil
}

func (c Caller) Delete(ctx context.Context, url string) (int, error) {
	requestOptions := &requestOption{
		URL:         combineUrl(c.url, url),
		Method:      http.MethodDelete,
		RequestBody: struct{}{},
		ResponesMap: map[int]interface{}{http.StatusNoContent: nil},
		Client:      c.client,
		Headers:     c.headers,
	}

	st, rsp, err := doRequest(ctx, requestOptions)
	if err != nil {
		return st, err
	}
	rsp.Body.Close()
	return st, nil
}

func combineUrl(base, url string) string {
	if strings.HasPrefix(url, "http://") ||
		strings.HasPrefix(url, "https://") {
		return url
	}
	return base + url
}
