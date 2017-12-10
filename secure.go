package httputil

import (
	"fmt"
	"net/http"
)

// SecureEnforcer ensures HTTPS in environments that provide the X-Forwarded-Proto header.
func SecureEnforcer(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Forwarded-Proto") == "http" {
			http.Redirect(w, r, fmt.Sprintf("https://%s%s", r.Host, r.RequestURI), http.StatusFound)
			return
		}
		h.ServeHTTP(w, r)
	})
}

// SecureHostEnforcer ensures HTTPS in environments that provide the
// X-Forwarded-Proto header and also ensures the host is correct.
func SecureHostEnforcer(host string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Host != host || r.Header.Get("X-Forwarded-Proto") == "http" {
			http.Redirect(w, r, fmt.Sprintf("https://%s%s", host, r.RequestURI), http.StatusFound)
			return
		}
		h.ServeHTTP(w, r)
	})
}
