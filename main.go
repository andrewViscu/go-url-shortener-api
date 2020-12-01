package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"math/rand"
	"log"
	"time"
)

func ConfigureServer() *mux.Router{
	r := mux.NewRouter()
	r.HandleFunc("/", Index).Methods("GET")
	r.HandleFunc("/", ShortenURL).Methods("POST")
	r.HandleFunc("/{shorten_url}", RedirectFromURL).Methods("POST")
	return r
}

func main() {
	rand.Seed(time.Now().UnixNano()) // set seed for different url outcomes
	handler := ConfigureServer()
	server := &http.Server{
		Addr: ":4567",
		Handler: handler,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}

}