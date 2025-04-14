package translation

import (
	"app/pkg/i18n"
	"bytes"
	_ "embed"
	"fmt"
)

//go:embed messages.csv
var EmbeddedTranslations []byte

const (
	LangFr  = "fr_FR"
	LangEng = "en_US"
)

func GetTranslator() (i18n.Translator, error) {
	translationFile := bytes.NewReader(EmbeddedTranslations)
	tr, err := i18n.LoadFromCsv(translationFile)
	if err != nil {
		return nil, fmt.Errorf("cannot read Embedded tranlation file %v", err)
	}
	return tr, nil
}
