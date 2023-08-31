package togglwebhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

type Payload struct {
	Event   string      `json:"event"`
	Data    interface{} `json:"data"`
	Version string      `json:"version"`
}

func ValidateWebhook(secret, signature string, body []byte) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(expectedMAC, []byte(signature))
}

func HandleTogglWebhook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		secret := os.Getenv("TOGGL_WEBHOOK_SECRET")
		if secret == "" {
			http.Error(w, "Webhook secret not set", http.StatusInternalServerError)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		signature := r.Header.Get("x-webhook-signature-256") // Adjust this if the header is different
		if !ValidateWebhook(secret, signature, body) {
			http.Error(w, "Invalid signature", http.StatusUnauthorized)
			return
		}

		var payload Payload
		err = json.Unmarshal(body, &payload)
		if err != nil {
			http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
			return
		}

		// Process payload
		w.WriteHeader(http.StatusOK)
	}
}
