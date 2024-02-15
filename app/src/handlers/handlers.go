package handlers

import (
	"com/mapify/components"
	"com/mapify/middleware"
	"log"
	"net/http"

	"github.com/a-h/templ"
)

func HelloHandlerFunc(defaultName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var comp templ.Component
		ctx := r.Context()

		log.Println("asdsa", ctx.Value(middleware.QueryKey{Key: "test"}))

		name, ok := ctx.Value(middleware.QueryKey{Key: "test"}).([]string)
		log.Println(name, ok)
		if !ok {
			name, err := Read[string]("name")
			if err != nil {
				comp = components.Hello(defaultName)
			} else {
				comp = components.Hello(*name)
			}
		} else {
			comp = components.Hello(name[0])
		}

		hello := templ.Handler(comp)
		hello.ServeHTTP(w, r)
	}
}
