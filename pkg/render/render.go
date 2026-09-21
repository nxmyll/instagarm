
package render

import (
	"bytes"
	"html/template"
	"insta/pkg/config"
	"log"
	"net/http"
	"path/filepath"
)

var app *config.AppConfig

// NewTemplates sets the config for the template package.
func NewTemplates(a *config.AppConfig) {
	app = a
}

// RenderTemplate renders an HTML template to the browser.
func RenderTemplate(w http.ResponseWriter, tpml string) {
	if app == nil {
		log.Println("Render error: app config is nil")
		http.Error(w, "Application configuration error", http.StatusInternalServerError)
		return
	}

	var (
		tc  map[string]*template.Template
		err error
	)

	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, err = CreateTemplateCache()
		if err != nil {
			log.Println("Template cache error:", err)
			http.Error(w, "Template cache error", http.StatusInternalServerError)
			return
		}
	}

	// Get the requested template from the cache.
	t, ok := tc[tpml]
	if !ok {
		log.Println("Template not found:", tpml)
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	// Execute the template into a buffer first.
	buf := new(bytes.Buffer)

	err = t.Execute(buf, nil)
	if err != nil {
		log.Println("Template execution error:", err)
		http.Error(w, "Template execution error", http.StatusInternalServerError)
		return
	}

	// Set the response type before writing the HTML.
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Send the rendered HTML to the browser.
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println("Error writing template response:", err)
	}
}

// CreateTemplateCache loads all HTML templates from the templates folder.
func CreateTemplateCache() (map[string]*template.Template, error) {
	theCache := make(map[string]*template.Template)

	// Find all HTML templates.
	pages, err := filepath.Glob("./templates/*.html")
	if err != nil {
		return theCache, err
	}

	if len(pages) == 0 {
		log.Println("Warning: no HTML templates found in ./templates")
	}

	// Parse each page and the shared layout.
	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return theCache, err
		}

		layoutPath := "./templates/layout.html"

		matches, err := filepath.Glob(layoutPath)
		if err != nil {
			return theCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseFiles(layoutPath)
			if err != nil {
				return theCache, err
			}
		}

		theCache[name] = ts
	}

	return theCache, nil
}
