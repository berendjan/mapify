package main

import (
	"com/mapify/handlers"
	"com/mapify/middleware"
	"log"
	"net/http"
)

func main() {

	handlers.Populate()

	m := middleware.NewMiddleWare(middleware.ContentTypeHeader)

	http.HandleFunc("/", handlers.HomePage)
	http.HandleFunc("/contacts", handlers.ContactsPage)
	http.HandleFunc("/contacts/", handlers.ContactPage)
	http.HandleFunc("/static/", m.Chain(handlers.StaticFiles))

	log.Fatal(http.ListenAndServe("localhost:3000", nil))
}
