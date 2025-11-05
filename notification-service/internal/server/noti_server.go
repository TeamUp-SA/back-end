package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

func notify(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	message := r.URL.Query().Get("message")

	if title == "" || message == "" {
		http.Error(w, "missing title or message", http.StatusBadRequest)
		return
	}

	cmd := exec.Command("osascript", "-e", fmt.Sprintf(`display notification "%s" with title "%s"`, message, title))
	if err := cmd.Run(); err != nil {
		log.Println("❌ Failed to show notification:", err)
		http.Error(w, "failed to show notification", 500)
		return
	}

	log.Printf("✅ Notification: [%s] %s\n", title, message)
	fmt.Fprintln(w, "Notification sent!")
}

func main() {
	http.HandleFunc("/notify", notify)
	log.Println("Listening on http://localhost:9999/notify")
	log.Fatal(http.ListenAndServe(":9999", nil))
}
