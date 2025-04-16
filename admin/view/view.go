package view

import (
	"app/pkg/flash"
	"io"
)

type ViewFunc func(io.Writer, func(string) string, flash.Message) error
