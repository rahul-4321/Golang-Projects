package main

import (
	"html/template"
	"log"
	"net/http"
	"https://github.com/rahul-4321/Golang-Projects/UnitConverter/converter"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/length", http.StatusSeeOther)
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/length", converter.lengthHandler)
	http.HandleFunc("/weight", converter.weightHandler)
	http.HandleFunc("/temp", converter.tempHandler)

	log.Println("Server started at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
