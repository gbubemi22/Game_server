package utils

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err) // Changed from Fatal to Printf
	}
}

func SendMail(subject, body string, to []string) error {
	loadEnv()

	// Zoho SMTP configuration
	host := os.Getenv("ZOHO_SMTP_HOST") // smtp.zoho.com
	portStr := os.Getenv("ZOHO_SMTP_PORT") // 465 for SSL, 587 for TLS
	user := os.Getenv("ZOHO_SMTP_USER") // Your Zoho email address
	pass := os.Getenv("ZOHO_SMTP_PASSWORD") // Zoho app-specific password
	sender := os.Getenv("ZOHO_SMTP_SENDER") // Should match your Zoho email

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return fmt.Errorf("invalid port: %v", err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", sender)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	// Configure dialer with Zoho settings
	d := gomail.NewDialer(host, port, user, pass)
	
	// For Zoho, you might need to explicitly set TLS (especially if using port 587)
	d.TLSConfig = &tls.Config{
		ServerName: host,
	}

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("could not send email: %v", err)
	}

	return nil
}