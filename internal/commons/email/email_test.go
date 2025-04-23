package email

import (
	"testing"

	"github.com/joho/godotenv"
)

func TestSendEmail(t *testing.T) {
	err := godotenv.Load("../../../.env.local")
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}
	to := "christopher.khant.work@gmail.com"

	subject := "Test Email from Go"
	body := `
			<p>Test</p>
			<p>Sent using GoMail</p>
			`

	err = SendEmail(to, subject, body)
	if err != nil {
		t.Errorf("Failed to send email: %v", err)
	}
}
