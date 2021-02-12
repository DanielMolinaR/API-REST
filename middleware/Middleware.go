package middleware

import (
	"TFG/API-REST/lib"
	. "TFG/API-REST/structures"
	"encoding/json"
	"github.com/google/uuid"
	"strconv"
	"strings"
	"time"
)

func UsersLogin(reqBody []byte) (bool, map[string]interface{}) {
	var userToLogIn Users

	//The data from reqBody is filled in the newUser
	err := json.Unmarshal(reqBody, &userToLogIn)

	userToLogIn.DNI = strings.ToLower(userToLogIn.DNI)

	if err != nil{
		lib.TerminalLogger.Error("Impossible to retrieve the data from the JSON")
		lib.DocuLogger.Error("Impossible to retrieve the data from the JSON")
		return false, map[string]interface{}{"state": "Problemas con la lectura de los datos"}

	//Try to login in keycloak
	}else if ok, accessToken, refreshToken := UserCredentialsLogin(userToLogIn.DNI, userToLogIn.Password); !ok{
		return false, map[string]interface{}{"state": "DNI o contraseña incorrecto"}
	} else {
		//Return true with a msg of correct login,
		//the name of the user, the tokens and the role
		lib.TerminalLogger.Trace("User logged in with the DNI: ******", userToLogIn.DNI[6:])
		lib.DocuLogger.Trace("User logged in with the DNI: ******", userToLogIn.DNI[6:])
		_, name := getStringFromField("users", "name", "dni", userToLogIn.DNI)
		return true, map[string]interface{}{"state": "Sesión iniciada", "accessToken": accessToken,
			"refreshToken": refreshToken, "role": getTheRole(accessToken), "email": getEmail(accessToken), "userName": name}
	}
}

func VerifyDataAndSendUniqueEmail(reqBody []byte) (bool, map[string]interface{}){
	var userData Users

	//The data from reqBody is filled in the newUser
	err := json.Unmarshal(reqBody, &userData)

	userData.DNI = strings.ToLower(userData.DNI)

	if err != nil{
		lib.TerminalLogger.Error("Impossible to retrieve the data from the JSON")
		lib.DocuLogger.Error("Impossible to retrieve the data from the JSON")
		return false, map[string]interface{}{"state": "Problemas con la lectura de los datos"}

	//verify if the format of the email is correct
	} else if !verifyEmail(userData.Email){
		lib.TerminalLogger.Error("Email format incorect")
		lib.DocuLogger.Error("Email format incorect")
		return false, map[string]interface{}{"state": "Formato del correo electrónico incorrecto"}
	} else {

		//if it's correct we generate the slug, then we saved it with an expiration date and finally we send the email
		customUuid := generateUUID()
		if !insertUuidExpTimeAndUserId(customUuid, "", userData.Email){
			return false, map[string]interface{}{"state": "Imposible de generar el url unico"}
		} else {
			if ok, response := CreateVerificationEmail(customUuid, userData.Name, userData.Email,"employee-sign-up", 0); !ok{

				//If the email has not been sent we delete the new row of the uuid for avoiding duplicate keys
				DeleteUuidRow(customUuid)
				return ok, response
			} else {
				return ok, response
			}
		}
	}
}

func EmployeeSignUpVerification(reqBody []byte) (bool, map[string]interface{}){
	var newEmployee Employee

	//The data from reqBody is filled in the newUser
	err := json.Unmarshal(reqBody, &newEmployee)

	newEmployee.User.DNI = strings.ToLower(newEmployee.User.DNI)

	newEmployee.Admin = false
	newEmployee.Active = true

	if err != nil{
		lib.TerminalLogger.Error("Impossible to retrieve the data from the JSON")
		lib.DocuLogger.Error("Impossible to retrieve the data from the JSON")
		return false, map[string]interface{}{"state": "Problemas con la lectura de los datos"}
	} else {

		//verify user data
		ok, response := userDataVerifications("employee", newEmployee.User.DNI, newEmployee.User.Phone, newEmployee.User.Email, newEmployee.User.Password)
		if !ok {
			return false, response

		//If all the data is correct, the user is inserted in the DB and Keycloak but disabled for login
		} else if ok, id := doEmployeeInsert(newEmployee); !ok {
			return false, map[string]interface{}{"state": "Imposible añadir el usuario en la BBDD"}
		} else {

			//if the user has been created we must wait for the email verification so we send the email with a verification URL
			EmailUuid := generateUUID()
			if !insertUuidExpTimeAndUserId(EmailUuid, id, newEmployee.User.Email){
				return false, map[string]interface{}{"state": "Imposible de generar el url unico"}
			} else {
				if ok, response := CreateVerificationEmail(EmailUuid, newEmployee.User.Name, newEmployee.User.Email,"email-verification", 1); !ok{

					//If the email has not been sent we delete the new row of the uuid for avoiding duplicate keys
					DeleteUuidRow(EmailUuid)

					//Also we must delete the user from the DB and from Keycloak for future sign up
					DeleteUserStatement(newEmployee.User.DNI)
					DeleteKeycloakUser(id)
					return ok, response
				} else {
					return ok, response
				}
			}
		}
	}
}

func PatientSignUpVerification(reqBody []byte) (bool, map[string]interface{}){
	var newPatient Patient

	//The data from reqBody is filled in the newUser
	err := json.Unmarshal(reqBody, &newPatient)

	newPatient.User.DNI = strings.ToLower(newPatient.User.DNI)

	if err != nil{
		lib.TerminalLogger.Error("Impossible to retrieve the data from the JSON")
		lib.DocuLogger.Error("Impossible to retrieve the data from the JSON")
		return false, map[string]interface{}{"state": "Problemas con la lectura de los datos"}
	} else {

		//Verify if the format of the email is correct
		if !verifyEmail(newPatient.User.Email){
			return false, map[string]interface{}{"state": "Formato del correo no válido"}
		} else{

			//if it is correct verify if it exists in the DB
			if checkIfExists("patients", "email", "email", newPatient.User.Email){

				//If it exists in the database this is because it could be from a patient
				//registered or not registered but has appointments
				//If It is the one not registered but has appointments we must update the data
				//instead of creating it. For this we check if the dni is a random dni generated by us
				if !verifyIfDniIsRandom(newPatient.User.Email){

					//As the dni is normal and the email exists
					return false, map[string]interface{}{"state": "Este email ya esta en uso"}
				} else{

					//As the dni is one created by us. Verify the data and update it on cascade.
					ok, response := existingPatientVerification("patients", newPatient.User.DNI, newPatient.User.Phone, newPatient.User.Password)
					if !ok {
						return false, response
					} else {

						//If all the data is correct the user is updated in the DB and inserted in Keycloak but not enabled for the login
						return updatePatientAndSendEmail(newPatient)
					}
				}
			} else {

				//The email doesnt exist in the DB so we do the insert and send the verification email after verificating all the data
				ok, response := userDataVerifications("patients", newPatient.User.DNI, newPatient.User.Phone, newPatient.User.Email, newPatient.User.Password)
				if !ok {
					return false, response
				} else {

					//If all the data is correct the user is inserted in the DB and Keycloak but not enabled for the login
					return insertPatientAndSendEmail(newPatient)
				}
			}
		}

	}
}

func insertPatientAndSendEmail(newPatient Patient) (bool, map[string]interface{}) {

	//As all the data is correct we insert the user
	if ok, id := doPatientInsert(newPatient); !ok {
		return false, map[string]interface{}{"state": "Imposible añadir el usuario en la BBDD"}
	} else {

		//if the user has been created we must wait for the email verification so we send the email with a verification URL
		EmailUuid := generateUUID()
		if !insertUuidExpTimeAndUserId(EmailUuid, id, newPatient.User.Email){
			return false, map[string]interface{}{"state": "Imposible de generar el url para la verificación del correo"}
		} else {
			if ok, response := CreateVerificationEmail(EmailUuid, newPatient.User.Name, newPatient.User.Email,"email-verification", 1); !ok{

				//If the email has not been sent we delete the new row of the uuid for avoiding duplicate keys
				DeleteUuidRow(EmailUuid)

				//Also we must delete the user from the DB and from Keycloak
				DeleteUserStatement(newPatient.User.DNI)
				DeleteKeycloakUser(id)
				return ok, response
			} else {
				return ok, response
			}
		}
	}
	return false, nil
}

func updatePatientAndSendEmail(newPatient Patient) (bool, map[string]interface{}) {

	//As all the data is correct we insert the user
	if ok, id := doPatientUpdateAndInsert(newPatient); !ok {
		return false, map[string]interface{}{"state": "Imposible añadir el usuario en la BBDD"}
	} else {

		//if the user has been created we must wait for the email verification so we send the email with a verification URL
		EmailUuid := generateUUID()
		if !insertUuidExpTimeAndUserId(EmailUuid, id, newPatient.User.Email){
			return false, map[string]interface{}{"state": "Imposible de generar el url para la verificación del correo"}
		} else {
			if !createClinicalBackground(newPatient.User.DNI){
				//Delete the user from the DB and from Keycloak
				DeleteUserStatement(newPatient.User.DNI)
				DeleteKeycloakUser(id)
				return false, nil
			} else if ok, response := CreateVerificationEmail(EmailUuid, newPatient.User.Name, newPatient.User.Email,"email-verification", 1); !ok{

				//If the email has not been sent we delete the new row of the uuid for avoiding duplicate keys
				DeleteUuidRow(EmailUuid)

				//Also we must delete the user from the DB and from Keycloak
				DeleteUserStatement(newPatient.User.DNI)
				DeleteKeycloakUser(id)
				return ok, response
			} else {
				return ok, response
			}
		}
	}
	return false, nil
}

func UpdateEnabledUserFromSlug(uuid string) (bool, map[string]interface{}) {
	userId := getUserIdFromUuid(uuid)

	if userId == ""{
		lib.TerminalLogger.Warn("The slug doesnt have any user id")
		lib.DocuLogger.Warn("The slug doesnt have any user id")
		return false, map[string]interface{}{"state": "El slug no coincide con ningun usuario"}
	} else {
		if !updateUserEnabled(userId){
			return false, map[string]interface{}{"state": "No se ha podido actualizar el usuario"}
		}else{

			//Once the user is updated we must delete the SLUG row from email verification for avoiding duplicate keys
			DeleteUuidRow(uuid)
			lib.TerminalLogger.Info("Now the user " + userId + " can log")
			lib.DocuLogger.Info("Now the user" + userId + " can log")
			return true, map[string]interface{}{"state": "Usuario disponible para login"}
		}
	}
}

func ResendVerificationEmail(uuid string) (bool, map[string]interface{}) {

	//For resending the email we must retrieve the email from the table Uniqueurl
	email := getEmailFromUuid(uuid)
	if email == ""{
		lib.TerminalLogger.Warn("The slug doesnt have any user id")
		lib.DocuLogger.Warn("The slug doesnt have any user id")
		return false, map[string]interface{}{"state": "El slug no coincide con ningun usuario"}
	} else {
		_, name := getStringFromField("users", "name", "email", email)
		return CreateVerificationEmail(uuid, name, email," email-verification", 1)
	}
}

func AppointmentMiddleware(reqBody []byte) (bool, map[string]interface{}){
	var appointmentData Appointment

	err := json.Unmarshal(reqBody, &appointmentData)

	if err != nil{
		lib.TerminalLogger.Error("Impossible to retrieve the data from the JSON")
		lib.DocuLogger.Error("Impossible to retrieve the data from the JSON")
		return false, map[string]interface{}{"state": "Problemas con la lectura de los datos"}
	} else if appointmentData.New_user{
		if ok, response := verifyAppointmentData(appointmentData, true); !ok{
			return ok, response
		} else{

			answer, patient_id := getStringFromField("patients", "dni", "email", appointmentData.Patient_email)

			if (!answer && patient_id=="") {
				//generate a random id for the patient as dni
				patient_uuid := "Usuario nuevo: " + generateUUID()

				if !insertNewRandomUser(patient_id, appointmentData.Patient_name, appointmentData.Patient_email, appointmentData.Patient_phone) {
					return false, map[string]interface{}{"state": "Problemas añadiendo el usuario nuevo"}
				} else {
					return verifyAndSaveApointment(appointmentData, patient_uuid)
				}
			} else {
				return verifyAndSaveApointment(appointmentData, patient_id)
			}
		}
	} else{
		if ok, response := verifyAppointmentData(appointmentData, false); !ok{
			return ok, response
		} else{
			_, employee_dni := getStringFromField("employee", "dni", "email", appointmentData.Employee_email)
			date := time.Date(appointmentData.Year, time.Month(appointmentData.Month), appointmentData.Day, appointmentData.Hour, appointmentData.Minute,0, 0, time.UTC)
			dateAsString := date.String()[:20]
			_, patient_dni := getStringFromField("patients", "dni", "email", appointmentData.Patient_email)
			if ok, response := verifyAppointmentAvailableness(patient_dni, employee_dni, dateAsString); !ok{
				return false, response
			} else {
				return saveAppointmentAndSendNotification(patient_dni, employee_dni, dateAsString, appointmentData)
			}
		}
	}
}

func verifyAndSaveApointment(appointmentData Appointment, patient_id string) (bool, map[string]interface{}){
	date := time.Date(appointmentData.Year, time.Month(appointmentData.Month), appointmentData.Day, appointmentData.Hour, appointmentData.Minute, 0, 0, time.UTC)
	dateAsString := date.String()[:20]
	_, employee_dni := getStringFromField("employee", "dni", "email", appointmentData.Employee_email)
	if ok, response := verifyAppointmentAvailableness(patient_id, employee_dni, dateAsString); !ok{
		return false, response
	} else {
		return saveAppointmentAndSendNotification(patient_id, employee_dni, dateAsString, appointmentData)
	}
}

func saveAppointmentAndSendNotification(patient_id, employee_dni, date string, appointmentData Appointment,) (bool, map[string]interface{}) {
	if !insertAppointment(date, employee_dni, patient_id){
		return false, map[string]interface{}{"state": "No se ha podido crear la cita"}
	} else {
		minute := strconv.Itoa(appointmentData.Minute)
		if (appointmentData.Minute >= 0 && appointmentData.Minute<10){
			minute = "0" + minute
		}
		_, employee_name := getStringFromField("employee", "name", "dni", employee_dni)
		ok := sendReminder("Cita fisioterapia", "Tienes una cita pendiente con " + employee_name,
			strconv.Itoa(appointmentData.Day), strconv.Itoa(appointmentData.Hour) + ":" + minute,
			"http://localhost:8081/calendar", appointmentData.Patient_email, appointmentData.Month)
		if ok {
			lib.TerminalLogger.Trace("The reminder has been sent")
			lib.DocuLogger.Trace("The reminder has been sent")
		} else {
			lib.TerminalLogger.Error("The reminder has not been sent")
			lib.DocuLogger.Error("The reminder has not been sent")
		}
		setAppointmentReminder(appointmentData)
		return true, map[string]interface{}{"state": "Cita creada"}
	}
}

func ExerciseMiddleware(reqBody []byte) (bool, map[string]interface{}){
	var exerciseData Exercise

	err := json.Unmarshal(reqBody, &exerciseData)

	if err != nil{
		lib.TerminalLogger.Error("Impossible to retrieve the data from the JSON")
		lib.DocuLogger.Error("Impossible to retrieve the data from the JSON")
		return false, map[string]interface{}{"state": "Problemas con la lectura de los datos"}
	} else if ok, response := verifyExerciseData(exerciseData); !ok{
		return ok, response
	} else {
		_, patient_dni := getStringFromField("patients", "dni", "email", exerciseData.Patient_email)
		date := time.Date(exerciseData.Year, time.Month(exerciseData.Month), exerciseData.Day, exerciseData.Hour, exerciseData.Minute,0, 0, time.UTC)
		dateAsString := date.String()[:20]
		if !verifyPatientAvaliableness(patient_dni, dateAsString, "exercises"){
			return false, map[string]interface{}{"state": "El apciente ya tiene un ejercicio en esta fecha"}
		} else {
			return saveExerciseAndSendNotification(patient_dni, dateAsString, exerciseData)
		}
	}
}

func saveExerciseAndSendNotification(patient_dni string, date string, exerciseData Exercise) (bool, map[string]interface{}) {
	if !insertExercise(date, patient_dni, exerciseData.Exercise_name, exerciseData.Description){
		return false, map[string]interface{}{"state": "No se ha podido crear el ejercicio"}
	} else {
		minute := strconv.Itoa(exerciseData.Minute)
		if (exerciseData.Minute >= 0 && exerciseData.Minute<10){
			minute = "0" + minute
		}
		ok := sendReminder("Ejercicio pendiente: " + exerciseData.Exercise_name, exerciseData.Description,
			strconv.Itoa(exerciseData.Day), strconv.Itoa(exerciseData.Hour) + ":" + minute,
			"http://localhost:8081/calendar", exerciseData.Patient_email, exerciseData.Month)
		if ok {
			lib.TerminalLogger.Trace("The reminder has been sent")
			lib.DocuLogger.Trace("The reminder has been sent")
		} else {
			lib.TerminalLogger.Error("The reminder has not been sent")
			lib.DocuLogger.Error("The reminder has not been sent")
		}
		setExerciseReminder(exerciseData)
		return true, map[string]interface{}{"state": "Ejercicio creado"}
	}
}

func GetAppointmentsDataFromDni(token string) (bool, map[string]interface{}){
	dni := getUserDniFromToken(token)
	if ok, rows := getAppointmentsFromDB(dni); !ok{
		return ok, map[string]interface{}{"state": "Ha habido algún problema encontrando las citas"}
	} else {
		return true, map[string]interface{}{"dataToShow": getAppointmentsDataFromRows(rows)}
	}
}

func GetAllAppointmentsData() (bool, map[string]interface{}){
	if ok, rows := getAllAppointmentsFromDB(); !ok{
		return ok, map[string]interface{}{"state": "Ha habido algún problema encontrando las citas"}
	} else {
		return true, map[string]interface{}{"dataToShow": getAppointmentsDataFromRows(rows)}
	}
}

func GetExercisesDataFromDni(token string) (bool, map[string]interface{}){
	dni := getUserDniFromToken(token)
	if ok, rows := getExercisesFromDB(dni); !ok{
		return ok, map[string]interface{}{"state": "Ha habido algún problema encontrando las citas"}
	} else {
		return true, map[string]interface{}{"dataToShow": getExercisesDataFromRows(rows)}
	}
}

func DeleteAppointmentDataFromDni(token string, reqBody []byte) (bool, map[string]interface{}){

	var date Date

	err := json.Unmarshal(reqBody, &date)

	if err != nil{
		lib.TerminalLogger.Error("Impossible to retrieve the data from the JSON")
		lib.DocuLogger.Error("Impossible to retrieve the data from the JSON")
		return false, map[string]interface{}{"state": "Problemas con la lectura de los datos"}
	} else {
		dni := getUserDniFromToken(token)
		if !deleteAppointmentFromDB(dni, date.Date) {
			return false, map[string]interface{}{"state": "Ha habido algún problema encontrando la cita a borrar"}
		} else {
			return true, map[string]interface{}{"State": "Cita borrada"}
		}
	}
}

func DeleteExerciseDataFromDni(token string, reqBody []byte) (bool, map[string]interface{}){
	var date Date

	err := json.Unmarshal(reqBody, &date)

	if err != nil{
		lib.TerminalLogger.Error("Impossible to retrieve the data from the JSON")
		lib.DocuLogger.Error("Impossible to retrieve the data from the JSON")
		return false, map[string]interface{}{"state": "Problemas con la lectura de los datos"}
	} else {
		dni := getUserDniFromToken(token)
		if !deleteExerciseFromDB(dni, date.Date) {
			return false, map[string]interface{}{"state": "Ha habido algún problema encontrando el ejercicio a borrar"}
		} else {
			return true, map[string]interface{}{"State": "Ejercicio borrado"}
		}
	}
}

func GetClinicalBackgroundMiddleware(reqBody []byte) (bool, map[string]interface{}){
	var data ClinicalBackgroundData

	err := json.Unmarshal(reqBody, &data)

	if err != nil{
		lib.TerminalLogger.Error("Impossible to retrieve the data from the JSON")
		lib.DocuLogger.Error("Impossible to retrieve the data from the JSON")
		return false, map[string]interface{}{"state": "Problemas con la lectura de los datos"}
	} else {
		if ok, clinicalData := getClinicalBackground(data); !ok {
			return false, map[string]interface{}{"state": "Ha habido algún problema actualizando los datos del historila clínico"}
		} else {
			return true, map[string]interface{}{"Data": clinicalData}
		}
	}
}

func UpdateClinicalBackgroundMiddleware(reqBody []byte) (bool, map[string]interface{}){
	var data ClinicalBackgroundData

	err := json.Unmarshal(reqBody, &data)

	if err != nil{
		lib.TerminalLogger.Error("Impossible to retrieve the data from the JSON")
		lib.DocuLogger.Error("Impossible to retrieve the data from the JSON")
		return false, map[string]interface{}{"state": "Problemas con la lectura de los datos"}
	} else {
		if !updateClinicalBackground(data) {
			return false, map[string]interface{}{"state": "Ha habido algún problema actualizando los datos del historila clínico"}
		} else {
			return true, map[string]interface{}{"State": "Historial clínico actualizado"}
		}
	}
}

func GetAllEmployeesData() (bool, map[string]interface{}){
	if ok, rows := getAllEmployeesFromDB(); !ok{
		return ok, map[string]interface{}{"state": "Ha habido algún problema encontrando las citas"}
	} else {
		return true, map[string]interface{}{"dataToShow": getEmployeeDataFromRows(rows)}
	}
}

func GetAllPatientsData() (bool, map[string]interface{}){
	if ok, rows := getAllPatientsFromDB(); !ok{
		return ok, map[string]interface{}{"state": "Ha habido algún problema encontrando las citas"}
	} else {
		return true, map[string]interface{}{"dataToShow": getPatientDataFromRows(rows)}
	}
}

func generateUUID() string {
	uuidWithHyphen := uuid.New()
	customUuid := strings.Replace(uuidWithHyphen.String(), "-", "", -1)
	return customUuid
}


