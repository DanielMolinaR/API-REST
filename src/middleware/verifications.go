package middleware

import (
	"strconv"
	"strings"
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
	for i, ch := range email{
		if string(ch) == "@"{
			//Now It checks if the address exists
			if emailAddressIsCorrect(email[i:]){
				return true
			}
		}
	}
	return false
}

func emailAddressIsCorrect(address string) bool {
	for _, s := range emailAddress{
		if address == s{
			return true
		}
	}
	return false
}

func verifyPasswordIsSafe(password string) bool {
	if len(password) < 6 {
		return false
	}
	return true
}