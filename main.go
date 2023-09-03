// main.go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ricotheque/webhooks-test/config"
	"github.com/ricotheque/webhooks-test/safelog"
	"github.com/ricotheque/webhooks-test/togglwebhook"
)

func main() {
	// Initialize the default logger
	if err := safelog.InitDefaultLogger("payloads.log"); err != nil {
		panic(fmt.Sprintf("Failed to initialize default logger: %v", err))
	}
	defer safelog.CloseDefaultLogger()

	// Load configuration
	config.LoadConfig("./config.yaml")

	// Set up route
	http.HandleFunc("/ttwh", togglwebhook.HandleTogglWebhook())

	// Start receiving payloads
	certFile := config.Get("certFile").(string)
	keyFile := config.Get("keyFile").(string)
	log.Println("Listening on :443")
	log.Fatal(
		http.ListenAndServeTLS(":443", certFile, keyFile, nil),
	)
}
