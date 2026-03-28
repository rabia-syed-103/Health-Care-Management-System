package utils

import (
	"os"

	"gopkg.in/gomail.v2"
)

func SendEmail(to string, subject string, body string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", os.Getenv("EMAIL_USER"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(
		"smtp.gmail.com",
		587,
		os.Getenv("EMAIL_USER"),
		os.Getenv("EMAIL_PASS"),
	)

	return d.DialAndSend(m)
}
