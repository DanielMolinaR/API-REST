package middleware

import (
	"TFG/API-REST/src/lib"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type Keys struct{
	PwdKey string `json:"pwdKey"`
	DataKey string `json:"dataKey"`
	Nonce string `json:"nonce"`
}

func encryptPwd(pwdToEncrypt string) (string, error) {

	//Since the key is in string, we need to convert decode it to bytes,
	//in this case we use the PwdKey for an as
	keyString, nonceString := getThePwdKey()
	key, _ := hex.DecodeString(keyString)
	plaintext := []byte(pwdToEncrypt)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return pwdToEncrypt, err
	}

	//Create a new GCM - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	//https://golang.org/pkg/crypto/cipher/#NewGCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return pwdToEncrypt, err
	}

	//Create a nonce. Nonce should be from GCM
	nonce, _ := hex.DecodeString(nonceString)

	//Encrypt the data using aesGCM.Seal
	//Since we don't want to save the nonce somewhere else in this case, we add it as a prefix
	//to the encrypted data. The first nonce argument in Seal is the prefix.
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return hex.EncodeToString(ciphertext), nil
}

func ComparePwdAndHash(pwd, pwdHashed string) bool{

	//Encrypt the inserted password
	cipherPwd, err := encryptPwd(pwd)
	if err != nil {
		lib.ErrorLogger.Println("could not ecrypt the password: %v", err)
		return false
	}

	//Compare it with the password encrypted in the DB
	if cipherPwd == pwdHashed{
		return true
	}
	return false
}

func encryptData(stringToEncrypt string) (string, error) {

	//Since the key is in string, we need to convert decode it to bytes
	key, _ := hex.DecodeString(getTheDataKey())
	plaintext := []byte(stringToEncrypt)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return stringToEncrypt, err
	}

	//Create a new GCM - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	//https://golang.org/pkg/crypto/cipher/#NewGCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return stringToEncrypt, err
	}

	//Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return stringToEncrypt, err
	}

	//Encrypt the data using aesGCM.Seal
	//Since we don't want to save the nonce somewhere else in this case, we add it as a prefix
	//to the encrypted data. The first nonce argument in Seal is the prefix.
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return string(ciphertext), nil
}

func decryptData(encryptedString string) (string, error) {

	key, _ := hex.DecodeString(getTheDataKey())
	enc, _ := hex.DecodeString(encryptedString)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	//Get the nonce size
	nonceSize := aesGCM.NonceSize()

	//Extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	//Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func getTheDataKey() string {
	dataconfig, err := os.Open("./API-REST/src/middleware/keys.json")
	if err != nil {
		fmt.Println(err)
	}
	jsonBody, _ := ioutil.ReadAll(dataconfig)
	var key Keys
	json.Unmarshal(jsonBody, &key)

	return key.DataKey
}

func getThePwdKey() (string, string) {
	dataconfig, err := os.Open("./API-REST/src/middleware/keys.json")
	if err != nil {
		fmt.Println(err)
	}
	jsonBody, _ := ioutil.ReadAll(dataconfig)
	var key Keys
	json.Unmarshal(jsonBody, &key)

	return key.PwdKey, key.Nonce
}