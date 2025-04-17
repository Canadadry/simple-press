package repository

import (
	"app/pkg/clock"
	"app/pkg/sqlutil"
	"regexp"
	"strings"
	"unicode"
)

type Repository struct {
	Db    sqlutil.DBTX
	Clock clock.Clock
}

func slugify(title string) string {
	slug := strings.ToLower(title)
	slug = removeAccents(slug)
	re := regexp.MustCompile(`[^a-z0-9]+`)
	slug = re.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")

	return slug
}

func removeAccents(s string) string {
	var result []rune
	for _, r := range s {
		switch r {
		case 'à', 'á', 'â', 'ã', 'ä', 'å':
			result = append(result, 'a')
		case 'è', 'é', 'ê', 'ë':
			result = append(result, 'e')
		case 'ì', 'í', 'î', 'ï':
			result = append(result, 'i')
		case 'ò', 'ó', 'ô', 'õ', 'ö':
			result = append(result, 'o')
		case 'ù', 'ú', 'û', 'ü':
			result = append(result, 'u')
		case 'ç':
			result = append(result, 'c')
		case 'ñ':
			result = append(result, 'n')
		default:
			if unicode.IsLetter(r) || unicode.IsDigit(r) {
				result = append(result, r)
			} else {
				result = append(result, ' ')
			}
		}
	}
	return string(result)
}
