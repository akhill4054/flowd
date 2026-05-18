package backend

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("MetaWebhookHandler", MetaWebhookHandler)
}

// MetaWebhookHandler handles incoming Meta (WhatsApp) webhook calls
func MetaWebhookHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Handle GET request (Webhook Verification)
	if r.Method == http.MethodGet {
		mode := r.URL.Query().Get("hub.mode")
		token := r.URL.Query().Get("hub.verify_token")
		challenge := r.URL.Query().Get("hub.challenge")

		verifyToken := os.Getenv("META_WA_WEBHOOK_VERIFY_TOKEN")
		if verifyToken == "" {
			log.Println("Warning: META_WA_WEBHOOK_VERIFY_TOKEN is not set")
		}

		if mode == "subscribe" && token == verifyToken {
			log.Println("Webhook Verified Successfully!")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(challenge))
			return
		}
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// 2. Handle POST request (Incoming Messages/Events)
	if r.Method == http.MethodPost {
		// Read and restore the body to log it without breaking downstream logic
		bodyBytes, _ := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		log.Printf("Received Webhook: %s", string(bodyBytes))

		// IMPORTANT: Always return 200 OK immediately to Meta
		// to avoid retries and timeouts.
		w.WriteHeader(http.StatusOK)
	}
}
