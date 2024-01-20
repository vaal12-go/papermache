package main

import (
	"net/http"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	// fmt.Printf("got / request\n")
	// fmt.Printf("r.URL: %v\n", r.URL)
	http.Redirect(w, r, "static/index.html", http.StatusFound)
	// if r.URL.Path != "/" {
	// 	http.NotFound(w, r)
	// 	return
	// }
	// io.WriteString(w, "Root handler\n")
} //func getRoot(w http.ResponseWriter, r *http.Request) {
