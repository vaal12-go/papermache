package handlers

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
)

var Content embed.FS
var UseEmbeddedFS = false

func RenderTemplate(w http.ResponseWriter, r *http.Request) {
	var tmpl *template.Template
	if UseEmbeddedFS {
		// fSystem := fs.FS(content)
		// staticFS, err := fs.Sub(fSystem, "static")
		// if err != nil {
		// 	fmt.Printf("Error getting subdir: %v\n", err)
		// }
		// fmt.Printf("staticFS: %v\n", staticFS)
		// dr, err := staticFS.ReadDir(".")
		// if err != nil {
		// 	fmt.Printf("Have error reading file system:%s\n", err)
		// }
		// fmt.Printf("Content of embedded:%s\n", dr)

		// dr, err := Content.ReadDir("static")
		// if err != nil {
		// 	fmt.Printf("Have error reading file system:%s\n", err)
		// }
		// // fmt.Printf("Content of embedded2:%s\n", dr)

		// dr, err = Content.ReadDir("/static")
		// if err != nil {
		// 	fmt.Printf("Have error reading file system:%s\n", err)
		// }
		// fmt.Printf("Content of embedded3:%s\n", dr)

		// fmt.Printf("r.URL.Path: %v\n", r.URL.Path)

		// https: //blog.jetbrains.com/go/2021/06/09/how-to-use-go-embed-in-go-1-16/
		tmpl = template.Must(template.ParseFS(Content,
			"static"+r.URL.Path, "static/base.html"))
	} else {
		tmpl = template.Must(template.ParseFiles(
			"./static"+r.URL.Path, "./static/base.html"))
	}
	err := tmpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
}
