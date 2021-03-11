package middleware

import (
	"TFG/API-REST/lib"
	. "TFG/API-REST/structures"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4"
	"html/template"
	"io/ioutil"
	"net/smtp"
	"os"
	"strconv"
	"time"
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
		return false, map[string]interface{}{"state": "No se ha podido enviar el correo de verificación"}
	} else {
		lib.TerminalLogger.Trace("Email for Sign Up sent to: " + email)
		lib.DocuLogger.Trace("Email for Sign Up sent to: " + email)
		return true, map[string]interface{}{"state": "Correo de verificación enviado"}
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

func sendNotification(header, data, day, hour, url, email string, month int) bool {
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

	err := r.ParseTemplate("C:/Users/Daniel/go/src/TFG/API-REST/middleware/Notification_template.html", templateData)
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

func sendDailyReminder(rows pgx.Rows, email string){
	var date int64
	var name string
	middlePart := ""
	firstPart := "<!DOCTYPE html><html lang=\"en\" xmlns=\"http://www.w3.org/1999/xhtml\" xmlns:v=\"urn:schemas-microsoft-com:vml\" xmlns:o=\"urn:schemas-microsoft-com:office:office\"><head>\n    <title> Recordatorio MFS </title>\n    <!--[if !mso]><!-- -->\n    <meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge\" />\n    <!--<![endif]-->\n    <meta http-equiv=\"Content-Type\" content=\"text/html; charset=UTF-8\" />\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n    <style type=\"text/css\">\n      #outlook a {\n        padding: 0;\n      }\n  \n      body {\n        margin: 0;\n        padding: 0;\n        -webkit-text-size-adjust: 100%;\n        -ms-text-size-adjust: 100%;\n      }\n  \n      table,\n      td {\n        border-collapse: collapse;\n        mso-table-lspace: 0pt;\n        mso-table-rspace: 0pt;\n      }\n  \n      img {\n        border: 0;\n        height: auto;\n        line-height: 100%;\n        outline: none;\n        text-decoration: none;\n        -ms-interpolation-mode: bicubic;\n      }\n  \n      p {\n        display: block;\n        margin: 13px 0;\n      }\n    </style>\n    <!--[if mso]>\n          <xml>\n          <o:OfficeDocumentSettings>\n            <o:AllowPNG/>\n            <o:PixelsPerInch>96</o:PixelsPerInch>\n          </o:OfficeDocumentSettings>\n          </xml>\n          <![endif]-->\n    <!--[if lte mso 11]>\n          <style type=\"text/css\">\n            .mj-outlook-group-fix { width:100% !important; }\n          </style>\n          <![endif]-->\n    <!--[if !mso]><!-->\n    <link href=\"https://fonts.googleapis.com/css?family=Roboto:100,300,400,700\" rel=\"stylesheet\" type=\"text/css\" />\n    <style type=\"text/css\">\n      @import url(https://fonts.googleapis.com/css?family=Roboto:100,300,400,700);\n    </style>\n    <!--<![endif]-->\n    <style type=\"text/css\">\n      @media only screen and (min-width:480px) {\n        .mj-column-per-100 {\n          width: 100% !important;\n          max-width: 100%;\n        }\n      }\n    </style>\n    <style type=\"text/css\">\n      @media only screen and (max-width:480px) {\n        table.mj-full-width-mobile {\n          width: 100% !important;\n        }\n  \n        td.mj-full-width-mobile {\n          width: auto !important;\n        }\n      }\n    </style>\n    <style type=\"text/css\">\n      a,\n      span,\n      td,\n      th {\n        -webkit-font-smoothing: antialiased !important;\n        -moz-osx-font-smoothing: grayscale !important;\n      }\n    </style>\n  </head>\n  \n  <body style=\"background-color:#f3f3f5;\">\n    <div style=\"display:none;font-size:1px;color:#ffffff;line-height:1px;max-height:0px;max-width:0px;opacity:0;overflow:hidden;\"> MFS </div>\n    <div style=\"background-color:#f3f3f5;\">\n      <!--[if mso | IE]>\n        <table\n           align=\"center\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" class=\"\" style=\"width:600px;\" width=\"600\"\n        >\n          <tr>\n            <td style=\"line-height:0px;font-size:0px;mso-line-height-rule:exactly;\">\n        <![endif]-->\n      <div style=\"margin:0px auto;max-width:600px;\">\n        <table align=\"center\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" role=\"presentation\" style=\"width:100%;\">\n          <tbody>\n            <tr>\n              <td style=\"direction:ltr;font-size:0px;padding:20px 0;text-align:center;\">\n                <!--[if mso | IE]>\n                    <table role=\"presentation\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\">\n                  \n          <tr>\n        \n              <td\n                 class=\"\" style=\"vertical-align:top;width:600px;\"\n              >\n            <![endif]-->\n                <div class=\"mj-column-per-100 mj-outlook-group-fix\" style=\"font-size:0px;text-align:left;direction:ltr;display:inline-block;vertical-align:top;width:100%;\">\n                  <table border=\"0\" cellpadding=\"0\" cellspacing=\"0\" role=\"presentation\" style=\"vertical-align:top;\" width=\"100%\">\n                    <tbody><tr>\n                      <td style=\"font-size:0px;word-break:break-word;\">\n                        <!--[if mso | IE]>\n      \n          <table role=\"presentation\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\"><tr><td height=\"20\" style=\"vertical-align:top;height:20px;\">\n        \n      <![endif]-->\n                        <div style=\"height:20px;\">   </div>\n                        <!--[if mso | IE]>\n      \n          </td></tr></table>\n        \n      <![endif]-->\n                      </td>\n                    </tr>\n                  </tbody></table>\n                </div>\n                <!--[if mso | IE]>\n              </td>\n            \n          </tr>\n        \n                    </table>\n                  <![endif]-->\n              </td>\n            </tr>\n          </tbody>\n        </table>\n      </div>\n      <!--[if mso | IE]>\n            </td>\n          </tr>\n        </table>\n        \n        <table\n           align=\"center\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" class=\"\" style=\"width:600px;\" width=\"600\"\n        >\n          <tr>\n            <td style=\"line-height:0px;font-size:0px;mso-line-height-rule:exactly;\">\n        <![endif]-->\n      <div style=\"background:#ffffff;background-color:#ffffff;margin:0px auto;border-radius:4px 4px 0 0;max-width:600px;\">\n        <table align=\"center\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" role=\"presentation\" style=\"background:#ffffff;background-color:#ffffff;width:100%;border-radius:4px 4px 0 0;\">\n          <tbody>\n            <tr>\n              <td style=\"direction:ltr;font-size:0px;padding:20px 0;text-align:center;\">\n                <!--[if mso | IE]>\n                    <table role=\"presentation\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\">\n                  \n              <tr>\n                <td\n                   class=\"\" width=\"600px\"\n                >\n            \n        <table\n           align=\"center\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" class=\"\" style=\"width:600px;\" width=\"600\"\n        >\n          <tr>\n            <td style=\"line-height:0px;font-size:0px;mso-line-height-rule:exactly;\">\n        <![endif]-->\n                <div style=\"margin:0px auto;max-width:600px;\">\n                  <table align=\"center\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" role=\"presentation\" style=\"width:100%;\">\n                    <tbody>\n                      <tr>\n                        <td style=\"direction:ltr;font-size:0px;padding:0px;text-align:center;\">\n                          <!--[if mso | IE]>\n                    <table role=\"presentation\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\">\n                  \n          <tr>\n        \n              <td\n                 class=\"\" style=\"vertical-align:top;width:600px;\"\n              >\n            <![endif]-->\n                          <div class=\"mj-column-per-100 mj-outlook-group-fix\" style=\"font-size:0px;text-align:left;direction:ltr;display:inline-block;vertical-align:top;width:100%;\">\n                            <table border=\"0\" cellpadding=\"0\" cellspacing=\"0\" role=\"presentation\" style=\"vertical-align:top;\" width=\"100%\">\n                              <tbody><tr>\n                                <td align=\"center\" style=\"font-size:0px;word-break:break-word;\">\n                                  <table border=\"0\" cellpadding=\"0\" cellspacing=\"0\" role=\"presentation\" style=\"border-collapse:collapse;border-spacing:0px;\">\n                                    <tbody>\n                                      <tr>\n                                        <td style=\"width:300px; height: 300;\">\n                                          <img height=\"auto\" src=\"C:\\Users\\Daniel\\Documents\\GitHub\\front_ApiRest\\front-fisioterapia\\src\\assets\\images\\logo-letras.png\" style=\"border:0;display:block;outline:none;text-decoration:none;height:auto;width:100%;font-size:13px;\" width=\"150\" />\n                                        </td>\n                                      </tr>\n                                    </tbody>\n                                  </table>\n                                </td>\n                              </tr>\n                              <tr>\n                                <td style=\"font-size:0px;word-break:break-word;\">\n                                  <!--[if mso | IE]>\n      \n          <table role=\"presentation\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\"><tr><td height=\"20\" style=\"vertical-align:top;height:20px;\">\n        \n      <![endif]-->\n                                  <div style=\"height:20px;\">   </div>\n                                  <!--[if mso | IE]>\n      \n          </td></tr></table>\n        \n      <![endif]-->\n                                </td>\n                              </tr>\n                              <tr>\n                                <td align=\"center\" style=\"font-size:0px;word-break:break-word;\">\n                                  <div style=\"font-family:Roboto, Helvetica, Arial, sans-serif;font-size:24px;font-weight:400;line-height:30px;text-align:center;color:#54595f;\">\n                                    <h1 style=\"font-size: 24px; line-height: normal; font-weight: 400;\">¡Bienvenido de nuevo! </h1>\n                                  </div>\n                                </td>\n                              </tr>\n                              <tr>\n                                <td align=\"left\" style=\"font-size:0px;padding:10px 25px;word-break:break-word;\">\n                                  <div style=\"font-family:Roboto, Helvetica, Arial, sans-serif;font-size:14px;font-weight:400;line-height:20px;text-align:left;color:#54595f;\">\n                                    <p style=\"margin-bottom: 0;\">Por si se te había olvidado: </p>\n                                  </div>\n                                </td>\n                              </tr>\n                              <tr>\n                            </tbody></table>\n                          </div>\n                          <!--[if mso | IE]>\n              </td>\n            \n          </tr>\n        \n                    </table>\n                  <![endif]-->\n                        </td>\n                      </tr>\n                    </tbody>\n                  </table>\n                </div>\n                <!--[if mso | IE]>\n            </td>\n          </tr>\n        </table>\n        \n                </td>\n              </tr>\n            \n                    </table>\n                  <![endif]-->\n              </td>\n            </tr>\n          </tbody>\n        </table>\n      </div>"
	secondPart := "      <!--[if mso | IE]>\n            </td>\n          </tr>\n        </table>\n        \n        <table\n           align=\"center\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" class=\"\" style=\"width:600px;\" width=\"600\"\n        >\n          <tr>\n            <td style=\"line-height:0px;font-size:0px;mso-line-height-rule:exactly;\">\n        <![endif]-->\n      <div style=\"background:#ffffff;background-color:#ffffff;margin:0px auto;max-width:600px;\">\n        <table align=\"center\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" role=\"presentation\" style=\"background:#ffffff;background-color:#ffffff;width:100%;\">\n          <tbody>\n            <tr>\n              <td style=\"direction:ltr;font-size:0px;padding:0px;text-align:center;\">\n                <!--[if mso | IE]>\n                    <table role=\"presentation\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\">\n                  \n              <tr>\n                <td\n                   class=\"\" width=\"600px\"\n                >\n            \n        <table\n           align=\"center\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" class=\"\" style=\"width:600px;\" width=\"600\"\n        >\n          <tr>\n            <td style=\"line-height:0px;font-size:0px;mso-line-height-rule:exactly;\">\n        <![endif]-->\n                <div style=\"margin:0px auto;max-width:600px;\">\n                  <table align=\"center\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" role=\"presentation\" style=\"width:100%;\">\n                    <tbody>\n                      <tr>\n                        <td style=\"direction:ltr;font-size:0px;padding:20px 0;text-align:center;\">\n                          <!--[if mso | IE]>\n                    <table role=\"presentation\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\">\n                  \n          <tr>\n        \n              <td\n                 class=\"\" style=\"vertical-align:top;width:600px;\"\n              >\n            <![endif]-->\n                          <div class=\"mj-column-per-100 mj-outlook-group-fix\" style=\"font-size:0px;text-align:left;direction:ltr;display:inline-block;vertical-align:top;width:100%;\">\n                            <table border=\"0\" cellpadding=\"0\" cellspacing=\"0\" role=\"presentation\" style=\"vertical-align:top;\" width=\"100%\">\n                              <tbody><tr>\n                                <td align=\"left\" style=\"font-size:0px;padding:10px 25px;word-break:break-word;\">\n                                    <div style=\"font-family:Roboto, Helvetica, Arial, sans-serif;font-size:14px;font-weight:400;line-height:20px;text-align:left;color:#54595f;\">\n                                      <p style=\"margin-bottom: 0;\">Si tienes algun problema no dudes en contactar con nuestro servicio técnico en: <a href=\"#\" style=\"color: #4793d7; text-decoration: none;\">info@clinicafortia.com</a></p>\n                                    </div>\n                                </td>\n                              </tr>\n                            </tbody></table>\n                          </div>\n                          <!--[if mso | IE]>\n              </td>\n            \n          </tr>\n        \n                    </table>\n                  <![endif]-->\n                        </td>\n                      </tr>\n                    </tbody>\n                  </table>\n                </div>\n                <!--[if mso | IE]>\n            </td>\n          </tr>\n        </table>\n        \n                </td>\n              </tr>\n            \n                    </table>\n                  <![endif]-->\n              </td>\n            </tr>\n          </tbody>\n        </table>\n      </div>\n        <!--[if mso | IE]>\n              </td>\n            </tr>\n          </table>\n          \n          <table\n             align=\"center\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" class=\"\" style=\"width:600px;\" width=\"600\"\n          >\n            <tr>\n              <td style=\"line-height:0px;font-size:0px;mso-line-height-rule:exactly;\">\n          <![endif]-->\n          <div style=\"background:#ffffff;background-color:#ffffff;margin:0px auto;border-radius:0 0 4px 4px;max-width:600px;\">\n              <table align=\"center\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" role=\"presentation\" style=\"background:#ffffff;background-color:#ffffff;width:100%;border-radius:0 0 4px 4px;\">\n                <tbody>\n                  <tr>\n                    <td style=\"direction:ltr;font-size:0px;padding:20px 0;text-align:center;\">\n                      <!--[if mso | IE]>\n                          <table role=\"presentation\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\">\n                        \n                <tr>\n              \n                    <td\n                       class=\"\" style=\"vertical-align:top;width:600px;\"\n                    >\n                  <![endif]-->\n                      <div class=\"mj-column-per-100 mj-outlook-group-fix\" style=\"font-size:0px;text-align:left;direction:ltr;display:inline-block;vertical-align:top;width:100%;\">\n                        <table border=\"0\" cellpadding=\"0\" cellspacing=\"0\" role=\"presentation\" style=\"vertical-align:top;\" width=\"100%\">\n                          <td style=\"font-size:0px;padding:10px 25px;word-break:break-word;\">\n                            <p style=\"border-top: dashed 1px lightgrey; font-size: 1px; margin: 0px auto; width: 100%;\">\n                            </p>\n                            <!--[if mso | IE]>\n              <table\n                 align=\"center\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" style=\"border-top:dashed 1px lightgrey;font-size:1px;margin:0px auto;width:550px;\" role=\"presentation\" width=\"550px\"\n              >\n                <tr>\n                  <td style=\"height:0;line-height:0;\">\n      \n                  </td>\n                </tr>\n              </table>\n            <![endif]-->\n                          </td>\n                        </tr>\n                          \n                          \n                          \n                          <tbody><tr>\n                            <td align=\"center\" style=\"font-size:0px;padding:10px 25px;word-break:break-word;\">\n                              <div style=\"font-family:Roboto, Helvetica, Arial, sans-serif;font-size:14px;font-weight:400;line-height:20px;text-align:center;color:#93999f;\"><a class=\"footer-link\" href=\"#\" style=\"color: #4793d7; text-decoration: none;\">info@clinicafortia.com</a> <br /><a class=\"footer-link\" href=\"https://clinicafortia.com/\" style=\"color: #4793d7; text-decoration: none;\">https://clinicafortia.com/</a></div>\n                            </td>\n                          </tr>\n                          <tr>\n                            <td align=\"center\" style=\"font-size:0px;padding:10px 25px;word-break:break-word;\">\n                              <!--[if mso | IE]>\n              <table\n                 align=\"center\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" role=\"presentation\"\n              >\n                <tr>\n              \n                      <td>\n                    <![endif]-->\n                              <table align=\"center\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" role=\"presentation\" style=\"float:none;display:inline-table;\">\n                                <tbody><tr>\n                                  <td style=\"padding:4px;\">\n                                    <table border=\"0\" cellpadding=\"0\" cellspacing=\"0\" role=\"presentation\" style=\"border-radius:3px;width:24px;\">\n                                      <tbody><tr>\n                                        <td style=\"font-size:0;height:24px;vertical-align:middle;width:24px;\">\n                                          <a href=\"#\" target=\"_blank\" style=\"color: #4793d7; text-decoration: none;\">\n                                            <img alt=\"twitter-logo\" height=\"24\" src=\"https://codedmails.com/images/social/color/twitter-logo-transparent.png\" style=\"border-radius:3px;display:block;\" width=\"24\" />\n                                          </a>\n                                        </td>\n                                      </tr>\n                                    </tbody></table>\n                                  </td>\n                                </tr>\n                              </tbody></table>\n                              <!--[if mso | IE]>\n                      </td>\n                    \n                      <td>\n                    <![endif]-->\n                              <table align=\"center\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" role=\"presentation\" style=\"float:none;display:inline-table;\">\n                                <tbody><tr>\n                                  <td style=\"padding:4px;\">\n                                    <table border=\"0\" cellpadding=\"0\" cellspacing=\"0\" role=\"presentation\" style=\"border-radius:3px;width:24px;\">\n                                      <tbody><tr>\n                                        <td style=\"font-size:0;height:24px;vertical-align:middle;width:24px;\">\n                                          <a href=\"#\" target=\"_blank\" style=\"color: #4793d7; text-decoration: none;\">\n                                            <img alt=\"facebook-logo\" height=\"24\" src=\"https://codedmails.com/images/social/color/facebook-logo-transparent.png\" style=\"border-radius:3px;display:block;\" width=\"24\" />\n                                          </a>\n                                        </td>\n                                      </tr>\n                                    </tbody></table>\n                                  </td>\n                                </tr>\n                              </tbody></table>\n                              <!--[if mso | IE]>\n                      </td>\n                    \n                      <td>\n                    <![endif]-->\n                              <table align=\"center\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" role=\"presentation\" style=\"float:none;display:inline-table;\">\n                                <tbody><tr>\n                                  <td style=\"padding:4px;\">\n                                    <table border=\"0\" cellpadding=\"0\" cellspacing=\"0\" role=\"presentation\" style=\"border-radius:3px;width:24px;\">\n                                      <tbody><tr>\n                                        <td style=\"font-size:0;height:24px;vertical-align:middle;width:24px;\">\n                                          <a href=\"https://www.instagram.com/fisioterapiamovimientoysalud/\" target=\"_blank\" style=\"color: #4793d7; text-decoration: none;\">\n                                            <img alt=\"instagram-logo\" height=\"24\" src=\"https://codedmails.com/images/social/color/insta-logo-transparent.png\" style=\"border-radius:3px;display:block;\" width=\"24\" />\n                                          </a>\n                                        </td>\n                                      </tr>\n                                    </tbody></table>\n                                  </td>\n                                </tr>\n                              </tbody></table>\n                              <!--[if mso | IE]>\n                      </td>\n                    \n                      <td>\n                    <![endif]-->\n                              <table align=\"center\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" role=\"presentation\" style=\"float:none;display:inline-table;\">\n                                <tbody><tr>\n                                  <td style=\"padding:4px;\">\n                                    <table border=\"0\" cellpadding=\"0\" cellspacing=\"0\" role=\"presentation\" style=\"border-radius:3px;width:24px;\">\n                                      <tbody><tr>\n                                        <td style=\"font-size:0;height:24px;vertical-align:middle;width:24px;\">\n                                          <a href=\"https://www.youtube.com/channel/UChVWPZblpBRMpkT3b0Ukiqw\" target=\"_blank\" style=\"color: #4793d7; text-decoration: none;\">\n                                            <img alt=\"youtube-logo\" height=\"24\" src=\"https://codedmails.com/images/social/color/youtube-logo-transparent.png\" style=\"border-radius:3px;display:block;\" width=\"24\" />\n                                          </a>\n                                        </td>\n                                      </tr>\n                                    </tbody></table>\n                                  </td>\n                                </tr>\n                              </tbody></table>\n                              <!--[if mso | IE]>\n                      </td>\n                    \n                  </tr>\n                </table>\n              <![endif]-->\n                            </td>\n                          </tr>\n                        </tbody></table>\n                      </div>\n                      <!--[if mso | IE]>\n                    </td>\n                  \n                </tr>\n              \n                          </table>\n                        <![endif]-->\n                    </td>\n                  </tr>\n                  <tr>\n                    <td align=\"center\" style=\"font-size:0px;padding:10px 25px;word-break:break-word;\">\n                      <div style=\"font-family:Roboto, Helvetica, Arial, sans-serif;font-size:14px;font-weight:400;line-height:20px;text-align:center;color:#93999f;\">© 2021 Clínica FORTIA, Derechos Reservados \n                    </td>\n                  </tr>\n                </tbody>\n              </table>\n            </div>\n            <!--[if mso | IE]>\n                  </td>\n                </tr>\n              </table>\n              \n              <table\n                 align=\"center\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" class=\"\" style=\"width:600px;\" width=\"600\"\n              >\n                <tr>\n                  <td style=\"line-height:0px;font-size:0px;mso-line-height-rule:exactly;\">\n              <![endif]-->\n            <div style=\"margin:0px auto;max-width:600px;\">\n              <table align=\"center\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" role=\"presentation\" style=\"width:100%;\">\n                <tbody>\n                  <tr>\n                    <td style=\"direction:ltr;font-size:0px;padding:20px 0;text-align:center;\">\n                      <!--[if mso | IE]>\n                          <table role=\"presentation\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\">\n                        \n                <tr>\n              \n                    <td\n                       class=\"\" style=\"vertical-align:top;width:600px;\"\n                    >\n                  <![endif]-->\n                      <div class=\"mj-column-per-100 mj-outlook-group-fix\" style=\"font-size:0px;text-align:left;direction:ltr;display:inline-block;vertical-align:top;width:100%;\">\n                        <table border=\"0\" cellpadding=\"0\" cellspacing=\"0\" role=\"presentation\" style=\"vertical-align:top;\" width=\"100%\">\n                          <tbody><tr>\n                            <td style=\"font-size:0px;word-break:break-word;\">\n                              <!--[if mso | IE]>\n            \n                <table role=\"presentation\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\"><tr><td height=\"1\" style=\"vertical-align:top;height:1px;\">\n              \n            <![endif]-->\n                              <div style=\"height:1px;\">   </div>\n                              <!--[if mso | IE]>\n            \n                </td></tr></table>\n              \n            <![endif]-->\n                            </td>\n                          </tr>\n                        </tbody></table>\n                      </div>\n                      <!--[if mso | IE]>\n                    </td>\n                  \n                </tr>\n              \n                          </table>\n                        <![endif]-->\n                    </td>\n                  </tr>\n                </tbody>\n              </table>\n            </div>\n            <!--[if mso | IE]>\n                  </td>\n                </tr>\n              </table>\n              <![endif]-->\n          </div>\n  \n  \n  </body></html>"
	for (rows.Next()){
		rows.Scan(&date, &name)
		dateAsSomething := time.Unix(date, 0)
		minute := strconv.Itoa(dateAsSomething.Minute())
		if ( dateAsSomething.Minute() >= 0 && dateAsSomething.Minute()<10){
			minute = "0" + minute
		}
		middlePart += "      <!--[if mso | IE]>\n            </td>\n          </tr>\n        </table>\n        \n        <table\n           align=\"center\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" class=\"\" style=\"width:600px;\" width=\"600\"\n        >\n          <tr>\n            <td style=\"line-height:0px;font-size:0px;mso-line-height-rule:exactly;\">\n        <![endif]-->\n      <div style=\"background:#ffffff;background-color:#ffffff;margin:0px auto;max-width:600px;\">\n        <table align=\"center\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" role=\"presentation\" style=\"background:#ffffff;background-color:#ffffff;width:100%;\">\n          <tbody>\n            <tr>\n              <td style=\"direction:ltr;font-size:0px;padding:0 20px;text-align:center;\">\n                <!--[if mso | IE]>\n                    <table role=\"presentation\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\">\n                  \n              <tr>\n                <td\n                   class=\"\" width=\"600px\"\n                >\n            \n        <table\n           align=\"center\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" class=\"\" style=\"width:560px;\" width=\"560\"\n        >\n          <tr>\n            <td style=\"line-height:0px;font-size:0px;mso-line-height-rule:exactly;\">\n        <![endif]-->\n\n                <div style=\"background:#F5914D;background-color:#F5914D;margin:0px auto;border-radius:4px;max-width:560px;\">\n                  <table align=\"center\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" role=\"presentation\" style=\"background:#F5914D;background-color:#F5914D;width:100%;border-radius:4px;\">\n                    <tbody>\n                      <tr>\n                        <td style=\"direction:ltr;font-size:0px;padding:5px;text-align:center;\">\n                          <!--[if mso | IE]>\n                    <table role=\"presentation\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\">\n                  \n          <tr>\n        \n              <td\n                 class=\"\" style=\"vertical-align:top;width:550px;\"\n              >\n            <![endif]-->\n                          <div class=\"mj-column-per-100 mj-outlook-group-fix\" style=\"font-size:0px;text-align:left;direction:ltr;display:inline-block;vertical-align:top;width:100%;\">\n                            <table border=\"0\" cellpadding=\"0\" cellspacing=\"0\" role=\"presentation\" style=\"vertical-align:top;\" width=\"100%\">\n                              <tbody><tr>\n                                <td align=\"left\" style=\"font-size:0px;padding:10px 25px;word-break:break-word;\">\n                                  <div style=\"font-family:Roboto, Helvetica, Arial, sans-serif;font-size:24px;font-weight:400;line-height:30px;text-align:left;color:#ffffff;\">\n                                    <h2 style=\"margin: 0; font-size: 20px; font-weight: 300; line-height: 24px;\">Cita Fisioterapia</h2>\n                                  </div>\n                                </td>\n                              </tr>\n                              <tr>\n                                <td align=\"left\" style=\"font-size:0px;padding:10px 25px;word-break:break-word;\">\n                                  <div style=\"font-family:Roboto, Helvetica, Arial, sans-serif;font-size:14px;font-weight:400;line-height:20px;text-align:left;color:#ffffff;\">\n                                    <p style=\"margin: 0;\">Tienes una cita pendiente con " + name + ".<br>Día: " + strconv.Itoa(dateAsSomething.Day()) + " de " + months[dateAsSomething.Month()-1] + "<br>Hora: " + strconv.Itoa(dateAsSomething.Hour() - 1)  + ":" + minute + " </p>\n                                  </div>\n                                </td>\n                              </tr>\n                              <tr>\n                                <td align=\"right\" vertical-align=\"middle\" style=\"font-size:0px;padding:10px 25px;word-break:break-word;\">\n                                  <table border=\"0\" cellpadding=\"0\" cellspacing=\"0\" role=\"presentation\" style=\"border-collapse:separate;line-height:100%;\">\n                                    <tbody><tr>\n                                      <td align=\"center\" bgcolor=\"#ffffff\" role=\"presentation\" style=\"border:none;border-radius:3px;cursor:auto;mso-padding-alt:10px 25px;background:#ffffff;\" valign=\"middle\">\n                                        <a href=\"http://localhost:8081/calendar\" style=\"display: inline-block; background: #ffffff; color: #000000; font-family: Roboto, Helvetica, Arial, sans-serif; font-size: 13px; font-weight: normal; line-height: 120%; margin: 0; text-decoration: none; text-transform: none; padding: 10px 25px; mso-padding-alt: 0px; border-radius: 3px;\" target=\"_blank\">\n                                          <strong>VER MI CALENDARIO</strong>\n                                        </a>\n                                      </td>\n                                    </tr>\n                                  </tbody></table>\n                                </td>\n                              </tr>\n                            </tbody></table>\n                          </div>\n                          <!--[if mso | IE]>\n              </td>\n            \n          </tr>\n        \n                    </table>\n                  <![endif]-->\n                        </td>\n                      </tr>\n                    </tbody>\n                  </table>\n                </div>\n                <!--[if mso | IE]>\n            </td>\n          </tr>\n        </table>\n        \n                </td>\n              </tr>\n            \n                    </table>\n                  <![endif]-->\n              </td>\n            </tr>\n          </tbody>\n        </table>\n      </div>"
	}
	if (middlePart!="") {
		finalPart := firstPart + middlePart + secondPart

		// Sender data.
		from := secret.Email
		password := secret.Password
		// smtp server configuration.
		smtpServer := smtpServer{host: "smtp.gmail.com", port: "587"}
		// Authentication.
		auth := smtp.PlainAuth("", from, password, smtpServer.host)

		emailToSend := []string{email}

		r := NewRequest(emailToSend, "BIENVENIDO", "")
		r.body = finalPart

		mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
		to := "To: " + emailToSend[0] + " \n"
		subject := "Subject: " + r.subject + "!\n"
		msg := []byte(to + subject + mime + "\n" + r.body)

		// Sending email.
		err := smtp.SendMail(smtpServer.Address(), auth, "pruebaapifisio@gmail.com", r.to, msg)
		if err != nil {
			lib.TerminalLogger.Error("Problems sending the email: " + err.Error())
			lib.DocuLogger.Error("Problems sending the email: " + err.Error())
		} else{
			lib.TerminalLogger.Trace("Emails send succesfully")
			lib.DocuLogger.Trace("Emails send succesfully")
		}
	} else {
		lib.TerminalLogger.Trace("No appointments today for: ", email)
		lib.DocuLogger.Trace("No appointments today for: ", email)
	}
}