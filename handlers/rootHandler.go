package handlers

import (
	"net/http"
)

func GetRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.String() == "/favicon.ico" {
		http.Redirect(w, r, "static/img/favicon.png", http.StatusFound)
	} else {
		http.Error(w, "Invalid request:"+r.URL.String(), 404)
	}
} //func getRoot(w http.ResponseWriter, r *http.Request) {
