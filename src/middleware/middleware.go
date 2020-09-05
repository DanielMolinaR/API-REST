package middleware

import (
	"TFG/API-REST/src/lib"
	. "TFG/API-REST/src/structures"
	"encoding/json"
	"fmt"
	"github.com/gbrlsnchs/jwt"
	"math/rand"
	"strings"
	"time"
)

type CustomPayload struct {
	jwt.Payload
}

var hs = jwt.NewHS256([]byte("secret"))

func UsersLogin(reqBody []byte, token string) (bool, map[string]interface{}) {

	//Check if the token is valid
	if validateToken(token){
		lib.InfoLogger.Println("A user has logged in thanks to his token")
		return true, map[string]interface{}{"state": "Sesión iniciada"}
	}

	var userToLog Users

	//The data from reqBody is filled in the newUser
	json.Unmarshal(reqBody, &userToLog)

	//Verify that the DNI or the Email
	if len(userToLog.DNI) == 0 && len(userToLog.Email) != 0{
		if !checkIfExists(userToLog.Email, "email"){
			lib.WarningLogger.Println("The email doesn't exist")
			return false, map[string]interface{}{"state": "El usuario no existe"}

		//If exists check if the password is correct
		} else if bool, response := checkIfPassswordIsCorrect(userToLog.Email, userToLog.Password); !bool{
			lib.WarningLogger.Println("The password is incorrect")
			return false, map[string]interface{}{"state": response}
		}
	} else if len(userToLog.DNI) != 0 && len(userToLog.Email) == 0 {
		if !checkIfExists(userToLog.DNI, "dni"){
			lib.WarningLogger.Println("The DNI doesn't exist")
			return false, map[string]interface{}{"state": "El usuario no existe"}
		}

		//If exists check if the password is correct
		if bool, response := checkIfPassswordIsCorrect(userToLog.DNI, userToLog.Password); !bool{
			lib.WarningLogger.Println("The password is incorrect")
			return false, map[string]interface{}{"state": response}
		}
	}

	//Return true with a msg of correct login,
	//the name of the user and the position
	if len(userToLog.DNI) == 0 && len(userToLog.Email) != 0 {
		lib.InfoLogger.Println("User logged with the email: %v", userToLog.Email)
		return true, map[string]interface{}{"state": "Sesión inicada", "name": getUserName(userToLog.Email, "email"),
			"userId": getUserId(userToLog.Email, "email"), "token": generateToken()}
	} else {
		lib.InfoLogger.Println("User logged with the DNI: %v", userToLog.DNI)
		return true, map[string]interface{}{"state": "Sesión inicada", "name": getUserName(userToLog.DNI, "dni"),
			"userId": getUserId(userToLog.DNI, "dni"), "token": generateToken()}
	}
}

func EmployeeSignInVerification(reqBody []byte) (bool, string){
	var newEmployee Employee
	//The data from reqBody is filled in the newUser
	json.Unmarshal(reqBody, &newEmployee)
	bool, response := signInVerifications(newEmployee.User.DNI, newEmployee.User.Phone, newEmployee.User.Email, newEmployee.User.Password)
	if  !bool{
		return false, response
	}

	if !DoEmployeeInsert(newEmployee){
		return false, ""
	}
	return true, "Usuario creado"
}

func PatientSignInVerification(reqBody []byte) (bool, string){
	var newPatient Patient
	//The data from reqBody is filled in the newUser
	json.Unmarshal(reqBody, &newPatient)
	bool, response := signInVerifications(newPatient.User.DNI, newPatient.User.Phone, newPatient.User.Email, newPatient.User.Password)
	if  !bool{
		return false, response
	}
	if !DoPatientInsert(newPatient){
		return false, ""
	}
	return true, "Usuario creado"
}

func signInVerifications(dni, phone, email, password string) (bool, string){
	//verifyDNI verify if the DNI is correct
	// and if it exists in the DB
	if !verifyDNI(dni){
		return false, "DNI incorrecto"
	} else if checkIfExists(dni, "dni"){
		return false, "Este DNI ya existe"
	}
	//Phone number verification
	if !verifyPhoneNumber(phone){
		return false, "El numero de telefono no existe"
	}
	//Email verification
	if !verifyEmail(email){
		return false, "Email no váido"
	}
	//Verify if the password is strong
	if !verifyPasswordIsSafe(string(password)){
		return false, "La contraseña es muy débil"
	}
	return true, ""
}

func generateToken() string {
	now := time.Now()
	rand.Seed(time.Now().UnixNano())
	pl := CustomPayload{
		Payload: jwt.Payload{
			Audience:       jwt.Audience{"https://golang.org", "https://jwt.io"},
			ExpirationTime: jwt.NumericDate(now.Add(24 * 7 * time.Hour)),
			NotBefore:      jwt.NumericDate(now.Add(1 * time.Second)), //30 min es mucho reducirlo a menos
			IssuedAt:       jwt.NumericDate(now),
			JWTID:			string(rand.Intn(100)),
		},
	}

	//sign the token
	token, err := jwt.Sign(pl, hs)
	if err != nil {
		lib.ErrorLogger.Println("Error in generating the token: %v", err)
		return ""
	}

	//verify the token
	var pl2 CustomPayload
	_, err = jwt.Verify(token, hs, &pl2)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	//Enconde the []byte token to a string for the json
	return string(token)
}

func validateToken(token string) bool{

	//Extract the Bearer expec from the data of the header
	tokenData := strings.Replace(token, "Bearer ", "", -1)

	//Decode the string token to a []byte
	tokenInBytes := []byte(tokenData)

	//validate the token
	var (
		now = time.Now()
		aud = jwt.Audience{"https://golang.org", "https://jwt.io"}

		// Validate claims "iat", "exp" and "aud".
		iatValidator = jwt.IssuedAtValidator(now)
		expValidator = jwt.ExpirationTimeValidator(now)
		audValidator = jwt.AudienceValidator(aud)
		nbValidator  = jwt.NotBeforeValidator(now)

		// Use jwt.ValidatePayload to build a jwt.VerifyOption.
		// Validators are run in the order informed.
		pl              CustomPayload
		validatePayload = jwt.ValidatePayload(&pl.Payload, iatValidator, expValidator, audValidator, nbValidator)
	)
	_, err := jwt.Verify(tokenInBytes, hs, &pl, validatePayload)
	if err != nil {
		lib.ErrorLogger.Println("Token not valid")
		return false
	}
	return true
}
