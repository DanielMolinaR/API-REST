package main

import (
	. "TFG/API-REST/src/middleware"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"io/ioutil"
	"log"
	"net/http"
)

func Login (w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, r.Host)

	//Read the authorization header
	tokenString := r.Header.Get("Authorization")

	// Convert r.Body into a readable formart
	reqBody, err := ioutil.ReadAll(r.Body)
	if err == nil {

		//In UsersLogin the user data is verified
		if	bool, response := UsersLogin(reqBody, tokenString); !bool{
			json.NewEncoder(w).Encode(response)
		} else {
			json.NewEncoder(w).Encode(response)
		}
	} else{
		json.NewEncoder(w).Encode(err)
	}
}

func employeeSignIn(w http.ResponseWriter,r *http.Request){
	fmt.Println(r.Method, r.Host)
	// Convert r.Body into a readable formart
	reqBody, err := ioutil.ReadAll(r.Body)
	if err == nil {
		//In EmployeeSignInVerification the user
		//is created with data and verified
		if	bool, response := EmployeeSignInVerification(reqBody); !bool{
			response = "No se ha podido crear el usuario: " + response
			json.NewEncoder(w).Encode(response)
		} else {
			fmt.Fprintf(w, response)
		}
	} else{
		json.NewEncoder(w).Encode(err)
	}
}

func patientSignIn(w http.ResponseWriter,r *http.Request){
	fmt.Println(r.Method, r.Host)
	// Convert r.Body into a readable formart
	reqBody, err := ioutil.ReadAll(r.Body)
	if err == nil {
		//In PatientSignInVerification the user
		//is created with data and verified
		if	bool, response := PatientSignInVerification(reqBody); !bool{
			fmt.Fprintf(w, "No se ha podido crear el usuario: %v", response)
		} else {
			fmt.Fprintf(w, response)
		}
	} else{
		json.NewEncoder(w).Encode(err)
	}
}

func main() {
	//routes
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Login).Methods("POST", "OPTIONS")
	router.HandleFunc("/employee-signIn", employeeSignIn).Methods("POST")
	router.HandleFunc("/patient-signIn", patientSignIn).Methods("POST")
	handler := cors.Default().Handler(router)
	log.Fatal(http.ListenAndServe("localhost:3000", handler))
}
