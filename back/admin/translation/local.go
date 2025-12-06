package translation

import (
	"app/pkg/cookie"
	"net/http"
)

const (
	localCookieName = "_local"
)

func GetLocal(w http.ResponseWriter, r *http.Request) string {
	queryLang := r.URL.Query()[localCookieName]
	if len(queryLang) > 0 && validateLang(queryLang[0]) != "" {
		lang := validateLang(queryLang[0])
		SetLocal(w, lang)
		return lang
	}
	return GetLangFromCookie(r)
}

func GetLangFromCookie(r *http.Request) string {
	lang := cookie.Get(r, localCookieName)
	if lang == "" {
		return LangFr
	}
	return lang
}

func SetLocal(w http.ResponseWriter, next string) {
	cookie.SetForOneYear(w, localCookieName, next)
}

func validateLang(lang string) string {
	switch lang {
	case LangFr:
		return LangFr
	case LangEng:
		return LangEng
	}
	return ""
}

func SwitchLocal(previous string) string {
	switch previous {
	case LangFr:
		return LangEng
	case LangEng:
		return LangFr
	}
	return LangFr
}
