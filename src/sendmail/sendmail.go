package sendmail

import (
	// "fmt"
	"net/smtp"
	"strings"
)

// SendToMail : send mail function
func SendToMail(user, password, host, to, subject, body, mailtype string) error {
	hp := strings.Split(host, ":")
	// fmt.Println(hp)
	auth := smtp.PlainAuth("", user, password, hp[0])
	// fmt.Println(auth)
	var contentType string
	if mailtype == "html" {
		contentType = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		contentType = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	msg := []byte("to: " + to + "\r\nfrom: goalive <" + user + ">\r\nsubject: " + subject + "\r\n" + contentType + "\r\n\r\n" + body)
	sendTo := strings.Split(to, ",")
	// fmt.Println(string(msg))
	// fmt.Println(host, auth, user, sendTo, msg)
	err := smtp.SendMail(host, auth, user, sendTo, msg)
	return err
}
