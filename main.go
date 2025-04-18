package main

import (
	"html/template"
	"log"
	"net/http"
)

type Page struct {
	Title string
}

func handler(w http.ResponseWriter, r *http.Request) {
	file, err := template.ParseFiles("index.html")
	if err != nil {
		return
	}
	err = file.Execute(w, Page{Title: "lol"})
	if err != nil {
		return
	}
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
