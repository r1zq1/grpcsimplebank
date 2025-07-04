package mail

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/jordan-wright/email"
)

func TestKirimEmail(t *testing.T) {
	// Load .env config
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("cannot load .env: %v", err)
	}

	// Ambil dari env
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USERNAME")
	smtpPass := os.Getenv("SMTP_PASSWORD")
	smtpFrom := os.Getenv("SMTP_FROM")
	// target := os.Getenv("SMTP_TEST_RECEIVER") // misal: test@example.com

	if smtpFrom == "" {
		log.Fatalf("SMTP_FROM or SMTP_TEST_RECEIVER is missing")
	}

	// Buat email
	e := email.NewEmail()
	e.From = smtpFrom
	e.To = []string{smtpFrom}
	e.Subject = "Test Email from GRPC Bank"
	e.Text = []byte("Halo! Ini adalah pengujian kirim email dari sistem GRPC.")

	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	fmt.Printf("Sending email from: %s to %s...\n", smtpFrom, smtpFrom)
	err = e.Send(addr, auth)
	if err != nil {
		log.Fatalf("Failed to send email: %v", err)
	}

	fmt.Println("✅ Email sent successfully!")
}
