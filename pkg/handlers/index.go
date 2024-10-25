package handlers

import (
	"fmt"
	"html/template"
	"net/http"
)

func BuildIndexHandler(opts *Options) (func(http.ResponseWriter, *http.Request), error) {

	tmpl, err := template.ParseFS(templates, "templates/index.html", "templates/base.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %s", err)
	}

	return func(w http.ResponseWriter, r *http.Request) {

		err := tmpl.ExecuteTemplate(w, "base", struct {
			Opts *Options
		}{Opts: opts})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to render template: %s", err), http.StatusInternalServerError)
			return
		}

	}, nil
}
