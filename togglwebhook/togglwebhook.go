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

// Event represents the top-level object structure
type Event struct {
	Timestamp string `json:"timestamp"`
	EventID   string
	Payload   string
	CreatorID string
	Metadata  *Metadata `json:"metadata"`
}

// Metadata captures the nested metadata structure
type Metadata struct {
	Model       string
	Action      string
	EventUserID string
}

// UnmarshalJSON implements the custom JSON unmarshalling for Event
func (e *Event) UnmarshalJSON(data []byte) error {
	type Alias Event

	aux := &struct {
		EventID   int64           `json:"event_id"`
		CreatorID int64           `json:"creator_id"`
		Payload   json.RawMessage `json:"payload"`
		*Alias
	}{
		Alias: (*Alias)(e),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	e.EventID = strconv.FormatInt(aux.EventID, 10)
	e.CreatorID = strconv.FormatInt(aux.CreatorID, 10)
	e.Payload = string(aux.Payload)
	return nil
}

// UnmarshalJSON implements the custom JSON unmarshalling for Metadata
func (m *Metadata) UnmarshalJSON(data []byte) error {
	type Alias Metadata

	aux := &struct {
		Model       json.RawMessage  `json:"model"`
		Action      *json.RawMessage `json:"action"`
		EventUserID json.RawMessage  `json:"event_user_id"`
		*Alias
	}{
		Alias: (*Alias)(m),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Model != nil {
		m.Model = string(aux.Model)
	}
	if aux.Action != nil {
		m.Action = string(*aux.Action)
	}
	m.EventUserID = string(aux.EventUserID)

	return nil
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

		// Only do the secret check if the payload is not a subscription
		if !isSubscription(payloadAsString) {
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
		ParsePayload(payloadAsString)

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

func ParsePayload(payload string) {
	event := Event{}
	err := json.Unmarshal([]byte(payload), &event)
	if err != nil {
		panic(err)
	}

	// Use default values in case metadata or its fields are missing
	var model, action, eventUserID string
	if event.Metadata != nil {
		model = stripQuotes(event.Metadata.Model)
		action = stripQuotes(event.Metadata.Action)
		eventUserID = stripQuotes(event.Metadata.EventUserID)
	}

	// Make payloads one-liners
	compactPayload, payloadErr := compactJSON(event.Payload)
	if payloadErr != nil {
		fmt.Println("Error compacting payload:", payloadErr)
		return
	}
	event.Payload = compactPayload

	// Combined print using fmt.Sprintf
	combinedOutput := fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s",
		event.Timestamp, event.EventID, event.CreatorID,
		model, action, eventUserID, event.Payload)

	fmt.Println(combinedOutput)
	safelog.Log(combinedOutput)
}

func compactJSON(jsonStr string) (string, error) {
	var jsonObj interface{}

	// Unmarshal the JSON string into an interface
	if err := json.Unmarshal([]byte(jsonStr), &jsonObj); err != nil {
		return "", err
	}

	// Marshal the interface back into a compact JSON string
	compactJSONBytes, err := json.Marshal(jsonObj)
	if err != nil {
		return "", err
	}

	// Convert the compact JSON bytes to a string
	compactStr := string(compactJSONBytes)

	// If the compacted JSON is a string, it will have surrounding quotes.
	// Here, we safely remove those quotes.
	compactStr = stripQuotes(compactStr)

	return compactStr, nil
}

func stripQuotes(s string) string {
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	if s == "" {
		s = "-"
	}
	return s
}
