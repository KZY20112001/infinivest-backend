package email

import (
	"fmt"
	"os"

	"gopkg.in/gomail.v2"
)

func SendEmail(to, subject, body string) error {

	from := os.Getenv("EMAIL_FROM")
	pass := os.Getenv("EMAIL_PASS")

	if from == "" || pass == "" {
		return fmt.Errorf("EMAIL_FROM or EMAIL_PASS is not set")
	}
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer("smtp.gmail.com", 587, from, pass)

	return d.DialAndSend(m)
}
