// togglwebhook.go
package togglwebhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ricotheque/webhooks-test/config"
	"github.com/ricotheque/webhooks-test/safelog"
)

func ValidateWebhook(secret, signature string, body []byte) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(expectedMAC, []byte(signature))
}

func HandleTogglWebhook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Save the payload as a string
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		payloadAsString := string(body)

		// If the payload is a webhooks subscription, don't do the secret check
		if isSubscription(payloadAsString) {

		} else {
			secret := config.Get("togglWebhooks.secret").(string)
			if secret == "" {
				http.Error(w, "togglWebhooks.secret not set on config.yaml", http.StatusInternalServerError)
				return
			}

			signature := r.Header.Get("x-webhook-signature-256")
			fmt.Println(signature)
			if !ValidateWebhook(secret, signature, body) {
				http.Error(w, "Invalid signature", http.StatusUnauthorized)
				return
			}
		}

		// Process payload
		parsePayload(payloadAsString)

		w.WriteHeader(http.StatusOK)
	}
}

func isSubscription(payload string) bool {
	type payloadJSON struct {
		Metadata struct {
			RequestType string `json:"request_type"`
		} `json:"metadata"`
		Payload           string `json:"payload"`
		ValidationCodeUrl string `json:"validation_code_url"`
	}

	var data payloadJSON

	if err := json.Unmarshal([]byte(payload), &data); err != nil {
		// Invalid payload = this isn't a subscription attempt
		return false
	}

	if data.Metadata.RequestType == "POST" && data.Payload == "ping" {
		fmt.Printf("Subscription payload received. Validation URL %s\n", data.ValidationCodeUrl)
		return true
	}

	return false
}

func parsePayload(payload string) {
	type Meta struct {
		Model       string `json:"model"`
		Action      string `json:"action"`
		EventUserID string `json:"event_user_id"`
	}
	type Event struct {
		Timestamp string `json:"timestamp"`
		EventID   string `json:"event_id"`
		Payload   string `json:"payload"`
		Metadata  *Meta  `json:"metadata"`
	}

	event := Event{}
	err := json.Unmarshal([]byte(payload), &event)
	if err != nil {
		panic(err)
	}

	// Defaulting missing fields
	if event.Metadata == nil {
		event.Metadata = &Meta{}
	}

	fmt.Printf("Timestamp: %s\n", event.Timestamp)
	fmt.Printf("Event ID: %s\n", event.EventID)
	fmt.Printf("Metadata Model: %s\n", event.Metadata.Model)
	fmt.Printf("Metadata Action: %s\n", event.Metadata.Action)
	fmt.Printf("Metadata Event User ID: %s\n", event.Metadata.EventUserID)
	fmt.Printf("Payload: %s\n", event.Payload)

	safelog.Log(fmt.Sprintf(
		"%s\t%s\t%s\t%s\t%s\t%s\n",
		event.Timestamp,
		event.EventID,
		event.Metadata.Model,
		event.Metadata.Action,
		event.Metadata.EventUserID,
		event.Payload,
	),
	)
}
