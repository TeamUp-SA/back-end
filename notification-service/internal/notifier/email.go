package notifier

import "log"

func SendEmail(to, message string) {
    // Mock implementation – replace with SMTP or SendGrid API
    log.Printf("[EMAIL] to=%s message=%s", to, message)
}
