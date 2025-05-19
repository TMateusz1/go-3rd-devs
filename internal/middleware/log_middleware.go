package middleware

import (
	"log"
	"net/http"
)

func LogMiddleware(fn func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		fn(w, r)
	}

}
