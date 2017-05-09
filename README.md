httputil
========

Simple set of middleware and other utility functions to make running a
web server in Go just a little bit more convenient.


Examples
--------

### Log every request:

```go
package main

import "fmt"
import "net/http"
import "github.com/fika-io/httputil"

func main() {
  http.HandleFunc("/hello", RequestHello)
  fmt.Println("Serving on port 8080...")
  http.ListenAndServe(":8080", httputil.Logger(http.DefaultServeMux))
}

func RequestHello(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("Hello World!\n"))
}
```

You will now get log lines whenever a request is made:

```
2017/05/09 12:18:06 [::1]:61963 "GET /hello HTTP/1.1" 200 39 "" "curl/7.51.0"
2017/05/09 12:18:08 [::1]:61964 "GET /bye HTTP/1.1" 404 19 "" "curl/7.51.0"
```

### Get Gzip compression on all requests:

```go
func main() {
  http.HandleFunc("/hello", RequestHello)
  fmt.Println("Serving on port 8080...")
  http.ListenAndServe(":8080", httputil.Logger(httputil.Gzipper(http.DefaultServeMux)))
}
```

**Notes:**

* The middleware will check `Accept-Encoding` before serving gzip
* Make sure to set `Content-Type` first if you‘re calling `WriteHeader`

### Add public cache headers to endpoints:

```go
func main() {
  http.HandleFunc("/hello", RequestHello)
  http.Handle("/s/", httputil.Cacher(168*time.Hour, http.StripPrefix("/s/", http.FileServer(http.Dir("static")))))
  http.Handle("/favicon.ico", httputil.FileWithCache("static/favicon.ico", 168*time.Hour))
  http.Handle("/robots.txt", httputil.FileWithCache("static/robots.txt", 24*time.Hour))
  fmt.Println("Serving on port 8080...")
  http.ListenAndServe(":8080", httputil.Logger(httputil.Gzipper(http.DefaultServeMux)))
}
```

**Notes:**

* This is intended for seamless CDN integration and won‘t do any caching for you
* This adds `Vary: Accept-Encoding` to work with and without gzip

### JSON handlers

```go
func main() {
  http.HandleFunc("/hello", httputil.Handler(RequestHello))
  fmt.Println("Serving on port 8080...")
  http.ListenAndServe(":8080", httputil.Logger(httputil.Gzipper(http.DefaultServeMux)))
}

func RequestHello(r *http.Request) (interface{}, error) {
  name := r.URL.Query().Get("name")
  if name == "" {
    return nil, fmt.Errorf("name is missing")
  }
  return map[string]string{"hello": name}, nil
}
```