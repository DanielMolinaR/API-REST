package main

import (
	. "TFG/API-REST/src/middleware"
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
	fmt.Println(r.Method, r.Host)
	// Convert r.Body into a readable formart
	reqBody, err := ioutil.ReadAll(r.Body)
	if err == nil {
		//In UsersLogin the user is created with data
		//and then It's verified
		if	bool, response := UsersLogin(reqBody); !bool{
			fmt.Fprintf(w, response)
		} else {
			fmt.Fprintf(w, response)
		}
	} else{
		panic(err)
	}
}

func employeeSignIn(w http.ResponseWriter,r *http.Request){
	fmt.Println(r.Method, r.Host)
	// Convert r.Body into a readable formart
	reqBody, err := ioutil.ReadAll(r.Body)
	if err == nil {
		//In UsersLogin the user is created with data
		//and then It's verified
		if	bool, response := EmployeeSignInVerification(reqBody); !bool{
			fmt.Fprintf(w, "No se ha podido crear el usuario: %v", response)
		} else {
			fmt.Fprintf(w, response)
		}
	} else{
		panic(err)
	}
}

func main() {
	//routes
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Homelink)
	router.HandleFunc("/login", Login).Methods("GET")
	router.HandleFunc("/employee-signIn", employeeSignIn).Methods("POST")
	log.Fatal(http.ListenAndServe("localhost:8080", router))
}
