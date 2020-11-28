package middleware

import (
	"TFG/API-REST/lib"
	. "TFG/API-REST/structures"
	"encoding/json"
	"fmt"
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
			"Refresh token": refreshToken, "Role": GetTheRole(accessToken)}
	}
}

func GetTheRole(token string) int {
	roles := DecodeToken(token)
	allTheRoles := (*roles)["realm_access"].(map[string]interface{})["roles"]
	data := allTheRoles.([]interface{})
	answer := 0
	for i := 0; i < len(data); i++{
		userRole := fmt.Sprintf("%v", data[i])
		if userRole == "EMPLOYEE_ROLE"{
			answer = 1
		}
		if userRole == "ADMIN_ROLE"{
			answer = 2
		}
	}
	return answer
}

func VerifyAndSendEmail(reqBody []byte) (bool, map[string]interface{}){
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
		uuid := generateUUID()
		if !insertUuidExpTimeAndUserId(uuid, "", ""){
			return false, map[string]interface{}{"state": "Imposible de generar el url unico"}
		} else {
			if ok, response := CreateEmail(uuid, userData.Name, userData.Email,"employee-sign-up"); !ok{

				//If the email has not been sent we delete the new row of the uuid for avoiding duplicate keys
				DeleteUuidRow(uuid)
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
		ok, response := signUpVerifications("employee", newEmployee.User.DNI, newEmployee.User.Phone, newEmployee.User.Email, newEmployee.User.Password)
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
				if ok, response := CreateEmail(EmailUuid, newEmployee.User.Name, newEmployee.User.Email,"email-verification"); !ok{

					//If the email has not been sent we delete the new row of the uuid for avoiding duplicate keys
					DeleteUuidRow(EmailUuid)

					//Also we must delete the user from the DB and from Keycloak
					//DeleteUser(dni)
					//DeleteUserFromKeycloak(id)
					return ok, response
				} else {
					return ok, response
				}
			}
		}
	}
}

func PatientSignInVerification(reqBody []byte) (bool, map[string]interface{}){
	var newPatient Patient

	//The data from reqBody is filled in the newUser
	err := json.Unmarshal(reqBody, &newPatient)

	if err != nil{
		lib.TerminalLogger.Error("Impossible to retrieve the data from the JSON")
		lib.DocuLogger.Error("Impossible to retrieve the data from the JSON")
		return false, map[string]interface{}{"state": "Problemas con la lectura de los datos"}
	} else {
		ok, response := signUpVerifications("patients", newPatient.User.DNI, newPatient.User.Phone, newPatient.User.Email, newPatient.User.Password)
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
				return false, map[string]interface{}{"state": "Imposible de generar el url unico"}
			} else {
				if ok, response := CreateEmail(EmailUuid, newPatient.User.Name, newPatient.User.Email,"email-verification"); !ok{

					//If the email has not been sent we delete the new row of the uuid for avoiding duplicate keys
					DeleteUuidRow(EmailUuid)

					//Also we must delete the user from the DB and from Keycloak
					//DeleteUser(dni)
					//DeleteUserFromKeycloak(id)
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
		return CreateEmail(uuid, getNameFromUsersWithEmail(email), email," email-verification")
	}
}




func generateUUID() string {
	uuidWithHyphen := uuid.New()
	uuid := strings.Replace(uuidWithHyphen.String(), "-", "", -1)
	return uuid
}


