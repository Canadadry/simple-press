package httpquery

import (
	"app/pkg/http/httpquery/parsestr"
	"net/http"
	"strconv"
)

func Has(r *http.Request, p string) bool {
	_, ok := r.URL.Query()[p]
	return ok
}

func ReadString(r *http.Request, p string, d string) string {
	v, ok := r.URL.Query()[p]
	if !ok {
		return d
	}
	return v[0]
}

func ReadInt(r *http.Request, p string, d int) int {
	arrayOfValues, ok := r.URL.Query()[p]
	if !ok || len(arrayOfValues) != 1 {
		return d
	}
	i, err := strconv.Atoi(arrayOfValues[0])
	if err != nil {
		return d
	}
	return i
}

func ReadArray(r *http.Request, p string, d map[string]string) map[string]string {
	q := r.URL.RawQuery
	values, err := parsestr.ParseQuery(q)
	if err != nil {
		return d
	}

	start := p + "["
	end := "]"

	result := map[string]string{}

	for k, v := range values {
		if len(k) < (len(start) + 2) {
			continue
		}
		if k[:len(start)] != start {
			continue
		}
		if k[len(k)-1:] != end {
			continue
		}
		realKey := k[len(start) : len(k)-1]
		if len(v) > 0 {
			result[realKey] = v[0]
		}
	}
	if len(result) == 0 {
		return d
	}

	return result
}
