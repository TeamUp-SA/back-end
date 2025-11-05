package notifier

import (
	"fmt"
	"log"
	"net/http"
)

func SendPush(title, message string) {
	url := fmt.Sprintf("http://host.docker.internal:9999/notify?title=%s&message=%s", title, message)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("❌ Failed to send notification to host:", err)
		return
	}
	defer resp.Body.Close()

	log.Println("📡 Sent notification request to host.")
}
