package main

import (
	"TFG/API-REST/lib"
	. "TFG/API-REST/middleware"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func login (w http.ResponseWriter, r *http.Request) {
	lib.TerminalLogger.Info("An user is trying to login from: ", r.Host)
	lib.DocuLogger.Info("An user is trying to login from: ", r.Host)

	//Read the authorization header
	tokenString := r.Header.Get("Authorization")

	//Extract the Bearer from the data of the header
	token := strings.Replace(tokenString, "Bearer ", "", -1)

	//Check if the token is valid
	if Verify(token){
		lib.TerminalLogger.Info("The user has logged in thanks to his token")
		lib.DocuLogger.Info("The user has logged in thanks to his token")
		response := map[string]interface{}{"state": "Sesión iniciada gracias al token"}
		w.WriteHeader(http.StatusAccepted)
		_ = json.NewEncoder(w).Encode(response)
	} else {
		// Convert r.Body into a readable formart
		reqBody, err := ioutil.ReadAll(r.Body)
		if err == nil {
			//In UsersLogin the user data is verified
			if ok, response := UsersLogin(reqBody); !ok {
				w.WriteHeader(http.StatusPreconditionFailed)
				_ = json.NewEncoder(w).Encode(response)
			} else {
				w.WriteHeader(http.StatusAccepted)
				_ = json.NewEncoder(w).Encode(response)
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(err)
		}
	}
}

func Appointments(w http.ResponseWriter, r *http.Request){

}

func employeeSignUp(w http.ResponseWriter,r *http.Request){
	lib.TerminalLogger.Info("Signing up an user from: ", r.Host)
	lib.DocuLogger.Info("Signing up an user from: ", r.Host)

	//Read the authorization header
	tokenString := r.Header.Get("Authorization")

	//Extract the Bearer from the data of the header
	token := strings.Replace(tokenString, "Bearer ", "", -1)

	if (Verify(token) && GetTheRole(DecodeToken(token))==2){
		// Convert r.Body into a readable formart
		reqBody, err := ioutil.ReadAll(r.Body)
		if err == nil {
			//In EmployeeSignInVerification the user is created with verified data
			if ok, response := EmployeeSignInVerification(reqBody); !ok {
				_ = json.NewEncoder(w).Encode(response)
			} else {
				lib.TerminalLogger.Info("User created")
				lib.DocuLogger.Info("User created")
				w.WriteHeader(http.StatusAccepted)
				_ = json.NewEncoder(w).Encode(response)
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(err)
		}
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
	/*c := cors.New(cors.Options{
		AllowedOrigins: []string{"*", "http://localhost:8080"},
		AllowCredentials: true,
		AllowedHeaders: []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
	})*/
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/login", login).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/appointments", Appointments).Methods("GET", "OPTIONS")
	router.HandleFunc("/employee-signUp", employeeSignUp).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/patient-signIn", patientSignIn).Methods("POST")
	//handler := c.Handler(router)
	router.Use(mux.CORSMethodMiddleware(router))
	log.Fatal(http.ListenAndServe("localhost:3000", router))
}


