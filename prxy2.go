package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	log.Printf("Starting...")
	http.HandleFunc("/", ProxyFunc)
	log.Fatal(http.ListenAndServe(":8888", nil))
}

func ProxyFunc(w http.ResponseWriter, r *http.Request) {
	u, _ := url.Parse(r.URL.Scheme + "://" + r.URL.Host)
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ServeHTTP(w, r)
}
