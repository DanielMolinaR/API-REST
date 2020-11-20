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
	}else if ok, accessToken, refreshToken := UserCredentialsLogin(userToLogIn.DNI, userToLogIn.Password); !ok{ //Try to login in keycloak
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

func GenerateAndSendURL(reqBody []byte) (bool, map[string]interface{}){
	var userData Users

	//The data from reqBody is filled in the newUser
	err := json.Unmarshal(reqBody, &userData)

	if err != nil{
		lib.TerminalLogger.Error("Impossible to retrieve the data from the JSON")
		lib.DocuLogger.Error("Impossible to retrieve the data from the JSON")
		return false, map[string]interface{}{"state": "Problemas con la lectura de los datos"}
	} else if !verifyEmail(userData.Email){ //verify if the format of the email is correct
		lib.TerminalLogger.Error("Email format incorect")
		lib.DocuLogger.Error("Email format incorect")
		return false, map[string]interface{}{"state": "Formato del correo electrónico incorrecto"}
	} else {
		uuid := generarteUUID()
		if !insertUuidAndExpTime(uuid){
			return false, map[string]interface{}{"state": "Imposible de generar el url unico"}
		} else {
			url := "http://localhost:8081/employee-sign-up/" + uuid
			if !SendEmailForSignUp(userData.Name, userData.Email, url) {
				deleteUuidRow(uuid)
				lib.TerminalLogger.Error("Impossible to send the email")
				lib.DocuLogger.Error("Impossible to send the email")
				return false, map[string]interface{}{"state": "Imposible enviar el correo"}
			} else {
				lib.TerminalLogger.Trace("Email for Sign Up sent to: " + userData.Email)
				lib.DocuLogger.Trace("Email for Sign Up sent to: " + userData.Email)
				return true, map[string]interface{}{"state": "Correo enviado"}
			}
		}
	}
}

func generarteUUID() string {
	uuidWithHyphen := uuid.New()
	uuid := strings.Replace(uuidWithHyphen.String(), "-", "", -1)
	return uuid
}

func EmployeeSignInVerification(reqBody []byte) (bool, map[string]interface{}){
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
		bool, response := signInVerifications("employee", newEmployee.User.DNI, newEmployee.User.Phone, newEmployee.User.Email, newEmployee.User.Password)
		if !bool {
			return false, response
			//insertar el usuario en keycloak
		} else if !doEmployeeInsert(newEmployee) {
			return false, map[string]interface{}{"state": "Imposible añadir el usuario en la BBDD"}
		}
		return true, map[string]interface{}{"state": "Usuario Creado"}
	}
}

func PatientSignInVerification(reqBody []byte) (bool, map[string]interface{}){
	var newPatient Patient
	//The data from reqBody is filled in the newUser
	json.Unmarshal(reqBody, &newPatient)
	bool, response := signInVerifications("patients", newPatient.User.DNI, newPatient.User.Phone, newPatient.User.Email, newPatient.User.Password)
	if  !bool{
		return false, response
	}
	if !doPatientInsert(newPatient){
		return false, map[string]interface{}{"state": "Imposible de añadir en la BBDD"}
	}
	return true, map[string]interface{}{"state": "Usuario Creado"}
}

