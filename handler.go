package httputil

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// HandlerFunc is a function that takes an HTTP request and returns a (data, err) tuple.
type HandlerFunc func(r *http.Request) (interface{}, error)

// Error returns an error object for the provided HTTP status code.
func Error(code int) error {
	return httpError{code}
}

// ErrorHandler will always return an error with the provided HTTP status code.
func ErrorHandler(code int) http.HandlerFunc {
	return Handler(func(r *http.Request) (interface{}, error) {
		return nil, Error(code)
	})
}

// Handler turns a function that returns a (data, err) tuple into a http.HandlerFunc.
func Handler(f HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		v, err := f(r)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err != nil {
			writeError(w, err)
			return
		}
		b, err := json.Marshal(v)
		if err != nil {
			writeError(w, err)
			return
		}
		_, err = w.Write(b)
		if err != nil {
			writeError(w, err)
		}
	}
}

type httpError struct {
	code int
}

func (e httpError) Error() string {
	return fmt.Sprintf("HTTP %d", e.code)
}

func writeError(w http.ResponseWriter, err error) {
	if httpErr, ok := err.(httpError); ok {
		w.WriteHeader(httpErr.code)
	} else {
		log.Printf("Unknown non-HTTP error %v", err)
		w.WriteHeader(500)
	}
	payload := map[string]string{"error": err.Error()}
	data, _ := json.Marshal(payload)
	w.Write(data)
}
