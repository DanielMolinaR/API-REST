package main

import (
	"github.com/gorilla/mux"
	"TFG/API-REST/src/middleware"
	"net/http"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		go middleware.Homelink(w, r)
	})
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		go middleware.Login(w, r)
	})
}
