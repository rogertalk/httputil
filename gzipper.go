package httputil

import (
	"compress/gzip"
	"log"
	"net/http"
	"strings"
	"sync"
)

// Compresses the response.
func Gzipper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			h.ServeHTTP(w, r)
			return
		}
		g := pool.Get().(*gzipResponseWriter)
		g.status = 0
		g.wroteHeader = false
		g.ResponseWriter = w
		g.w.Reset(w)
		defer func() {
			// g.w.Close will write a footer even if no data has been written.
			// StatusNotModified and StatusNoContent expect an empty body so don't close it.
			if g.status != http.StatusNotModified && g.status != http.StatusNoContent {
				if err := g.w.Close(); err != nil {
					log.Printf("ERROR: %v", err)
				}
			}
			pool.Put(g)
		}()
		h.ServeHTTP(g, r)
	})
}

// Writes gzip compressed data (used by Gzipper).
type gzipResponseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
	w           *gzip.Writer
}

var (
	pool = sync.Pool{
		New: func() interface{} {
			w, _ := gzip.NewWriterLevel(nil, gzip.BestSpeed)
			return &gzipResponseWriter{w: w}
		},
	}
)

func (g *gzipResponseWriter) Write(b []byte) (int, error) {
	h := g.Header()
	if _, ok := h["Content-Type"]; !ok {
		h.Set("Content-Type", http.DetectContentType(b))
	}
	if !g.wroteHeader {
		g.WriteHeader(http.StatusOK)
	}
	return g.w.Write(b)
}

func (g *gzipResponseWriter) WriteHeader(code int) {
	g.wroteHeader = true
	g.status = code
	if g.status != http.StatusNotModified && g.status != http.StatusNoContent {
		h := g.Header()
		h.Del("Content-Length")
		h.Set("Content-Encoding", "gzip")
	}
	g.ResponseWriter.WriteHeader(code)
}

func (g *gzipResponseWriter) Flush() {
	if g.w != nil {
		g.w.Flush()
	}
	if f, ok := g.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}
