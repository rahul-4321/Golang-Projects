package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/rahul-4321/Golang-Projects/UnitConverter/converter"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/length", http.StatusSeeOther)
}

func main() {
	converter.InitTemplates(templates)
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/length", converter.LengthHandler)
	http.HandleFunc("/weight", converter.WeightHandler)
	http.HandleFunc("/temp", converter.TempHandler)

	log.Println("Server started at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
