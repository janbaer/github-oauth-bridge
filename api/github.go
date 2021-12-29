package api

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func Handle(w http.ResponseWriter, req *http.Request) {
	log.Println(req.URL.Path)

	setupResponse(&w, req)

	if (*req).Method == http.MethodOptions {
		return
	}

	newPath := strings.Replace(req.URL.Path, "/api/github", "", -1)
	url, _ := url.Parse(fmt.Sprintf("https://api.github.com/%s", newPath))
	log.Printf("Redirect to %s\n", url)

	proxy := httputil.NewSingleHostReverseProxy(url)

	req.Host = "api.github.com"
	req.URL.Host = "api.github.com"
	req.URL.Scheme = "https"

	proxy.ServeHTTP(w, req)
}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", req.URL.Hostname())
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
