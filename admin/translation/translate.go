package translation

import (
	"app/pkg/i18n"
	"bytes"
	_ "embed"
	"fmt"
)

//go:embed messages.csv
var EmbeddedTranslations []byte

//go:embed dashboard.en.md
var EmbeddedDashboardEn []byte

// //go:embed dashboard.fr.md
// var EmbeddedDashboardFr []byte

const dashboardkey = "dashboard.content"

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
	tr[dashboardkey] = map[string]string{
		LangEng: string(EmbeddedDashboardEn),
		LangFr:  string(EmbeddedDashboardEn),
	}
	return tr, nil
}
