package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/jba/muxpatterns"
)

var tmpl *template.Template = nil

func main() {
	fmt.Println("Hello World!")

	mux := muxpatterns.NewServeMux()
	tmpl, _ = template.ParseGlob("./templates/*.tmpl.html")

	mux.Handle("/", http.FileServer(http.Dir("public")))
	log.Fatal(http.ListenAndServe(":8000", mux))
}
