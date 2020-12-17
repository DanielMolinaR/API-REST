package middleware

import (
	"TFG/API-REST/lib"
	"TFG/API-REST/structures"
	"github.com/badoux/checkmail"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var letters = []string{"T", "R", "W", "A", "G", "M", "Y", "F", "P", "D", "X", "B", "N",	"J", "Z", "S", "Q",	"V", "H", "L", "C", "K", "E"}

func userDataVerifications(condition, dni, phone, email, password string) (bool, map[string]interface{}){


	//verify if the email is correct and if it exists in the DB
	if !verifyEmail(email){
		return false, map[string]interface{}{"state": "Correo no válido"}
	} else if checkIfExists(condition,"email", email){
		if ok, response := newUserUpdateVerifications(email); ok{
			return ok, response
		}
		//if the email exists we must take the dni a check if is a dni or a new user
		//if the dni is from a new user we must update the user and return true
		//if the dni is not from a new user we return false
		//before updating the user we have to verify the data
		return false, map[string]interface{}{"state": "Ya existe este correo"}
	} else{
		//If the email doesnt exists in the DB we continue with the verification
		if ok, response := dataFormVerification(dni, phone, password); !ok{
			return ok, response
		} else{
			return verifyIfSomeDataExist(condition, dni, phone)
		}
	}

}

func dataFormVerification(dni, phone, password string) (bool, map[string]interface{}) {

	//verify if the DNI is correct
	if !verifyDNI(dni){
		return false, map[string]interface{}{"state": "DNI incorrecto"}
	}

	//Verify if the password is strong
	if !verifyPasswordIsSafe(password){
		return false, map[string]interface{}{"state": "Contraseña débil"}
	}

	//verify if the phone is correct and if it exists in the DB
	if !verifyPhoneNumber(phone){
		return false, map[string]interface{}{"state": "teléfono no válido"}
	}
	//If everything is correct return true
	return true, nil
}

func verifyIfSomeDataExist(condition, dni, phone string) (bool, map[string]interface{}) {

	//verify if the dni doesnt exist in the DDBB
	if checkIfExists(condition,"dni", dni){
		return false, map[string]interface{}{"state": "Ya existe este DNI"}
	}

	//verify if the phone doesnt exist in the DDBB
	if checkIfExists(condition,"phone", phone){
		return false, map[string]interface{}{"state": "Ya existe este número de teléfono"}
	}

	//If everything is correct return true
	return true, nil

}

func newUserUpdateVerifications(email, dni string) (bool, map[string]interface{}){
	//if the email exists we must take the dni a check if is a dni or a new user
	//if the dni is from a new user to verify that is him we comparrer the phone numbers we must update the user and return true
	//if the dni is not from a new user we return false
	//before updating the user we have to verify the data
	//PROBLEM I dont have all the user's data to update it
	//SOLUTION dont have data pe data as parameter and have the struct with all the data
	//or do all the verifications one method before
	dniFromDB := getStringFromUsersWithEmail("dni", email)
	if (verifyDNI(dniFromDB)){
		return false, nil
	}

	phoneFromDB := getStringFromUsersWithEmail("phone", email)



}

func verifyDNI(dni string) bool{
	//The DNI must has 9 characters
	if len(dni)!=9{
		return false
	} else {
		//The last char of the DNI must be a Letter
		if !verifyLastCharIsALetter(dni){
			return false
		} else {
			//Verify if the Letter is correct with the numbers of the DNI
			if !verifyLetterIsCorrect(dni){
				return false
			}
		}
	}
	return true
}

func verifyLastCharIsALetter(dni string) bool{
	//Take the last char
	c := strings.ToUpper(dni[8:])
	//Verified if the last char is a Letter
	// parsing it to and int and using ASCII
	asciiValue := int(c[0])
	if asciiValue < 65 || asciiValue > 90 {
		return false
	} else {
		return true
	}
}

func verifyLetterIsCorrect (dni string) bool {
	//Parse to int the DNI except the last char
	dniNumber, err := strconv.Atoi(dni[0:8])
	if err!=nil{
		return false
	}
	//The module of the division of the number of the DNI
	// by 23, must be the position of the Letter in dniLetter[]
	//This rule is established by Spain
	if letters[dniNumber%23] != dni[8:]{
		return false
	}
	return true
}

func verifyPhoneNumber(phone string) bool {

	if len(phone)!=9 {
		return false
	} else if !allCharAreNumbers(phone){
		return false
	}
	return true
}

func allCharAreNumbers(phone string) bool{
	//The range of a string return an int32
	//because It iterates over UTF-8-encoded
	//code points in the string
	for i, ch := range phone{
		if int(ch) < 48 || int(ch) > 57{
			return false
		}
		if i == 0{
			//Verify if the first digit of the number
			//matches with one of the three types of
			//phone numbers in Spain (6,7 or 9)
			if !verifyDigit(int(ch)){
				return false
			}
		}
	}
	return true
}

func verifyDigit(c int) bool{
	if c == 54 || c == 55 || c == 57{
		return true
	}
	return false
}

func verifyEmail (email string) bool {
	//Validate Format
	err := checkmail.ValidateFormat(email)
	if err != nil {
		return false
	}
	//Validate Domain
	err = checkmail.ValidateHost(email)
	if err != nil {
		return false
	}
	//Validate User
	err = checkmail.ValidateHost(email)
	if _, ok := err.(checkmail.SmtpError); ok && err != nil {
		return false
	}
	return true
}

func verifyPasswordIsSafe(s string) bool {
	//Validate if the password has at least
	//one letter in upper case, another one
	//in lower case, a special character,
	//a number and if It's longer than 6
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
		hasntSpace = true
	)
	if len(s) >= 7 {
		hasMinLen = true
	}
	for _, char := range s {
		switch {
		//
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		case int(char) == 32:
			hasntSpace = false
		}
	}
	//If every value is true the password is safe
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial && hasntSpace
}

func VerifyUuidAndGetExpTime(uuid string) int64 {

	//If the uuid doesnt exists the time is 0
	var unixExpTime int64 = 0
	unixExpTime = getExpTimeFromUuid(uuid)
	return unixExpTime
}

func VerifyExpTime(unixExpTime int64) (bool) {
	expTime := time.Unix(unixExpTime, 0)
	if time.Now().After(expTime){
		lib.TerminalLogger.Error("The uuid has expired. Expiration date: ", expTime)
		lib.DocuLogger.Error("The uuid has expired. Expriration Date: ", expTime)
		return false
	} else {
		lib.TerminalLogger.Error("The slug is correct and it is not expired")
		lib.DocuLogger.Error("The slug is correct and it is not expired")
		return true
	}
}

func verifyAppointmentData(appointmentData structures.Appointment) (bool, map[string]interface{}) {
	if !verifyEmail(appointmentData.Employee_email){
		return false, map[string]interface{}{"state": "Correo no válido"}
	} else if checkIfExists("employee","email", appointmentData.Employee_email){
		return false, map[string]interface{}{"state": "Ya existe este correo"}
	}

	if !verifyEmail(appointmentData.Patient_email){
		return false, map[string]interface{}{"state": "Correo no válido"}
	} else if checkIfExists("patients","email", appointmentData.Patient_email){
		return false, map[string]interface{}{"state": "Ya existe este correo"}
	}

	verifyTime(appointmentData.Year, appointmentData.Month, appointmentData.Day, appointmentData.Hour, appointmentData.Minute)

		//verify time

	return true, nil
}

func verifyTime(year, month, day, hour, minute int) bool{

}

