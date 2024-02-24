package main

// import (
// 	"html/template"
// 	"log"
// 	"net/http"
// 	"os"
// 	"time"

// 	"com/mapify/handlers"
// 	mymiddleware "com/mapify/middleware"

// 	"github.com/go-chi/chi/v5"
// 	"github.com/go-chi/chi/v5/middleware"
// )

// func main_old() {
// 	r := chi.NewRouter()

// 	handlers.Populate()

// 	tmpl, err := template.ParseGlob("template/**")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	style, err := os.ReadFile("style.css")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// A good base middleware stack
// 	r.Use(middleware.RequestID)
// 	r.Use(middleware.RealIP)
// 	r.Use(middleware.Logger)
// 	r.Use(middleware.Recoverer)

// 	// Set a timeout value on the request context (ctx), that will signal
// 	// through ctx.Done() that the request has timed out and further
// 	// processing should be stopped.
// 	r.Use(middleware.Timeout(60 * time.Second))

// 	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
// 		tmpl.ExecuteTemplate(w, "index.html", "mapify")
// 	})

// 	r.Get("/style.css", func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Add("Content-Type", "text/css")
// 		w.Write(style)
// 	})

// 	r.Route("/contacts", func(r chi.Router) {
// 		r.Get("/", handlers.IndexWrapper(tmpl, handlers.ContactsHandlerFunc))

// 		r.Post("/new", handlers.IndexWrapper(tmpl, handlers.ContactsNewHandlerFunc))

// 		r.Get("/details/{id}", handlers.IndexWrapper(tmpl, handlers.ContactsDetailsHandlerFunc))

// 		r.Get("/{id}/edit", handlers.IndexWrapper(tmpl, handlers.ContactsEditHandlerFunc))

// 		r.Post("/{id}/edit", handlers.IndexWrapper(tmpl, nil))

// 		r.Post("/{id}/delete", handlers.IndexWrapper(tmpl, nil))
// 	})

// 	r.Route("/components", func(r chi.Router) {
// 		r.Use(mymiddleware.QueryParamsCtx)

// 	})

// 	http.ListenAndServe(":3000", r)
// }
