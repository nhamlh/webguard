package web

import (
	"embed"
)

//go:embed static
var staticAssets embed.FS

type templateData struct {
	Errors []string
}
