package parse

import (
	"html/template"
	"path/filepath"
)

func Parse() map[string]*template.Template {
	// Parse common files into base template
	baseTemplate := template.Must(template.ParseGlob("templates/layout/*.html"))
	baseTemplate = template.Must(baseTemplate.ParseGlob("templates/partials/*.html"))

	templates := make(map[string]*template.Template)
	
	// Get all page template files
	pageFiles, err := filepath.Glob("templates/*.html")
	if err != nil {
		panic(err)
	}

	// Parse each page template with base files
	for _, file := range pageFiles {
		name := filepath.Base(file)
		// Clone base template (includes all common files)
		tmpl := template.Must(baseTemplate.Clone())
		// Add the specific page template to the clone
		templates[name] = template.Must(tmpl.ParseFiles(file))
	}

	return templates
}
