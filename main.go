package main

import (
	. "TFG/lib"
	"TFG/middleware"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func Login (w http.ResponseWriter, r *http.Request) {
	TerminalLogger.Info("An user is trying to login")
	DocuLogger.Info("An user is trying to login")
	fmt.Println(r.Method, r.Host)

	//Read the authorization header
	tokenString := r.Header.Get("Authorization")

	//Extract the Bearer expecification from the data of the header
	token := strings.Replace(tokenString, "Bearer ", "", -1)

	//Check if the token is valid
	if (len(token)!=0){ //if token no vacio se valida el token
		TerminalLogger.Info("The user has logged in thanks to his token")
		DocuLogger.Info("The user has logged in thanks to his token")
		response := map[string]interface{}{"state": "Sesión iniciada por token"}
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(response)
	} else {

		// Convert r.Body into a readable formart
		reqBody, err := ioutil.ReadAll(r.Body)
		if err == nil {

			//In UsersLogin the user data is verified
			if bool, response := middleware.UsersLogin(reqBody); !bool {
				w.WriteHeader(http.StatusPreconditionFailed)
				json.NewEncoder(w).Encode(response)
			} else {
				w.WriteHeader(http.StatusAccepted)
				json.NewEncoder(w).Encode(response)
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err)
		}
	}
}

func Appointments(w http.ResponseWriter, r *http.Request){

}

func employeeSignIn(w http.ResponseWriter,r *http.Request){
	fmt.Println(r.Method, r.Host)
	// Convert r.Body into a readable formart
	reqBody, err := ioutil.ReadAll(r.Body)
	if err == nil {
		//In EmployeeSignInVerification the user
		//is created with data and verified
		if	bool, response := middleware.EmployeeSignInVerification(reqBody); !bool{
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
		if	bool, response := middleware.PatientSignInVerification(reqBody); !bool{
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
	/*c := cors.New(cors.Options{
		AllowedOrigins: []string{"*", "http://localhost:8080"},
		AllowCredentials: true,
		AllowedHeaders: []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
	})*/
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/login", Login).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/appointments", Appointments).Methods("GET", "OPTIONS")
	router.HandleFunc("/employee-signIn", employeeSignIn).Methods("POST")
	router.HandleFunc("/patient-signIn", patientSignIn).Methods("POST")
	//handler := c.Handler(router)
	router.Use(mux.CORSMethodMiddleware(router))
	log.Fatal(http.ListenAndServe("localhost:3000", router))
}


