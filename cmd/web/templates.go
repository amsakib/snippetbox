package main

import (
	"html/template"
	"path/filepath"
	"snippetbox.amsakib.com/internal/models"
	"time"
)

type templateData struct {
	CurrentYear int
	Snippet     models.Snippet
	Snippets    []models.Snippet
	Form        any
	Flash       string
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		// get filename from path
		name := filepath.Base(page)
		// parse the base file
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/pages/base.html")
		if err != nil {
			return nil, err
		}
		// parse all partials
		ts, err = ts.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}
		// parse the current page
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}
