package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

var routes = map[string]string{
    "/service1/": "http://localhost:8081",
    "/service2/": "http://localhost:8082",
}

func main() {
    http.HandleFunc("/", gatewayHandler)
    log.Println("API Gateway listening on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func gatewayHandler(w http.ResponseWriter, r *http.Request) {
    for prefix, target := range routes {
        if strings.HasPrefix(r.URL.Path, prefix) {
            remote, err := url.Parse(target)
            if err != nil {
                http.Error(w, "Bad backend URL", http.StatusInternalServerError)
                return
            }
            proxy := httputil.NewSingleHostReverseProxy(remote)
            proxy.ServeHTTP(w, r)
            return
        }
    }
    http.Error(w, "Not found", http.StatusNotFound)
}