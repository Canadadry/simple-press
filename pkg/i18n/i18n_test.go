package i18n

import (
	"bytes"
	"reflect"
	"testing"
)

func TestLoadFromCsv(t *testing.T) {
	tests := []struct {
		in       string
		local    string
		expected Translation
	}{
		{
			in: `,fr,en,de,it,es,nl,pt,ls
appendices.title,Annexes,Appendices,Anhänge,Appendici,Apéndices,Bijlagen,Appêndices,ls
authorization.applicant_signature_title,Signature du candidat*,Applicant signature*,Unterschrift des Antragstellers*,Apporre la firma*,Firma del solicitante*,Handtekening kandidaat*,Assinatura do candidato*,ls
`,
			expected: Translation{
				"appendices.title": map[string]string{
					"fr": "Annexes",
					"en": "Appendices",
					"de": "Anhänge",
					"it": "Appendici",
					"es": "Apéndices",
					"nl": "Bijlagen",
					"pt": "Appêndices",
					"ls": "ls",
				},
				"authorization.applicant_signature_title": map[string]string{
					"fr": "Signature du candidat*",
					"en": "Applicant signature*",
					"de": "Unterschrift des Antragstellers*",
					"it": "Apporre la firma*",
					"es": "Firma del solicitante*",
					"nl": "Handtekening kandidaat*",
					"pt": "Assinatura do candidato*",
					"ls": "ls",
				},
			},
		},
		{
			in:       `,fr,en,de,it,es,nl,pt,ls`,
			expected: Translation{},
		},
	}

	for i, tt := range tests {
		result, err := LoadFromCsv(bytes.NewBufferString(tt.in))
		if err != nil {
			t.Fatalf("[%d] failed %v", i, err)
		}
		if !reflect.DeepEqual(result, tt.expected) {
			t.Fatalf("[%d] exp \n%#v\n got \n%#v\n", i, tt.expected, result)
		}
	}
}

func TestTrans(t *testing.T) {
	tests := []struct {
		in       Translation
		key      string
		lang     string
		expected string
	}{
		{
			in:       Translation{},
			key:      "fake",
			lang:     "en",
			expected: "fake",
		},
		{
			in: Translation{
				"key": map[string]string{
					"en": "value",
				},
			},
			key:      "key",
			lang:     "fr",
			expected: "MISSING LOCAL: fr",
		},
		{
			in: Translation{
				"key": map[string]string{
					"en": "value",
				},
			},
			key:      "key",
			lang:     "en",
			expected: "value",
		},
	}

	for i, tt := range tests {
		result := tt.in.Trans(tt.key, tt.lang)
		if result != tt.expected {
			t.Fatalf("[%d] exp \n%#v\n got \n%#v\n", i, tt.expected, result)
		}
	}
}
