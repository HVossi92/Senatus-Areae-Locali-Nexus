package main

import (
	"fmt"
	"log"
	"net/http"
	"senatus/src/templates"

	"github.com/a-h/templ"
)

func main() {
	// Serve static files from the "static" directory
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	component := templates.Index()
	http.Handle("/", templ.Handler(component))

	fmt.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
