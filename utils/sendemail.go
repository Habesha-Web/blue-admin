package utils

import (
	"context"
	"fmt"
	"net/smtp"
	"os"
	"strings"

	"blue-admin.com/configs"
)

type Mail struct {
	Sender  string
	To      []string
	Subject string
	Body    string
}

func SendEmailConsumer(ctx context.Context, msg string, sub string, emails []string) {
	// email server
	host := configs.AppConfig.Get("MAIL_SERVER")
	port := configs.AppConfig.Get("MAIL_PORT")
	sender_email := configs.AppConfig.Get("MAIL_USERNAME")
	password := configs.AppConfig.Get("MAIL_PASSWORD")

	body_msg := fmt.Sprintf("<p><b>%s</b></p>", msg)
	request := Mail{
		Sender:  sender_email,
		To:      emails,
		Subject: sub,
		Body:    body_msg,
	}
	my_msg := BuildMessage(request)

	// We can't send strings directly in mail,
	// strings need to be converted into slice bytes
	msg_body := []byte(my_msg)

	// PlainAuth uses the given username and password to
	// authenticate to host and act as identity.
	// Usually identity should be the empty string,
	// to act as username.
	auth := smtp.PlainAuth("", sender_email, password, host)

	// SendMail uses TLS connection to send the mail
	// The email is sent to all address in the toList,
	// the body should be of type bytes, not strings
	// This returns error if any occurred.
	err := smtp.SendMail(host+":"+port, auth, sender_email, emails, msg_body)

	// handling the errors
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Successfully sent mail to all user in toList")
}

func BuildMessage(mail Mail) string {
	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	msg += fmt.Sprintf("From: %s\r\n", mail.Sender)
	msg += fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";"))
	msg += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	msg += fmt.Sprintf("\r\n%s\r\n", mail.Body)

	return msg
}
