package view

import (
	"fmt"
	"html/template"
	"strings"
	"time"
)

func trans(tr func(string) string) func(string, ...string) (string, error) {
	return func(key string, p ...string) (string, error) {
		if len(p)%2 == 1 {
			return "", fmt.Errorf("Must have a pair number of param : a key and a value")
		}

		translated := tr(key)
		for i := 0; i < len(p); i += 2 {
			translated = strings.ReplaceAll(translated, p[i], p[i+1])
		}

		return translated, nil
	}
}

func mailTo(email string) template.HTML {
	return template.HTML("<a href=\"mailto:" + email + "\">" + email + "</a>")
}

func incr(idx int) int {
	return idx + 1
}

func decr(idx int) int {
	return idx - 1
}

func mergeArgAndApply(fn func(string) (float64, error)) func(...string) string {
	return func(values ...string) string {
		str := ""
		for _, v := range values {
			str = fmt.Sprintf("%s %v", str, v)
		}
		result, err := fn(str)
		if err != nil {
			return err.Error()
		}
		return fmt.Sprintf("%v", result)
	}
}

func replace(text string, p ...interface{}) (string, error) {
	if len(p)%2 == 1 {
		return "", fmt.Errorf("must have a pair number of param : a key and a value")
	}
	for i := 0; i < len(p); i += 2 {
		text = strings.ReplaceAll(text, fmt.Sprintf("%v", p[i]), fmt.Sprintf("%v", p[i+1]))
	}
	return text, nil
}

/* #nosec */
func safe(str interface{}) template.HTML {
	return template.HTML(fmt.Sprintf("%v", str))
}

func safeUrl(str string) template.URL {
	return template.URL(str)
}

func escapeJs(text string) string {
	return template.JSEscapeString(text)
}

func formatNumber(num float64) string {
	strNum := fmt.Sprintf("%.2f", num)
	parts := strings.Split(strNum, ".")
	wholeReversed := reverse(parts[0])
	withSpaces := ""
	for i, char := range wholeReversed {
		if i%3 == 0 && i != 0 {
			withSpaces += " "
		}
		withSpaces += string(char)
	}
	result := reverse(withSpaces) + "." + parts[1]
	return result
}

func reverse(s string) string {
	rns := []rune(s)
	for i, j := 0, len(rns)-1; i < j; i, j = i+1, j-1 {
		rns[i], rns[j] = rns[j], rns[i]
	}
	return string(rns)
}

func formatDateTemplate(tr func(string) string) func(t time.Time) string {
	return func(t time.Time) string {
		return formatDate(t, tr)
	}
}

func formatDate(t time.Time, tr func(string) string) string {
	if t.IsZero() {
		return "-"
	}
	month := t.Format("Jan")
	return strings.Replace(t.Format("02 Jan 2006, 15:04"), month, tr("month.short."+month), 1)
}
