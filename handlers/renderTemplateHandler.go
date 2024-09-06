package handlers

import (
	"fmt"
	"html/template"
	"net/http"
)

func RenderTemplate(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"./static"+r.URL.Path, "./static/base.html"))
	err := tmpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
}
