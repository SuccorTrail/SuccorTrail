package util

import (
	"html/template"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	templates = make(map[string]*template.Template)
	mu        sync.RWMutex
)

func LoadTemplates(directory string, templatesList []string) error {
	mu.Lock()
	defer mu.Unlock()

	for _, tmpl := range templatesList {
		path := filepath.Join(directory, tmpl)
		tmplObj, err := template.ParseFiles(path)
		if err != nil {
			logrus.WithError(err).Errorf("Error parsing template: %s", tmpl)
			return err
		}
		templates[tmpl] = tmplObj
	}
	return nil
}

func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	mu.RLock()
	defer mu.RUnlock()

	tmplObj, ok := templates[tmpl]
	if !ok {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	if err := tmplObj.Execute(w, data); err != nil {
		logrus.WithError(err).Errorf("Error rendering template: %s", tmpl)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
