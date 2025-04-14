package httpcache

import (
	"net/http"
)

const (
	DontCache = "dont-cache-this-request"
)

func NewCachedClient(cache Cache, client Client) Client {
	if cache == nil {
		return client
	}
	return &CachedClient{
		Client: client,
		Cache:  cache,
	}
}

type Client interface {
	Do(r *http.Request) (*http.Response, error)
}

type CachedClient struct {
	Client Client
	Cache  Cache
}

func (cc *CachedClient) Do(r *http.Request) (*http.Response, error) {
	if cc.Cache.IsHit(r) {
		return cc.Cache.Get(r)
	}
	rsp, err := cc.Client.Do(r)
	if err != nil {
		return nil, err
	}
	if r.Context().Value(DontCache) != nil {
		return rsp, nil
	}
	return cc.Cache.Store(r, rsp)
}
