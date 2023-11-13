package main

import (
	"html/template"
	"path/filepath"
	"time"

	"lincoln.boris/forum/pkg/forms"
	"lincoln.boris/forum/pkg/models"
)

type templateData struct {
	AuthenticatedUser int
	CurrentYear int
	ErrorMessage string
	Form *forms.Form
	FileTooLargeErr string
	Post *models.Post
	PostCategories []*models.PostCategory
	Comments []*models.Comment
	Posts []*models.Post
	Category *models.Category
	Categories []*models.Category
}

func humanDate(t time.Time) string {
	return t.Format("02.01.2006, 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}

	return cache, nil
}
