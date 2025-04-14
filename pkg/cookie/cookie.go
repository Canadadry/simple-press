package cookie

import (
	"net/http"
	"time"
)

var secureCookie = "true"

func Invalidate(w http.ResponseWriter, name string) {
	secure := secureCookie == "true"
	rcookie := http.Cookie{
		Name:     name,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		Path:     "/",
		Domain:   "",
		HttpOnly: true,
		Secure:   secure,
	}
	http.SetCookie(w, &rcookie)
}

func Get(r *http.Request, name string) string {
	for _, cookie := range r.Cookies() {
		if cookie.Name == name {
			return cookie.Value
		}
	}
	return ""
}

func SetForOneYear(w http.ResponseWriter, name, value string) {
	Set(w, name, value, 365*24*time.Hour)
}

func SetForSession(w http.ResponseWriter, name, value string) {
	secure := secureCookie == "true"
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Expires:  time.Now().Add(365 * 24 * time.Hour),
		Path:     "/",
		Domain:   "",
		HttpOnly: true,
		Secure:   secure,
	}
	http.SetCookie(w, &cookie)
}

func Set(w http.ResponseWriter, name, value string, duration time.Duration) {
	secure := secureCookie == "true"
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Expires:  time.Now().Add(duration),
		Path:     "/",
		Domain:   "",
		HttpOnly: true,
		Secure:   secure,
	}
	http.SetCookie(w, &cookie)
}
