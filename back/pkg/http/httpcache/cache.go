package httpcache

import (
	"crypto/sha1"
	"fmt"
	"net/http"
)

func Hash(r *http.Request) string {
	hash := sha1.New()
	hash.Write([]byte(r.URL.String()))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

type Cache interface {
	IsHit(request *http.Request) bool
	Get(request *http.Request) (*http.Response, error)
	Store(request *http.Request, response *http.Response) (*http.Response, error)
}
