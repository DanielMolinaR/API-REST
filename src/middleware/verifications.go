package middleware

import (
	"fmt"
	"strconv"
	"strings"
	"github.com/badoux/checkmail"
)

var letters = []string{"T", "R", "W", "A", "G", "M", "Y", "F", "P", "D", "X", "B", "N",	"J", "Z", "S", "Q",	"V", "H", "L", "C", "K", "E"}
var emailAddress = []string{"@gmail.com", "@icloud.es", "@hotmail.com"}

func verifyDNI(dni string) bool{
	//The DNI must has 9 char
	if len(dni)!=9{
		return false
	} else {
		//The last char of the DNI must be a Letter
		if !verifyLastCharIsALetter(dni){
			return false
		} else {
			//Verify if the Letter is correct
			//with the numbers of the DNI
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
	//It checks if the phone number has 9 characters
	if len(phone)!=9 {
		return false
		//It checks if all the characters are numbers
	} else if !allAreNumbers(phone){
		return false
	}
	return true
}

func allAreNumbers(phone string) bool{
	//The range of a string return an int32
	//because It iterates over UTF-8-encoded
	//code points in the string
	for _, ch := range phone{
		if int(ch) < 48 || int(ch) > 57{
			return false
		}
	}
	return true
}

func verifyEmail (email string) bool {
	//It search for the @ char
	err := checkmail.ValidateFormat(email)
	if err != nil {
		return false
	}
	err = checkmail.ValidateHost(email)
	if err != nil {
		return false
	}
	err = checkmail.ValidateHost(email)
	if _, ok := err.(checkmail.SmtpError); ok && err != nil {
		return false
	}
	return true
}

func verifyPasswordIsSafe(password string) bool {
	if len(password) < 6 {
		return false
	}
	return true
}

//mirar que la contraseña no tenga espacio
