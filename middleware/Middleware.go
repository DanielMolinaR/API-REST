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
		return true, map[string]interface{}{"state": "Sesión iniciada", "Access token": accessToken,
			"Refresh token": refreshToken, "Role": getTheRole(accessToken), "Email": getEmail(accessToken)}
	}
}

func VerifyDataAndSendUniqueEmail(reqBody []byte) (bool, map[string]interface{}){
	var userData Users

	//The data from reqBody is filled in the newUser
	err := json.Unmarshal(reqBody, &userData)

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
					ok, response := existingPatientVerification("patients", newPatient.User.DNI, newPatient.User.Phone, newPatient.User.Email, newPatient.User.Password)
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

			//generate a random id for the patient as dni
			patient_id := "Usuario nuevo: " + generateUUID()
			if !insertNewRandomUser(patient_id, appointmentData.Patient_name, appointmentData.Patient_email, appointmentData.Patient_phone){
				return false, map[string]interface{}{"state": "Problemas añadiendo el usuario nuevo"}
			} else {
				return saveAppointmentAndSendNotification(patient_id, appointmentData)
			}
		}
	} else{
		if ok, response := verifyAppointmentData(appointmentData, false); !ok{
			return ok, response
		} else{
			_, patient_dni := getStringFromField("patients", "dni", "email", appointmentData.Patient_email)
			return saveAppointmentAndSendNotification(patient_dni, appointmentData)
		}
	}
}

func saveAppointmentAndSendNotification(patient_id string, appointmentData Appointment) (bool, map[string]interface{}) {
	date := time.Date(appointmentData.Year, time.Month(appointmentData.Month), appointmentData.Day, appointmentData.Hour, appointmentData.Minute,0, 0, time.UTC)
	_, employee_dni := getStringFromField("employee", "dni", "email", appointmentData.Employee_email)
	if !insertAppointment(date.String()[:20], employee_dni, patient_id){
		return false, map[string]interface{}{"state": "No se ha podido crear la cita"}
	} else {
		minute := strconv.Itoa(appointmentData.Minute)
		if (minute == "0"){
			minute = minute + "0"
		}
		_, employee_name := getStringFromField("employee", "name", "dni", employee_dni)
		sendReminder("Cita fisioterapia", "Tienes una cita pendiente con " + employee_name,
			strconv.Itoa(appointmentData.Day), strconv.Itoa(appointmentData.Hour) + ":" + minute,
			"http://localhost:8081/calendar", appointmentData.Patient_email, appointmentData.Month)
		setReminder(appointmentData, true)
		return true, map[string]interface{}{"state": "Cita creada"}
	}

}

func generateUUID() string {
	uuidWithHyphen := uuid.New()
	customUuid := strings.Replace(uuidWithHyphen.String(), "-", "", -1)
	return customUuid
}


