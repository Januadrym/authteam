package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type (
	Router struct {
		m map[string]http.Handler
	}
)

func NewRouter() *Router {
	return &Router{
		m: make(map[string]http.Handler),
	}
}

func (router *Router) Handle(path string, method string, handler http.Handler) {
	router.m[path+method] = handler
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler, ok := router.m[r.URL.Path+r.Method]
	if ok {
		handler.ServeHTTP(w, r)
		return
	}
	http.Error(w, "not found", http.StatusNotFound)
}

func main() {
	router := mux.NewRouter()
	router.Use(LogMiddleware)

	router.Path("/api/v1/users").Methods(http.MethodGet).HandlerFunc(Handler1)
	router.Path("/api/v1/users").Methods(http.MethodPost).HandlerFunc(Handler2)

	server := http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func Handler1(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"name": "jack",
		"age":  22,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Handler2(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("handler2"))
}

func LogMiddleware(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("path: %s, method: %s\n", r.URL.Path, r.Method)
		inner.ServeHTTP(w, r)
	})
}
