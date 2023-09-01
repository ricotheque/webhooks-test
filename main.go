// main.go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ricotheque/webhooks-test/safelog"
	"github.com/ricotheque/webhooks-test/togglwebhook"
)

func main() {
	// Initialize the default logger
	if err := safelog.InitDefaultLogger("payloads.log"); err != nil {
		panic(fmt.Sprintf("Failed to initialize default logger: %v", err))
	}
	defer safelog.CloseDefaultLogger()

	http.HandleFunc("/ttwh", togglwebhook.HandleTogglWebhook())

	// Modify this line, provide the paths to your certificate and private key files
	log.Fatal(
		http.ListenAndServeTLS(":443", "/etc/letsencrypt/live/webhooks.mossesgeld.com/fullchain.pem", "/etc/letsencrypt/live/webhooks.mossesgeld.com/privkey.pem", nil),
	)

	log.Println("Server started on :443")
}
