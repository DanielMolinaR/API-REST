package main

import (
	"TFG/API-REST/lib"
	. "TFG/API-REST/middleware"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func login (w http.ResponseWriter, r *http.Request) {
	lib.TerminalLogger.Trace("User trying to login from: ", r.Host)
	lib.DocuLogger.Trace("User trying to login from: ", r.Host)

	//Read the authorization header
	authHeader := r.Header.Get("Authorization")

	//Extract the Bearer from the data of the header
	token := strings.Replace(authHeader, "Bearer ", "", -1)

	//Check if the token is valid
	if VerifyToken(token) {
		lib.TerminalLogger.Trace("The user has logged in thanks to his token")
		lib.DocuLogger.Trace("The user has logged in thanks to his token")
		setAnswer(map[string]interface{}{"state": "Sesión iniciada gracias al token"}, w, http.StatusAccepted)
	} else {

		if ok, reqBody := readBody(r); !ok {
			setAnswer(map[string]interface{}{"state": "Imposible leer la información"} ,w, http.StatusInternalServerError)
		} else {

			//In UsersLogin we try to login the user in keycloak, if it doesnt exist we receive a false
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
	authHeader := r.Header.Get("Authorization")

	if ok, response := VerifyTokenIsFromAdmin(authHeader); !ok {
		lib.TerminalLogger.Warn("Someone who is not an Amdin is trying to generate an unique URL: ", r.Host)
		lib.DocuLogger.Warn("Someone who is not an Amdin is trying to generate an unique URL: ", r.Host)
		setAnswer(response, w, http.StatusPreconditionFailed)
	} else {
		if ok, reqBody := readBody(r); !ok {
			setAnswer(map[string]interface{}{"state": "Imposible leer la información"} ,w, http.StatusInternalServerError)
		} else {

			//We get the info and verify it, if it's correct the email is sent
			if ok, response := VerifyDataAndSendUniqueEmail(reqBody); !ok{
				setAnswer(response ,w, http.StatusPreconditionFailed)
			} else {
				setAnswer(response ,w, http.StatusOK)
			}
		}
	}
}

func employeeSignUp(w http.ResponseWriter,r *http.Request){
	lib.TerminalLogger.Trace("Signing up an employee from: ", r.Host)
	lib.DocuLogger.Trace("Signing up an employee from: ", r.Host)

	//Read the authorization header
	authHeader := r.Header.Get("Authorization")

	//Extract the Bearer from the data of the header
	//In this request we dont expect a token, we expect the slug from the URL for creating the user
	SignUpUuid := strings.Replace(authHeader, "Bearer ", "", -1)

	//Verify the slug and get the expiration date
	expTime := VerifyUuidAndGetExpTime(SignUpUuid)

	//verify the expiration date
	if expTime == 0{
		lib.TerminalLogger.Error("Empty fields from the table uniqueUrl: ", SignUpUuid, " ", expTime)
		lib.DocuLogger.Error("Empty fields from the table uniqueUrl", SignUpUuid, " ", expTime)
		setAnswer(map[string]interface{}{"state": "El slug no existe"}, w, http.StatusNotAcceptable)
	}else if !VerifyExpTime(expTime){

		//If it's expired we must delete the row from the table and wait for a new one
		DeleteUuidRow(SignUpUuid)
		setAnswer(map[string]interface{}{"state": "Tiempo para crear la cuenta ha expirado"}, w, http.StatusInternalServerError)
	} else{

		//If the slug is correct and it's not expired we sign up the new employee
		if ok, reqBody := readBody(r); !ok {
			setAnswer(map[string]interface{}{"state": "Imposible leer la información"} ,w, http.StatusInternalServerError)
		} else {

			//In EmployeeSignUpVerification the data from the user is verified and if it's correct the user in saved in the DB
			if ok, response := EmployeeSignUpVerification(reqBody); !ok {
				setAnswer(response, w, http.StatusPreconditionFailed)
			} else {

				//If everything went good we delete the SLUG row from the signUp for avoiding duplicate keys
				DeleteUuidRow(SignUpUuid)
				lib.TerminalLogger.Info("Employee created")
				lib.DocuLogger.Info("Employee created")
				setAnswer(response, w, http.StatusCreated)
			}
		}
	}
}

func patientSignUp(w http.ResponseWriter,r *http.Request){
	lib.TerminalLogger.Trace("Verifying an email from: ", r.Host)
	lib.DocuLogger.Trace("Verifying an email from: ", r.Host)

	if ok, reqBody := readBody(r); !ok {
		setAnswer(map[string]interface{}{"state": "Imposible leer la información"} ,w, http.StatusInternalServerError)
	} else {

		//In PatientSignUpVerification the data from the user is verified and if it's correct the user in saved in the DB
		if	ok, response := PatientSignUpVerification(reqBody); !ok{
			setAnswer(response, w, http.StatusPreconditionFailed)
		} else {
			lib.TerminalLogger.Info("Patient created")
			lib.DocuLogger.Info("Patient created")
			setAnswer(response, w, http.StatusCreated)
		}
	}
}

func verifyEmail(w http.ResponseWriter, r *http.Request){
	lib.TerminalLogger.Info("Signing up an user from: ", r.Host)
	lib.DocuLogger.Info("Signing up an user from: ", r.Host)

	//Read the authorization header
	uuid := r.Header.Get("Authorization")

	//Extract the Bearer from the data of the header
	//In this request we dont expect a token, we expect the slug from the URL for verifying the email of the user
	uuid = strings.Replace(uuid, "Bearer ", "", -1)

	//Verify uuid and get expiration date
	expTime := VerifyUuidAndGetExpTime(uuid)

	//verify expiration date
	if expTime == 0{
		lib.TerminalLogger.Error("Empty fields from the table uniqueUrl: ", uuid, " ", expTime)
		lib.DocuLogger.Error("Empty fields from the table uniqueUrl", uuid, " ", expTime)
		setAnswer(map[string]interface{}{"state": "El slug no existe"}, w, http.StatusNotAcceptable)
	}else if !VerifyExpTime(expTime){

		//If it's expired we must update the expiration date and send the email again
		if !UpdateExpTimeFromUuid(uuid) {
			setAnswer(map[string]interface{}{"state": "Imposible verificar el correo"}, w, http.StatusInternalServerError)
		} else {

			//Once the expiration date is updated we must send the email
			if ok, response := ResendVerificationEmail(uuid); !ok{
				setAnswer(response, w, http.StatusPreconditionFailed)
			} else {
				lib.TerminalLogger.Info("Email sent")
				lib.DocuLogger.Info("Email sent")
				setAnswer(response, w, http.StatusCreated)
			}
		}
	} else {

		//If is not expired the user must be updated for login
		if ok, response := UpdateEnabledUserFromSlug(uuid); !ok{
			setAnswer(response, w, http.StatusPreconditionFailed)
		} else {
			setAnswer(response, w, http.StatusAccepted)
		}
	}
}

func createAppointments(w http.ResponseWriter, r *http.Request){
	lib.TerminalLogger.Trace("Creating an appointment from: ", r.Host)
	lib.DocuLogger.Trace("Creating an appointment from: ", r.Host)

	//Read the authorization header
	authHeader := r.Header.Get("Authorization")

	if ok, response := VerifyTokenIsFromEmployeeOrAdmin(authHeader); !ok{
		setAnswer(response, w, http.StatusPreconditionFailed)
	} else{
		if ok, reqBody := readBody(r); !ok {
			setAnswer(map[string]interface{}{"state": "Imposible leer la información"} ,w, http.StatusInternalServerError)
		} else {
			ok, response := AppointmentMiddleware(reqBody)
			if !ok {
				setAnswer(response, w, http.StatusNotAcceptable)
			} else {
				setAnswer(response, w, http.StatusCreated)
			}
		}
	}
}

func createExercises(w http.ResponseWriter, r *http.Request) {
	lib.TerminalLogger.Trace("Creating an exercise from: ", r.Host)
	lib.DocuLogger.Trace("Creating an exercise from: ", r.Host)

	//Read the authorization header
	authHeader := r.Header.Get("Authorization")

	if ok, response := VerifyTokenIsFromEmployeeOrAdmin(authHeader); !ok{
		setAnswer(response, w, http.StatusPreconditionFailed)
	} else{
		if ok, reqBody := readBody(r); !ok {
			setAnswer(map[string]interface{}{"state": "Imposible leer la información"} ,w, http.StatusInternalServerError)
		} else {
			ok, response := ExerciseMiddleware(reqBody)
			if !ok {
				setAnswer(response, w, http.StatusNotAcceptable)
			} else {
				setAnswer(response, w, http.StatusCreated)
			}
		}
	}
}

func getAppointment(w http.ResponseWriter, r *http.Request){
	lib.TerminalLogger.Trace("Getting appointments from: ", r.Host)
	lib.DocuLogger.Trace("Getting appointments from: ", r.Host)

	//Read the authorization header
	authHeader := r.Header.Get("Authorization")

	//Extract the Bearer from the data of the header
	token := strings.Replace(authHeader, "Bearer ", "", -1)

	//Check if the token is valid
	if !VerifyToken(token) {
		lib.TerminalLogger.Trace("The user is trying to retrieve appointments with an invalid token")
		lib.DocuLogger.Trace("The user is trying to retrieve appointments with an invalid token")
		setAnswer(map[string]interface{}{"state": "Token no valido"}, w, http.StatusNotAcceptable)
	} else {
		if ok, response := GetAppointmentsDataFromDni(token); !ok {
			setAnswer(response, w, http.StatusPreconditionFailed)
		} else {
			setAnswer(response, w, http.StatusAccepted)
		}
	}
}

func getAllAppointments(w http.ResponseWriter, r *http.Request){
	lib.TerminalLogger.Trace("Getting all the appointments from: ", r.Host)
	lib.DocuLogger.Trace("Getting all the appointments from: ", r.Host)

	//Read the authorization header
	authHeader := r.Header.Get("Authorization")

	//Check if the token is valid
	if ok, response := VerifyTokenIsFromAdmin(authHeader); !ok {
		lib.TerminalLogger.Trace("The user is trying to retrieve appointments with an invalid token")
		lib.DocuLogger.Trace("The user is trying to retrieve appointments with an invalid token")
		setAnswer(response, w, http.StatusNotAcceptable)
	} else {
		if ok, response := GetAllAppointmentsData(); !ok {
			setAnswer(response, w, http.StatusPreconditionFailed)
		} else {
			setAnswer(response, w, http.StatusAccepted)
		}
	}
}

func getExercises(w http.ResponseWriter, r *http.Request){
	lib.TerminalLogger.Trace("Getting exercises from: ", r.Host)
	lib.DocuLogger.Trace("Getting exercises from: ", r.Host)

	//Read the authorization header
	authHeader := r.Header.Get("Authorization")

	//Extract the Bearer from the data of the header
	token := strings.Replace(authHeader, "Bearer ", "", -1)

	//Check if the token is valid
	if !VerifyToken(token) {
		lib.TerminalLogger.Trace("The user is trying to retrieve exercises with an invalid token")
		lib.DocuLogger.Trace("The user is trying to retrieve exercises with an invalid token")
		setAnswer(map[string]interface{}{"state": "Token no valido"}, w, http.StatusNotAcceptable)
	} else {
		if ok, response := GetExercisesDataFromDni(token); !ok {
			setAnswer(response, w, http.StatusPreconditionFailed)
		} else {
			setAnswer(response, w, http.StatusAccepted)
		}
	}
}

func deleteAppointment(w http.ResponseWriter, r *http.Request){
	lib.TerminalLogger.Trace("Getting exercises from: ", r.Host)
	lib.DocuLogger.Trace("Getting exercises from: ", r.Host)

	//Read the authorization header
	authHeader := r.Header.Get("Authorization")

	//Extract the Bearer from the data of the header
	token := strings.Replace(authHeader, "Bearer ", "", -1)

	//Check if the token is valid
	if !VerifyToken(token) {
		lib.TerminalLogger.Trace("The user is trying to retrieve exercises with an invalid token")
		lib.DocuLogger.Trace("The user is trying to retrieve exercises with an invalid token")
		setAnswer(map[string]interface{}{"state": "Token no valido"}, w, http.StatusNotAcceptable)
	} else {
		if ok, reqBody := readBody(r); !ok {
			setAnswer(map[string]interface{}{"state": "Imposible leer la información"} ,w, http.StatusInternalServerError)
		} else {
			if ok, response := DeleteAppointmentDataFromDni(token, reqBody); !ok {
				setAnswer(response, w, http.StatusPreconditionFailed)
			} else {
				setAnswer(response, w, http.StatusAccepted)
			}
		}
	}
}

func deleteExercise(w http.ResponseWriter, r *http.Request){
	lib.TerminalLogger.Trace("Getting exercises from: ", r.Host)
	lib.DocuLogger.Trace("Getting exercises from: ", r.Host)

	//Read the authorization header
	authHeader := r.Header.Get("Authorization")

	//Extract the Bearer from the data of the header
	token := strings.Replace(authHeader, "Bearer ", "", -1)

	//Check if the token is valid
	if !VerifyToken(token) {
		lib.TerminalLogger.Trace("The user is trying to retrieve exercises with an invalid token")
		lib.DocuLogger.Trace("The user is trying to retrieve exercises with an invalid token")
		setAnswer(map[string]interface{}{"state": "Token no valido"}, w, http.StatusNotAcceptable)
	} else {
		if ok, reqBody := readBody(r); !ok {
			setAnswer(map[string]interface{}{"state": "Imposible leer la información"} ,w, http.StatusInternalServerError)
		} else {
			if ok, response := DeleteExerciseDataFromDni(token, reqBody); !ok {
				setAnswer(response, w, http.StatusPreconditionFailed)
			} else {
				setAnswer(response, w, http.StatusAccepted)
			}
		}
	}
}

func getClinicalBackground(w http.ResponseWriter, r *http.Request){
	lib.TerminalLogger.Trace("Getting exercises from: ", r.Host)
	lib.DocuLogger.Trace("Getting exercises from: ", r.Host)

	//Read the authorization header
	authHeader := r.Header.Get("Authorization")

	//Extract the Bearer from the data of the header
	token := strings.Replace(authHeader, "Bearer ", "", -1)

	if ok, response := VerifyTokenIsFromEmployeeOrAdmin(token); !ok{
		setAnswer(response, w, http.StatusPreconditionFailed)
	} else {
		if ok, reqBody := readBody(r); !ok {
			setAnswer(map[string]interface{}{"state": "Imposible leer la información"} ,w, http.StatusInternalServerError)
		} else {
			if ok, response := GetClinicalBackgroundMiddleware(reqBody); !ok {
				setAnswer(response, w, http.StatusPreconditionFailed)
			} else {
				setAnswer(response, w, http.StatusAccepted)
			}
		}
	}
}

func updateClinicalBackground(w http.ResponseWriter, r *http.Request){
	lib.TerminalLogger.Trace("Getting exercises from: ", r.Host)
	lib.DocuLogger.Trace("Getting exercises from: ", r.Host)

	//Read the authorization header
	authHeader := r.Header.Get("Authorization")

	//Extract the Bearer from the data of the header
	token := strings.Replace(authHeader, "Bearer ", "", -1)

	if ok, response := VerifyTokenIsFromEmployeeOrAdmin(token); !ok{
		setAnswer(response, w, http.StatusPreconditionFailed)
	} else {
		if ok, reqBody := readBody(r); !ok {
			setAnswer(map[string]interface{}{"state": "Imposible leer la información"} ,w, http.StatusInternalServerError)
		} else {
			if ok, response := UpdateClinicalBackgroundMiddleware(reqBody); !ok {
				setAnswer(response, w, http.StatusPreconditionFailed)
			} else {
				setAnswer(response, w, http.StatusAccepted)
			}
		}
	}
}

func getAllEmployees(w http.ResponseWriter, r *http.Request){
	lib.TerminalLogger.Trace("Getting all the employees from: ", r.Host)
	lib.DocuLogger.Trace("Getting all the employees from: ", r.Host)

	//Read the authorization header
	authHeader := r.Header.Get("Authorization")

	//Check if the token is valid
	if ok, response := VerifyTokenIsFromAdmin(authHeader); !ok {
		lib.TerminalLogger.Trace("Someone is trying to retrieve employee's data with an invalid token")
		lib.DocuLogger.Trace("Someone is trying to retrieve appointments with an invalid token")
		setAnswer(response, w, http.StatusNotAcceptable)
	} else {
		if ok, response := GetAllEmployeesData(); !ok {
			setAnswer(response, w, http.StatusPreconditionFailed)
		} else {
			setAnswer(response, w, http.StatusAccepted)
		}
	}
}

func getAllPatients(w http.ResponseWriter, r *http.Request){
	lib.TerminalLogger.Trace("Getting all the patients from: ", r.Host)
	lib.DocuLogger.Trace("Getting all the patients from: ", r.Host)

	//Read the authorization header
	authHeader := r.Header.Get("Authorization")

	//Check if the token is valid
	if ok, response := VerifyTokenIsFromEmployeeOrAdmin(authHeader); !ok {
		lib.TerminalLogger.Trace("Someone is trying to retrieve employee's data with an invalid token")
		lib.DocuLogger.Trace("Someone is trying to retrieve appointments with an invalid token")
		setAnswer(response, w, http.StatusNotAcceptable)
	} else {
		if ok, response := GetAllPatientsData(); !ok {
			setAnswer(response, w, http.StatusPreconditionFailed)
		} else {
			setAnswer(response, w, http.StatusAccepted)
		}
	}
}

func refreshToken(w http.ResponseWriter, r *http.Request){
	lib.TerminalLogger.Trace("Using a refresh token from: ", r.Host)
	lib.DocuLogger.Trace("Using a refresh token from: ", r.Host)

	//Read the authorization header
	authHeader := r.Header.Get("Authorization")

	if ok, accessToken, refreshToken := GettingNewTokens(authHeader); !ok{
		setAnswer(map[string]interface{}{"state": "Imposible conseguir los token"} ,w, http.StatusInternalServerError)
	} else {
		setAnswer(map[string]interface{}{"accessToken": accessToken, "refreshToken": refreshToken}, w, http.StatusAccepted)
	}
}

func setAnswer(response map[string]interface{}, w http.ResponseWriter, state int){
	w.WriteHeader(state)
	_ = json.NewEncoder(w).Encode(response)
}

func readBody(r *http.Request) (bool, []byte){

	// Convert r.Body into a readable formart
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		lib.TerminalLogger.Error("Impossible to read the data from the request")
		lib.DocuLogger.Error("Impossible to read the data from the request")
		return false, nil
	}
	return true, reqBody
}

func main() {
	//routes
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*", "http://localhost:8081"},
		AllowedMethods: []string{http.MethodDelete, http.MethodGet, http.MethodPost,http.MethodPatch},
		AllowCredentials: true,
		AllowedHeaders: []string{"Access-Control-Allow-Origin", "Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
	})
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/login", login).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/generate-and-send-unique-URL", generateAndSendUniqueUrlForSignUp).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/employee-signUp", employeeSignUp).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/patient-signUp", patientSignUp).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/verify-email", verifyEmail).Methods(http.MethodPatch, http.MethodOptions)
	router.HandleFunc("/create-appointments", createAppointments).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/create-exercises", createExercises).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/get-appointments", getAppointment).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/get-all-appointments", getAllAppointments).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/get-exercises", getExercises).Methods(http.MethodGet, http.MethodOptions)
	/*router.HandleFunc("/update-appointment", updateAppointments).Methods(http.MethodPut, http.MethodOptions)
	router.HandleFunc("/update-exercise", updateExercises).Methods(http.MethodPut, http.MethodOptions)*/
	router.HandleFunc("/delete-appointment", deleteAppointment).Methods(http.MethodDelete, http.MethodOptions)
	router.HandleFunc("/delete-exercise", deleteExercise).Methods(http.MethodDelete, http.MethodOptions)
	router.HandleFunc("/get-clinical-background", getClinicalBackground).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/update-clinical-background", updateClinicalBackground).Methods(http.MethodPatch, http.MethodOptions)
	router.HandleFunc("/get-all-employees", getAllEmployees).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/get-all-patients", getAllPatients).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/refresh-token", refreshToken).Methods(http.MethodGet, http.MethodOptions)



	handler := c.Handler(router)
	//router.Use(mux.CORSMethodMiddleware(handler))
	log.Fatal(http.ListenAndServe("localhost:3000", handler))
}


