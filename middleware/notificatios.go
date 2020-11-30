package middleware

import (
	"TFG/API-REST/lib"
	. "TFG/API-REST/structures"
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/smtp"
	"os"
)

var secret Users

func init() {
	dataconfig, err := os.Open("./API-REST/middleware/ePassword.json")
	if err != nil {
		fmt.Println(err)
	}
	jsonBody, _ := ioutil.ReadAll(dataconfig)
	json.Unmarshal(jsonBody, &secret)

}

//Request struct
type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

// smtpServer data to smtp server
type smtpServer struct {
	host string
	port string
}

func NewRequest(to []string, subject, body string) *Request {
	return &Request{
		to:      to,
		subject: subject,
		body:    body,
	}
}

// Address URI to smtp server
func (s *smtpServer) Address() string {
	return s.host + ":" + s.port
}

func CreateEmail(uuid, name, email, path string) (bool, map[string]interface{}) {

	url := "http://localhost:8081/" + path + "/" + uuid
	if !sendEmail(name, email, url) {
		lib.TerminalLogger.Error("Impossible to send the email")
		lib.DocuLogger.Error("Impossible to send the email")
		return false, map[string]interface{}{"state": "Imposible enviar el correo"}
	} else {
		lib.TerminalLogger.Trace("Email for Sign Up sent to: " + email)
		lib.DocuLogger.Trace("Email for Sign Up sent to: " + email)
		return true, map[string]interface{}{"state": "Correo enviado"}
	}
}

func sendEmail(name, email, url string) bool {
	// Sender data.
	from := secret.Email
	password := secret.Password
	// smtp server configuration.
	smtpServer := smtpServer{host: "smtp.gmail.com", port: "587"}
	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpServer.host)

	templateData := struct {
		Name string
		URL string
	}{
		Name: name,
		URL: url,
	}

	emailToSend := []string{email}

	r := NewRequest(emailToSend, "BIENVENIDO", "")
	err := r.ParseTemplate("C:/Users/Daniel/go/src/TFG/API-REST/middleware/template.html", templateData)
	if err != nil {
		lib.TerminalLogger.Error("Problems parsing the tempalte: " + err.Error())
		lib.DocuLogger.Error("Problems parsing the tempalte: " + err.Error())
		return false
	}

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	to := "To: "
	for i := 0; i < len(emailToSend); i++{
		if i == 0{
			to += emailToSend[i]
		} else{
			to += ", " + emailToSend[i]
		}
	}
	to += " \n"
	subject := "Subject: " + r.subject + "!\n"
	msg := []byte(to + subject + mime + "\n" + r.body)

	// Sending email.
	err = smtp.SendMail(smtpServer.Address(), auth, "pruebaapifisio@gmail.com", r.to, msg)
	if err != nil {
		lib.TerminalLogger.Error("Problems sending the email: " + err.Error())
		lib.DocuLogger.Error("Problems sending the email: " + err.Error())
		return false
	}
	return true
}

func (r *Request) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}