package main

import (
	"TFG/API-REST/src/middleware"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

func Homelink (w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, r.Host)
	fmt.Fprintf(w, "Welcome home!")
}

func Login (w http.ResponseWriter, r *http.Request) {
	// Convert r.Body into a readable formart
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the user dni, email, password and name only in order to update")
	}

	if	!middleware.UsersLogin(reqBody){
		fmt.Fprintf(w, "No se ha podido iniciar sesion")
	} else {
		fmt.Fprint(w, "Sesión iniciada")
	}
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Homelink)
	router.HandleFunc("/login", Login).Methods("GET")
	log.Fatal(http.ListenAndServe("localhost:8080", router))
}
