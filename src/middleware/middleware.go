package middleware

import (

	. "TFG/API-REST/src/structures"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gbrlsnchs/jwt"
	"math/rand"
	"time"
)

type CustomPayload struct {
	jwt.Payload
}

var hs = jwt.NewHS256([]byte("secret"))

func UsersLogin(reqBody []byte) (bool, map[string]interface{}) {
	var userToLog Users
	//The data from reqBody is filled in the newUser
	json.Unmarshal(reqBody, &userToLog)
	//verify that the DNI or the Email exists
	if len(userToLog.DNI) == 0 && len(userToLog.Email) != 0{
		if !checkIfExists(userToLog.Email, "email"){
			return false, map[string]interface{}{"state": "El usuario no existe"}
			//If exists check the password
		} else if bool, response := checkIfPassswordIsCorrect(userToLog.Email, userToLog.Password); !bool{
			return false, map[string]interface{}{"state": response}
		}
	} else if len(userToLog.DNI) != 0 && len(userToLog.Email) == 0 {
		if !checkIfExists(userToLog.DNI, "dni"){
			return false, map[string]interface{}{"state": "El usuario no existe"}
		}
		//If exists check the password
		if bool, response := checkIfPassswordIsCorrect(userToLog.DNI, userToLog.Password); !bool{
			return false, map[string]interface{}{"state": response}
		}
	}
	if len(userToLog.DNI) == 0 && len(userToLog.Email) != 0 {
		return true, map[string]interface{}{"state": "Sesión inicada", "name": getUserName(userToLog.Email, "email"),
			"userId": getUserId(userToLog.Email, "email"), "token": generateToken()}
	} else {
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
		fmt.Println(err)
		return ""
	}
	//verify the token
	var pl2 CustomPayload
	_, err = jwt.Verify(token, hs, &pl2)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return hex.EncodeToString(token)
}

func isLogged(token string) bool{
	tokenDecoded, _ := hex.DecodeString(token)
	//validate
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
	_, err := jwt.Verify(tokenDecoded, hs, &pl, validatePayload)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
