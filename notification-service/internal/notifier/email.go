package notifier

import (
	"log"
	"net/smtp"
	"os"
)

func SendEmail(to, message string) {
	from := "MS_BBVFQ8@test-zkq340e8z7kgd796.mlsender.net"
	smtpHost := os.Getenv("MAILERSEND_SMTP_SERVER")
	smtpPort := os.Getenv("MAILERSEND_SMTP_PORT")
	smtpUser := os.Getenv("MAILERSEND_SMTP_USER")
	smtpPass := os.Getenv("MAILERSEND_SMTP_PASS")

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	msg := []byte("To: " + to + "\r\n" +
		"Subject: Notification\r\n" +
		"\r\n" +
		message + "\r\n")

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	if err != nil {
		log.Println("Failed to send email:", err)
		return
	}

	log.Printf("[EMAIL] to=%s message=%s sent successfully!", to, message)
}
