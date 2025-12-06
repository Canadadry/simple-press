package httpcache

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	maxBufferSize = 30 * 1024 * 1024
)

type Entry struct {
	ExpireAt time.Time
	Key      string
	Response Response
}

func loadEntriesFromFile(path string) (map[string]Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		if err != os.ErrNotExist {
			return map[string]Entry{}, nil
		}
		return nil, fmt.Errorf("cannot read cache file %w", err)
	}
	defer f.Close()
	cache := map[string]Entry{}
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 8192), maxBufferSize)
	for scanner.Scan() {
		e := Entry{}
		content := scanner.Bytes()
		err := json.Unmarshal(content, &e)
		if err != nil {
			return nil, fmt.Errorf("invalid json found while scanning %w : %s", err, string(scanner.Bytes()))
		}
		if e.ExpireAt.After(time.Now()) {
			cache[e.Key] = e
		}
	}
	return cache, nil
}

func clearAndsaveEntriesToFile(entries map[string]Entry, path string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("cannot write to cache file %w", err)
	}
	defer f.Close()
	for _, e := range entries {
		err := json.NewEncoder(f).Encode(e)
		if err != nil {
			return fmt.Errorf("cannot encode entry %s : %w", e.Key, err)
		}
	}
	return nil
}

func saveEntryToFile(path string, e Entry) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("cannot write to cache file %w", err)
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(e)
	if err != nil {
		return fmt.Errorf("cannot encode entry %s : %w", e.Key, err)
	}
	return nil
}

type LocalCache struct {
	cache    map[string]Entry
	path     string
	expireAt time.Time
	canStore func(rsp *http.Response) bool
}

func Local(path string, expireAt time.Time, canStore func(rsp *http.Response) bool) (*LocalCache, error) {
	c, err := loadEntriesFromFile(path)
	if err != nil {
		return nil, fmt.Errorf("loading %w", err)
	}
	err = clearAndsaveEntriesToFile(c, path)
	if err != nil {
		return nil, fmt.Errorf("saving %w", err)
	}
	return &LocalCache{
		cache:    c,
		path:     path,
		expireAt: expireAt,
		canStore: canStore,
	}, nil
}

func (m *LocalCache) IsHit(r *http.Request) bool {
	_, ok := m.cache[Hash(r)]
	return ok
}

func (m *LocalCache) Get(r *http.Request) (*http.Response, error) {
	e, ok := m.cache[Hash(r)]
	if !ok {
		return nil, fmt.Errorf("no cache found for %s", Hash(r))
	}

	return newResponseFrom(e.Response)

}

func (m *LocalCache) Store(r *http.Request, rsp *http.Response) (*http.Response, error) {
	if !m.canStore(rsp) {
		return rsp, nil
	}
	savedRsp, err := saveResponse(rsp)
	if err != nil {
		return nil, fmt.Errorf("cannot read body to store %s : %w", Hash(r), err)
	}
	e := Entry{
		Key:      Hash(r),
		ExpireAt: m.expireAt,
		Response: savedRsp,
	}
	m.cache[e.Key] = e
	err = saveEntryToFile(m.path, e)
	if err != nil {
		return nil, fmt.Errorf("cannot store  %s : %w", e.Key, err)
	}
	return m.Get(r)
}
