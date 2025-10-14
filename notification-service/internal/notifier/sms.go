package notifier

import "log"

func SendSMS(to, message string) {
    // Mock implementation – replace with SMTP or SendGrid API
    log.Printf("[SMS] to=%s message=%s", to, message)
}
