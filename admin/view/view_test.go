package view

import (
	"app/pkg/flash"
	"bytes"
	"testing"
)

func fakeTr(key string) string {
	return key
}

func TestView(t *testing.T) {
	tests := map[string]ViewFunc{
		"article_list": ArticlesList(ArticlesListData{}),
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			err := tt(buf, fakeTr, flash.Message{})
			if err != nil {
				t.Fatalf("failed : %v", err)
			}
		})
	}
}
