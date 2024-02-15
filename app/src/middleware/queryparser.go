package middleware

import (
	"context"
	"log"
	"net/http"
)

type QueryKey struct {
	Key string
}

func QueryParamsCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		for k, v := range r.URL.Query() {
			log.Println(k, v)
			ctx = context.WithValue(ctx, QueryKey{k}, v)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
