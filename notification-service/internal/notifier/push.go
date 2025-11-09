package notifier

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func SendPush(title, message string) {
	addr := os.Getenv("NOTIFICATION_SERVICE_RPC_ADDR")
	if addr == "" {
		addr = "localhost:9999" // Default for local
	}
	// url := fmt.Sprintf("http://host.docker.internal:9999/notify?title=%s&message=%s", title, message)
	url := fmt.Sprintf("http://%s/notify?title=%s&message=%s", addr, title, message)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("❌ Failed to send notification to host:", err)
		return
	}
	defer resp.Body.Close()

	log.Println("📡 Sent notification request to host.")
}
