package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
)

var templates map[string]*template.Template

func init() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	templatesDir := "./static/"

	layouts, err := filepath.Glob(templatesDir + "layout/*.html")
	if err != nil {
		log.Fatal(err)
	}

	widgets, err := filepath.Glob(templatesDir + "widget/*.html")
	if err != nil {
		log.Fatal(err)
	}

	for _, layout := range layouts {
		files := append(widgets, layout)
		templates[filepath.Base(layout)] = template.Must(template.ParseFiles(files...))
	}
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/index", showindex)
	r.HandleFunc("/post", showpost)
	r.HandleFunc("/tag/{id:[0-9]+}", showtag)

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:9000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

func showindex(wrt http.ResponseWriter, req *http.Request) {
	str := []string{"hello", "world", "nba"}
	renderTemplate(wrt, "index.html", str)

	// fmt.Fprintln(wrt, "Hello World")
}

func showpost(wrt http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(wrt, "This is post")
}

func showtag(wrt http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(wrt, "This is tag")
}

func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	tmpl, ok := templates[name]
	if !ok {
		log.Fatal("template not exist")
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.ExecuteTemplate(w, name, data)
}
