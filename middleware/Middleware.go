package middleware

import (
	"TFG/API-REST/lib"
	. "TFG/API-REST/structures"
	"encoding/json"
	"github.com/google/uuid"
	"strings"
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
			"Refresh token": refreshToken, "Role": GetTheRole(accessToken), "Email": getEmail(accessToken)}
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

		//If all the data is correct the user is inserted in the DB and Keycloak but no enabled for login
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

					//Also we must delete the user from the DB and from Keycloak for new sign up
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
		ok, response := userDataVerifications("patients", newPatient.User.DNI, newPatient.User.Phone, newPatient.User.Email, newPatient.User.Password)
		if !ok {
			return false, response
		}

		//If all the data is correct the user is inserted in the DB and Keycloak but no enabled for login
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
	}
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
		return CreateVerificationEmail(uuid, getStringFromUsersWithEmail("name", email), email," email-verification", 1)
	}
}

func VerifySendAndNotifyAppointment(reqBody []byte) (bool, map[string]interface{}){
	var appointmentData Appointment

	err := json.Unmarshal(reqBody, &appointmentData)

	//newUser true o false

	//verificar los datos de dni que existan en la BBDD y que ninguno de los dos ya tenga una cita ese dia a esa hora

	if err != nil{
		lib.TerminalLogger.Error("Impossible to retrieve the data from the JSON")
		lib.DocuLogger.Error("Impossible to retrieve the data from the JSON")
		return false, map[string]interface{}{"state": "Problemas con la lectura de los datos"}
	} else {
		if appointmentData.New_user{
			addNewUserWithRandomDNI(appointmentData.Patient_name, appointmentData.Patient_email, appointmentData.Patient_phone)
		}
		verifyAppointmentData(appointmentData)
	}
}


func generateUUID() string {
	uuidWithHyphen := uuid.New()
	customUuid := strings.Replace(uuidWithHyphen.String(), "-", "", -1)
	return customUuid
}


