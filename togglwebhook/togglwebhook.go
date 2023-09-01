// togglwebhook.go
package togglwebhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"

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
		// secret := os.Getenv("TTWH_SECRET")
		// if secret == "" {
		// 	http.Error(w, "Webhook secret not set", http.StatusInternalServerError)
		// 	return
		// }

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		// signature := r.Header.Get("x-webhook-signature-256") // Adjust this if the header is different
		// if !ValidateWebhook(secret, signature, body) {
		// 	http.Error(w, "Invalid signature", http.StatusUnauthorized)
		// 	return
		// }

		// Now, 'body' contains the raw JSON payload as a string.
		payloadAsString := string(body)

		// Process payload (example: just print it for now)
		fmt.Println(payloadAsString)
		safelog.Log(payloadAsString)

		w.WriteHeader(http.StatusOK)
	}
}
