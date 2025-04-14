package httpcache

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

const (
	cachename = "test.cache"
)

var (
	ErrTestingClientCalled = fmt.Errorf("testing client called")
)

type TestingClient struct {
}

func (tc *TestingClient) Do(_ *http.Request) (*http.Response, error) {
	return nil, ErrTestingClientCalled
}

func TestStoreAndLoad(t *testing.T) {
	defer os.Remove(cachename)

	l, err := Local(cachename, time.Now().Add(time.Hour), func(rsp *http.Response) bool {
		return true
	})
	if err != nil {
		t.Fatalf("failed to build local cache : %v", err)
	}
	cc := NewCachedClient(l, http.DefaultClient)

	req, err := http.NewRequest("GET", "https://www.google.com", nil)
	if err != nil {
		t.Fatalf("failed to create google request : %v", err)
	}

	_, err = cc.Do(req)
	if err != nil {
		t.Fatalf("failed to request google %v", err)
	}

	f, err := os.Open(cachename)
	if err != nil {
		t.Fatalf("cannot open cache file %v", err)
	}
	content, err := io.ReadAll(f)
	f.Close()
	if err != nil {
		t.Fatalf("cannot open cache file %v", err)
	}
	if len(content) == 0 {
		t.Fatalf("file empty")
	}
	t.Logf("%s\n", string(content))

	l, err = Local(cachename, time.Now().Add(time.Hour), func(rsp *http.Response) bool {
		return true
	})
	if err != nil {
		t.Fatalf("failed to build local cache the second time : %v", err)
	}
	cc = NewCachedClient(l, &TestingClient{})
	_, err = cc.Do(req)
	if err != nil {
		t.Fatalf("failed requesting google the second time failed : %v", err)
	}

}

func TestStoreAndLoadAfterExpiry(t *testing.T) {
	defer os.Remove(cachename)

	l, err := Local(cachename, time.Now().Add(-time.Minute), func(rsp *http.Response) bool {
		return true
	})
	if err != nil {
		t.Fatalf("failed to build local cache : %v", err)
	}
	cc := NewCachedClient(l, http.DefaultClient)

	req, err := http.NewRequest("GET", "https://www.google.com", nil)
	if err != nil {
		t.Fatalf("failed to create google request %v", err)
	}

	_, err = cc.Do(req)
	if err != nil {
		t.Fatalf("failed to request google : %v", err)
	}

	l, err = Local(cachename, time.Now().Add(-time.Minute), func(rsp *http.Response) bool {
		return true
	})
	if err != nil {
		t.Fatalf("failed to build local cache the second time : %v", err)
	}

	cc = NewCachedClient(l, &TestingClient{})
	_, err = cc.Do(req)
	if err != ErrTestingClientCalled {
		t.Fatalf("should have called testing client got : %v", err)
	}

}

func TestDontCacheThisRequest(t *testing.T) {
	defer os.Remove(cachename)

	l, err := Local(cachename, time.Now().Add(time.Hour), func(rsp *http.Response) bool {
		return true
	})
	if err != nil {
		t.Fatalf("failed to build local cache : %v", err)
	}
	cc := NewCachedClient(l, http.DefaultClient)

	req, err := http.NewRequest("GET", "https://www.google.com", nil)
	if err != nil {
		t.Fatalf("failed to create google request : %v", err)
	}

	req = req.WithContext(context.WithValue(context.TODO(), DontCache, DontCache))

	_, err = cc.Do(req)
	if err != nil {
		t.Fatalf("failed to request google %v", err)
	}

	f, err := os.Open(cachename)
	if err != nil {
		t.Fatalf("cannot open cache file %v", err)
	}
	content, err := io.ReadAll(f)
	f.Close()
	if err != nil {
		t.Fatalf("cannot open cache file %v", err)
	}
	if len(content) != 0 {
		t.Logf("%s\n", string(content))
		t.Fatalf("file should be empty")
	}

	l, err = Local(cachename, time.Now().Add(time.Hour), func(rsp *http.Response) bool {
		return true
	})
	if err != nil {
		t.Fatalf("failed to build local cache the second time : %v", err)
	}
	cc = NewCachedClient(l, &TestingClient{})
	_, err = cc.Do(req)
	if err == nil {
		t.Fatalf("should have called testing client : %v", err)
	}

}

func TestCanCache(t *testing.T) {
	defer os.Remove(cachename)

	l, err := Local(cachename, time.Now().Add(time.Hour), func(rsp *http.Response) bool {
		return rsp.StatusCode != http.StatusNotFound
	})
	if err != nil {
		t.Fatalf("failed to build local cache : %v", err)
	}
	cc := NewCachedClient(l, http.DefaultClient)
	req, err := http.NewRequest("GET", "https://www.google.com/fake", nil)
	if err != nil {
		t.Fatalf("failed to create google request %v", err)
	}

	_, err = cc.Do(req)
	if err != nil {
		t.Fatalf("failed to request google : %v", err)
	}

	cc = NewCachedClient(l, &TestingClient{})
	_, err = cc.Do(req)
	if err == nil {
		t.Fatalf("should not have called testing client got : %v", err)
	}
}
