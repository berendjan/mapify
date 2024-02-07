package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Construct the file path based on the URL path
		filePath := "pages" + r.URL.Path
		if r.URL.Path == "/" {
			filePath = "pages/page.html"
		} else {
			filePath += "/page.html"
		}

		// Ensure the file exists
		if _, err := filepath.Abs(filePath); err != nil {
			http.NotFound(w, r)
			return
		}

		// Parse and serve the template
		tmpl, err := template.ParseFiles(filePath)
		if err != nil {
			log.Printf("Error parsing template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			log.Printf("Error executing template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	})

	log.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println(err.Error())
	}
}
