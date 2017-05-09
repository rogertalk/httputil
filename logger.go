package httputil

import (
	"fmt"
	"log"
	"net/http"
)

// Logger logs incoming requests, including response status.
func Logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		info := fmt.Sprintf("%s %s %s", r.Method, r.URL, r.Proto)
		o := &responseObserver{ResponseWriter: w}
		h.ServeHTTP(o, r)
		log.Printf("%s %q %d %d %q %q",
			r.RemoteAddr,
			info,
			o.status,
			o.written,
			r.Referer(),
			r.UserAgent())
	})
}

// Spies on http.ResponseWriter (used by Logger).
type responseObserver struct {
	http.ResponseWriter
	status      int
	written     int64
	wroteHeader bool
}

func (o *responseObserver) Write(p []byte) (n int, err error) {
	if !o.wroteHeader {
		o.WriteHeader(http.StatusOK)
	}
	n, err = o.ResponseWriter.Write(p)
	o.written += int64(n)
	return
}

func (o *responseObserver) WriteHeader(code int) {
	o.wroteHeader = true
	o.status = code
	o.ResponseWriter.WriteHeader(code)
}
