package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
)

type page struct {
	Title   string
	Content string
}

func main() {
	r := mux.NewRouter()
	r.PathPrefix("/static").Handler(http.FileServer(http.Dir("./")))
	r.HandleFunc("/{pageName}", viewHandler)
	r.HandleFunc("/", viewHandler)
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pageName := vars["pageName"]

	if pageName == "" {
		pageName = "home-page"
	}
	log.Printf("Serve %s", pageName)

	t, err := template.ParseFiles("templates/base.html")
	if err != nil {
		log.Printf("Error loading template [%s]: [%s]", pageName, err)
		return
	}

	page, err := loadPage(pageName)
	if err != nil {
		return // Page doesn't exist - just skip
	}
	t.Execute(w, page)
}

func loadPage(title string) (*page, error) {
	var thispage page

	content, err := ioutil.ReadFile("pages/" + title + ".html")
	if err != nil {
		return nil, err
	}

	thispage.Content = string(content)
	thispage.Title = strings.Title(strings.Replace(title, "-", " ", -1))

	return &thispage, nil
}
