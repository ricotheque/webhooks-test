// togglwebhook.go
package togglwebhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/ricotheque/webhooks-test/config"
	"github.com/ricotheque/webhooks-test/safelog"
)

type Event struct {
	Timestamp string    `json:"timestamp"`
	EventID   int64     `json:"event_id"`
	Payload   string    `json:"-"`
	Metadata  *Metadata `json:"metadata"`
}

type Metadata struct {
	Model       string `json:"model"`
	Action      string `json:"action"`
	EventUserID int64  `json:"-"`
}

func ValidateWebhook(secret string, signature string, payload string) bool {
	messageMAC, _ := hex.DecodeString(strings.TrimPrefix(signature, "sha256="))

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	expectedMAC := mac.Sum(nil)

	return hmac.Equal([]byte(messageMAC), expectedMAC)
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
				fmt.Println("togglWebhooks.secret not set on config.yaml")
				http.Error(w, "togglWebhooks.secret not set on config.yaml", http.StatusInternalServerError)
				return
			}

			signature := r.Header.Get("x-webhook-signature-256")
			if !ValidateWebhook(secret, signature, payloadAsString) {
				fmt.Println("Invalid signature")
				http.Error(w, "Invalid signature", http.StatusUnauthorized)
				return
			}
		}

		// Process payload
		fmt.Println(payloadAsString)
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
	event := Event{}
	err := json.Unmarshal([]byte(payload), &event)
	if err != nil {
		panic(err)
	}

	// Defaulting missing fields
	if event.Metadata == nil {
		event.Metadata = &Metadata{}
	}

	fmt.Printf("Timestamp: %s\n", event.Timestamp)
	fmt.Printf("Event ID: %d\n", event.EventID)
	fmt.Printf("Metadata Model: %s\n", event.Metadata.Model)
	fmt.Printf("Metadata Action: %s\n", event.Metadata.Action)
	fmt.Printf("Metadata Event User ID: %d\n", event.Metadata.EventUserID)
	fmt.Printf("Payload: %s\n", event.Payload)

	safelog.Log(fmt.Sprintf(
		"%s\t%d\t%s\t%s\t%d\t%s\n",
		event.Timestamp,
		event.EventID,
		event.Metadata.Model,
		event.Metadata.Action,
		event.Metadata.EventUserID,
		event.Payload,
	),
	)
}

func (e *Event) UnmarshalJSON(data []byte) error {
	type Alias Event
	aux := struct {
		Payload  interface{} `json:"payload"`
		Metadata struct {
			EventUserID interface{} `json:"event_user_id"`
			*Metadata
		} `json:"metadata"`
		*Alias
	}{
		Alias: (*Alias)(e),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Check if e.Metadata is nil and initialize if necessary
	if e.Metadata == nil {
		e.Metadata = &Metadata{}
	}

	// Convert Payload into string
	switch v := aux.Payload.(type) {
	case string:
		e.Payload = v
	case map[string]interface{}:
		bytes, err := json.Marshal(v)
		if err != nil {
			return err
		}
		e.Payload = string(bytes)
	default:
		return fmt.Errorf("unexpected type %T for Payload", v)
	}

	// Ensure EventUserID resolves to int64
	switch v := aux.Metadata.EventUserID.(type) {
	case float64:
		e.Metadata.EventUserID = int64(v)
	case string:
		val, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}
		e.Metadata.EventUserID = val
	default:
		return fmt.Errorf("unexpected type %T for EventUserID", v)
	}

	return nil
}
