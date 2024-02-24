package handlers

import (
	"com/mapify/structs"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

func getContact(id string) (*structs.Contact, error) {
	contacts, err := Read[[]structs.Contact]("contacts")
	if err == nil {
		for _, c := range *contacts {
			if c.Id == id {
				return &c, nil
			}
		}
	}
	return nil, errors.New("ContactNotFound")
}

func internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Println(err.Error())
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, "500 Internal Server Error")
}

func notFoundError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "404 Not Found")
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	// http.Redirect(w, r, "localhost:3000/contacts", http.StatusTemporaryRedirect)
	tmpl, err := template.ParseFiles("template/index.html")

	if err != nil {
		internalServerError(w, r, err)
		return
	}
	tmpl.ExecuteTemplate(w, "index.html", map[string]any{
		"title": "Home",
	})
}

func ContactsPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("template/index.html", "template/contacts.html")
	if err != nil {
		internalServerError(w, r, err)
		return
	}
	contacts, err := Read[[]structs.Contact]("contacts")
	if err != nil {
		internalServerError(w, r, err)
		return
	}
	tmpl.ExecuteTemplate(w, "contacts.html", map[string]any{
		"contacts": contacts,
	})
}

func ContactPage(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	elems := strings.Split(path, "/")
	if elems[3] == "edit" {
		editContactPage(w, r, elems[2])
	} else if elems[3] == "view" {
		viewContactPage(w, r, elems[2])
	}
}

func editContactPage(w http.ResponseWriter, r *http.Request, id string) {
	tmpl, err := template.ParseFiles("template/index.html", "template/editcontact.html")
	if err != nil {
		internalServerError(w, r, err)
		return
	}
	contact, err := getContact(id)
	if err != nil {
		internalServerError(w, r, err)
		return
	}
	tmpl.Funcs(template.FuncMap{
		"arr": func(args ...any) []any { return args },
	})
	log.Println(tmpl.DefinedTemplates())
	tmpl.ExecuteTemplate(log.Writer(), "editcontact.html", map[string]any{"contact": contact})
	tmpl.ExecuteTemplate(w, "editcontact.html", map[string]any{
		"contact": contact,
	})
}

func viewContactPage(w http.ResponseWriter, r *http.Request, id string) {
	tmpl, err := template.ParseFiles("template/index.html", "template/viewcontact.html")
	if err != nil {
		internalServerError(w, r, err)
		return
	}
	contact, err := getContact(id)
	if err != nil {
		internalServerError(w, r, err)
		return
	}
	tmpl.Funcs(template.FuncMap{
		"arr": func(args ...any) []any { return args },
	})
	tmpl.ExecuteTemplate(w, "viewcontact.html", map[string]any{
		"contact": contact,
	})
}

func StaticFiles(w http.ResponseWriter, r *http.Request) {
	files, _ := filepath.Glob("static/*")

	for _, file := range files {
		if strings.HasSuffix(r.URL.Path, file) {
			f, err := os.ReadFile(file)
			if err != nil {
				internalServerError(w, r, err)
				return
			}
			w.Write(f)
			return
		}
	}
	notFoundError(w, r)
}
