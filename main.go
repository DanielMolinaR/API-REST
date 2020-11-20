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
	if len(token)>0 {
		if !VerifyToken(token) {
			setAnswer(map[string]interface{}{"state": "Token no válido"}, w, http.StatusNotAcceptable)
		} else {

			lib.TerminalLogger.Trace("The user has logged in thanks to his token")
			lib.DocuLogger.Trace("The user has logged in thanks to his token")
			setAnswer(map[string]interface{}{"state": "Sesión iniciada gracias al token"}, w, http.StatusAccepted)
		}
	} else {
		// Convert r.Body into a readable formart
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			lib.TerminalLogger.Error("Impossible to read the data from the request")
			lib.DocuLogger.Error("Impossible to read the data from the request")
			setAnswer(map[string]interface{}{"state": "Imposible leer la información"} ,w, http.StatusInternalServerError)
		} else {
			//In UsersLogin the user data is verified
			if ok, response := UsersLogin(reqBody); !ok {
				setAnswer(response, w, http.StatusPreconditionFailed)
			} else {
				setAnswer(response, w, http.StatusAccepted)
			}
		}
	}
}

func generateAndSendUniqueUrlForSignUp (w http.ResponseWriter, r *http.Request) {
	//Read the authorization header
	tokenString := r.Header.Get("Authorization")

	//Extract the Bearer from the data of the header
	token := strings.Replace(tokenString, "Bearer ", "", -1)

	//Check if the token is valid
	if !VerifyToken(token){
		setAnswer(map[string]interface{}{"state": "Token no válido"} ,w, http.StatusNotAcceptable)
	} else if GetTheRole(token)!= 2{
		lib.TerminalLogger.Warn("Someone who is not an Amdin is trying to generate an unique URL")
		lib.DocuLogger.Warn("Someone who is not an Amdin is trying to generate an unique URL")
		setAnswer(map[string]interface{}{"state": "Acceso restringido"} ,w, http.StatusNotAcceptable)
	} else {
		// Convert r.Body into a readable formart
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			lib.TerminalLogger.Error("Impossible to read the data from the JSON")
			lib.DocuLogger.Error("Impossible to read the data from the JSON")
			setAnswer(map[string]interface{}{"state": "Imposible leer la información"} ,w, http.StatusInternalServerError)
		} else {
			if ok, response := GenerateAndSendURL(reqBody); !ok{
				setAnswer(response ,w, http.StatusPreconditionFailed)
			} else {
				setAnswer(response ,w, http.StatusOK)
			}
		}
	}
}

func employeeSignUp(w http.ResponseWriter,r *http.Request){
	lib.TerminalLogger.Info("Signing up an user from: ", r.Host)
	lib.DocuLogger.Info("Signing up an user from: ", r.Host)

	//Read the authorization header
	uuid := r.Header.Get("Authorization")

	//Extract the Bearer from the data of the header
	//In this reequest we dont expect a token, we expect the slug from the URL
	uuid = strings.Replace(uuid, "Bearer ", "", -1)

	//Verify the slug and check if it has not expired
	if (!VerifyUuid(uuid)) {
		setAnswer(map[string]interface{}{"state": "SLug no válido"}, w, http.StatusInternalServerError)
	} else{
		// Convert r.Body into a readable formart
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			lib.TerminalLogger.Error("Impossible to read the data from the JSON")
			lib.DocuLogger.Error("Impossible to read the data from the JSON")
			setAnswer(map[string]interface{}{"state": "Imposible leer la información"} ,w, http.StatusInternalServerError)
		} else {
			//In EmployeeSignInVerification the user is created with verified data
			if ok, response := EmployeeSignInVerification(reqBody); !ok {
				setAnswer(response, w, http.StatusPreconditionFailed)
				_ = json.NewEncoder(w).Encode(response)
			} else {
				lib.TerminalLogger.Info("Employee created")
				lib.DocuLogger.Info("Employee created")
				setAnswer(response, w, http.StatusCreated)
			}
		}
	}
}

func patientSignUp(w http.ResponseWriter,r *http.Request){
	fmt.Println(r.Method, r.Host)
	// Convert r.Body into a readable formart
	reqBody, err := ioutil.ReadAll(r.Body)
	if err == nil {
		//In PatientSignInVerification the user is created with data and verified
		if	bool, response := PatientSignInVerification(reqBody); !bool{
			_ = json.NewEncoder(w).Encode(response)
		} else {
			lib.TerminalLogger.Info("Patient created")
			lib.DocuLogger.Info("Patient created")
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(response)
		}
	} else{
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err)
	}
}


func Appointments(w http.ResponseWriter, r *http.Request){

}

func setAnswer(response map[string]interface{}, w http.ResponseWriter, state http.ConnState){
	w.WriteHeader(http.StatusPreconditionFailed)
	_ = json.NewEncoder(w).Encode(response)
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
	router.HandleFunc("/generate-and-send-unique-URL", generateAndSendUniqueUrlForSignUp).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/employee-signUp", employeeSignUp).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/appointments", Appointments).Methods("GET", "OPTIONS")
	router.HandleFunc("/patient-signIn", patientSignUp).Methods("POST")
	//handler := c.Handler(router)
	router.Use(mux.CORSMethodMiddleware(router))
	log.Fatal(http.ListenAndServe("localhost:3000", router))
}


