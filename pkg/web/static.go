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

type templateData map[string]interface{}

// renderTemplate renders `name`.tpl inside `templatePath`
func renderTemplate(name string, data templateData, w http.ResponseWriter) {
	files := []string{filepath.Join(templatePath, "*.tpl"), filepath.Join(templatePath, `partials/*.tpl`)}

	t := template.Must(template.ParseFS(staticAssets, files...))
	err := t.ExecuteTemplate(w, name + ".tpl", data)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Error rendering %s: %s", name, err.Error())))
	}
}
