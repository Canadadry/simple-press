package i18n

import (
	"encoding/csv"
	"io"
)

type Translation map[string]map[string]string

type Translator interface {
	Trans(key string, language string) string
}

func LoadFromCsv(file io.Reader) (Translation, error) {
	if file == nil {
		return Translation{}, nil
	}
	csvReader := csv.NewReader(file)
	csvReader.Comma = ','

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) <= 1 {
		return Translation{}, nil
	}

	t := Translation{}
	recordLen := len(records[0])
	for _, line := range records[1:] {
		if _, ok := t[line[0]]; !ok {
			t[line[0]] = make(map[string]string)
		}
		for langIdx := 1; langIdx < recordLen; langIdx++ {
			t[line[0]][records[0][langIdx]] = line[langIdx]
		}
	}
	return t, nil
}

func (t Translation) Trans(key, language string) string {
	transForKey, ok := t[key]
	if !ok {
		return key
	}
	str, ok := transForKey[language]
	if !ok {
		return "MISSING LOCAL: " + language
	}

	return str
}
