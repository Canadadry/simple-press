package urlutil

import (
	"net/url"
)

func GetFullPath(u *url.URL) string {
	ret := u.Path
	if ret == "" {
		ret += "/"
	}
	if u.RawQuery == "" {
		return ret
	}
	return ret + "?" + u.RawQuery
}
