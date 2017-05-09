package httputil

import (
	"fmt"
	"net/http"
	"time"
)

// Cacher returns a public Cache-Control header for all requests.
func Cacher(duration time.Duration, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()
		header.Set("Cache-Control", fmt.Sprintf("public, max-age=%d", int(duration.Seconds())))
		header.Set("Vary", "Accept-Encoding")
		h.ServeHTTP(w, r)
	})
}

// FileWithCache will serve the file at the provided path with a public Cache-Control header.
func FileWithCache(path string, duration time.Duration) http.Handler {
	return Cacher(duration, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path)
	}))
}
