package middleware

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"encoding/json"
	"io/ioutil"
	"os"
)

type Keys struct{
	pwd []byte `json:"PwdKey"`
	data []byte `json:"DataKey"`
}

func encrypt(plaintext []byte) ([]byte, error) {
	dataconfig, err := os.Open("./API-REST/src/middleware/keys.json")
	if err != nil {
		fmt.Println(err)
	}
	jsonBody, _ := ioutil.ReadAll(dataconfig)
	var key Keys
	json.Unmarshal(jsonBody, &key)

	c, err := aes.NewCipher(key.pwd)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func ComparePwdAndHash(pwd, pwdHashed []byte) bool{
	cipherPwd, err := encrypt(pwd)
	if err != nil {
		return false
	}
	if bytes.Compare(cipherPwd, pwdHashed) == 0{
		return true
	}
	return false
}