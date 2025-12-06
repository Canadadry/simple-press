package view

import (
	"io"
)

type ViewFunc func(io.Writer, func(string) string) error
