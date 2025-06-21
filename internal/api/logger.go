package api

import (
	"log"
	"net/http"
	"time"
)

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Printf("Start processing request %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		log.Printf("Request %s %s took %v", r.Method, r.URL.Path, time.Since(start))
	})
}
