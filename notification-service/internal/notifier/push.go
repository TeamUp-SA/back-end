package notifier

import "log"

func SendPush(to, message string) {
    log.Printf("[Push] to=%s message=%s", to, message)
}
