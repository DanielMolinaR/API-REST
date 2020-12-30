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
var months = []string{"Enero", "Febrero", "Marzo", "Abril", "Mayo", "Junio", "Julio", "Agosto", "Septiembre", "Octubre", "Noviembre", "Diciembre"}

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

func CreateVerificationEmail(uuid, name, email, path string, condition int) (bool, map[string]interface{}) {

	url := "http://localhost:8081/" + path + "/" + uuid
	if !sendEmail(name, url, email, condition) {
		lib.TerminalLogger.Error("Impossible to send the email")
		lib.DocuLogger.Error("Impossible to send the email")
		return false, map[string]interface{}{"state": "Imposible enviar el correo"}
	} else {
		lib.TerminalLogger.Trace("Email for Sign Up sent to: " + email)
		lib.DocuLogger.Trace("Email for Sign Up sent to: " + email)
		return true, map[string]interface{}{"state": "Correo enviado"}
	}
}

func sendEmail(name, url, email string, condition int) bool {
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

	switch condition{
	case 0: //Case where the email is sent for the SignUp of a new employee
		err := r.ParseTemplate("C:/Users/Daniel/go/src/TFG/API-REST/middleware/Employee_signUp_template.html", templateData)
		if err != nil {
			lib.TerminalLogger.Error("Problems parsing the tempalte: " + err.Error())
			lib.DocuLogger.Error("Problems parsing the tempalte: " + err.Error())
			return false
		}
		break;
	case 1: //Case where the email is sent for verify the email address of the user
		err := r.ParseTemplate("C:/Users/Daniel/go/src/TFG/API-REST/middleware/Email_verification_template.html", templateData)
		if err != nil {
			lib.TerminalLogger.Error("Problems parsing the tempalte: " + err.Error())
			lib.DocuLogger.Error("Problems parsing the tempalte: " + err.Error())
			return false
		}
		break;
	default:
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
	err := smtp.SendMail(smtpServer.Address(), auth, "pruebaapifisio@gmail.com", r.to, msg)
	if err != nil {
		lib.TerminalLogger.Error("Problems sending the email: " + err.Error())
		lib.DocuLogger.Error("Problems sending the email: " + err.Error())
		return false
	}
	return true
}

func sendReminder(header, data, day, hour, url, email string, month int) bool {
	// Sender data.
	from := secret.Email
	password := secret.Password
	// smtp server configuration.
	smtpServer := smtpServer{host: "smtp.gmail.com", port: "587"}
	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpServer.host)

	templateData := struct {
		Header	string
		Data	string
		Day		string
		Hour	string
		URL 	string
	}{
		Header: header,
		Data: data,
		Day: day + " de " + months[month-1],
		Hour: hour,
		URL: url,
	}

	emailToSend := []string{email}

	r := NewRequest(emailToSend, "BIENVENIDO", "")

	err := r.ParseTemplate("C:/Users/Daniel/go/src/TFG/API-REST/middleware/Reminder_template.html", templateData)
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

func sendReminderEmailToEmployees() {
	//coger las citas del dia de hoy
	//clasificar por dni de empleados
	//coger los email de los empleados y referenciar las citas a ese email
	//enviar a cada empleado un email con sus citas

}