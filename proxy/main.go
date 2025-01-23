package main

import (
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func main() {
	r := chi.NewRouter()

	h := NewHandler()

	r.Use(middlewareCORS, redirectMiddleware)
	r.Get("/swagger", h.swaggerUI)
	r.Route("/swagger", func(r chi.Router) {
		r.Get("/", h.swaggerUI)
		r.Get("/swagger.yaml", h.swaggerGET)
	})
	r.Route("/api", func(r chi.Router) {
		r.Route("/address", func(r chi.Router) {
			r.Post("/search", h.search)
			r.Post("/geocode", h.geocode)
		})
	})

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}

const content = ``

func NewProxy(targetURL string) *httputil.ReverseProxy {

	target, err := url.Parse(targetURL)
	if err != nil {
		log.Fatal(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	proxy.Director = func(req *http.Request) {
		// Меняем параметры запроса
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Host = target.Host
	}

	return proxy
}

func redirectMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		checkApi := strings.Contains(r.RequestURI, "/api")
		checkSwagger := strings.Contains(r.RequestURI, "/swagger")

		if !checkApi && !checkSwagger {
			proxy := NewProxy("http://hugo:1313")

			proxy.ServeHTTP(w, r)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func middlewareCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)

	})
}
