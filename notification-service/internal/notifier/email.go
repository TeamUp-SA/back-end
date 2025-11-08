package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	mailerSendEndpoint    = "https://api.mailersend.com/v1/email"
	hardcodedEmailAddress = "dalai2547@gmail.com"
)

type mailerSendAddress struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type mailerSendEmailRequest struct {
	From     mailerSendAddress   `json:"from"`
	To       []mailerSendAddress `json:"to"`
	Subject  string              `json:"subject"`
	Text     string              `json:"text,omitempty"`
	HTML     string              `json:"html,omitempty"`
	ReplyTo  *mailerSendAddress  `json:"reply_to,omitempty"`
	Personal string              `json:"personalization,omitempty"`
}

var httpClient = &http.Client{
	Timeout: 15 * time.Second,
}

func SendEmail(_ string, subject, message string) {
	apiKey := strings.TrimSpace(os.Getenv("MAILERSEND_API_KEY"))
	if apiKey == "" {
		log.Println("MAILERSEND_API_KEY is not set; skipping email send")
		return
	}

	fromEmail := strings.TrimSpace(os.Getenv("MAILERSEND_FROM_EMAIL"))
	if fromEmail == "" {
		fromEmail = hardcodedEmailAddress
	}

	emailSubject := subject
	if strings.TrimSpace(emailSubject) == "" {
		emailSubject = "Notification"
	}

	payload := mailerSendEmailRequest{
		From: mailerSendAddress{
			Email: fromEmail,
		},
		To: []mailerSendAddress{
			{Email: hardcodedEmailAddress},
		},
		Subject: emailSubject,
		Text:    message,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal MailerSend payload: %v", err)
		return
	}

	req, err := http.NewRequest(http.MethodPost, mailerSendEndpoint, bytes.NewReader(body))
	if err != nil {
		log.Printf("Failed to create MailerSend request: %v", err)
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("MailerSend request error: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		responseBody, _ := io.ReadAll(resp.Body)
		log.Printf("MailerSend request failed: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(responseBody)))
		return
	}

	log.Printf("[EMAIL] to=%s subject=%s message=%s sent successfully via MailerSend API!", hardcodedEmailAddress, emailSubject, message)
}
