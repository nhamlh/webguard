package web

import (
	"embed"
	"net/http"
	"html/template"
	"fmt"
	"path/filepath"
)

//go:embed static
var staticAssets embed.FS

const templatePath = "static/templates"

type templateData struct {
	Errors []string
}

// renderTemplate renders `file` inside `templatePath`
func renderTemplate(file string, data interface{}, w http.ResponseWriter) {
	t := template.Must(template.ParseFS(staticAssets, filepath.Join(templatePath, file)))
	err := t.Execute(w, data)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Error reading %s: %s", file, err.Error())))
	}
}
