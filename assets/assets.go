package public

import (
	"embed"
)

//go:embed public
var publicFiles embed.FS

func GetPublicFiles() embed.FS {
	return publicFiles
}
