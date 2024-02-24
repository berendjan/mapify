package handlers

import (
	"com/mapify/components"
	"com/mapify/structs"
	"context"
	"errors"
	"html/template"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
)

// func HelloHandlerFunc(defaultName string) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var comp templ.Component
// 		ctx := r.Context()

// 		name, ok := ctx.Value(middleware.QueryKey{Key: "test"}).([]string)
// 		if !ok {
// 			name, err := Read[string]("name")
// 			if err != nil {
// 				comp = components.Hello(defaultName)
// 			} else {
// 				comp = components.Hello(*name)
// 			}
// 		} else {
// 			comp = components.Hello(name[0])
// 		}

// 		hello := templ.Handler(comp)
// 		hello.ServeHTTP(w, r)
// 	}
// }

func IndexWrapper(tmpl *template.Template, templComp func(w http.ResponseWriter, r *http.Request) (templ.Component, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		comp, err := templComp(w, r)
		if err != nil {
			return
		}
		html, err := templ.ToGoHTML(context.Background(), comp)
		if err != nil {
			log.Fatalf("failed to convert to html: %v", err)
		}

		content := structs.Content{
			Title: "Mapify-Contacts",
			Body:  html,
		}

		tmpl.ExecuteTemplate(w, "index.html", content)
	}
}

func ContactsHandlerFunc(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	contacts, err := Read[[]structs.Contact]("contacts")
	if err != nil {
		contacts = new([]structs.Contact)
	}
	return components.Contacts(*contacts), nil
}

func ContactsNewHandlerFunc(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	return components.ContactsNew(), nil
}

func ContactsDetailsHandlerFunc(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	id := chi.URLParam(r, "id")
	if contact, err := getContact(id); err != nil {
		w.Write([]byte("404 Not Found"))
		return nil, errors.New("NotFound")
	} else {
		return components.ContactsDetails(*contact), nil
	}
}

func ContactsEditHandlerFunc(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	id := chi.URLParam(r, "id")
	if contact, err := getContact(id); err != nil {
		w.Write([]byte("404 Not Found"))
		return nil, errors.New("NotFound")
	} else {
		return components.ContactsEdit(*contact), nil
	}
}
