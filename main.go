// main.go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ricotheque/webhooks-test/safelog"
	"github.com/ricotheque/webhooks-test/togglwebhook"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var k = koanf.New(".")

func main() {
	// Initialize the default logger
	if err := safelog.InitDefaultLogger("payloads.log"); err != nil {
		panic(fmt.Sprintf("Failed to initialize default logger: %v", err))
	}
	defer safelog.CloseDefaultLogger()

	// Load configuration
	if err := k.Load(file.Provider("config.yaml"), yaml.Parser()); err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	certFile := ""
	if certFile = k.String("certFile"); certFile == "" {
		log.Fatalf("error loading config: %v", "certFile missing from config.yaml")
	}

	keyFile := ""
	if keyFile = k.String("keyFile"); keyFile == "" {
		log.Fatalf("error loading config: %v", "keyFile missing from config.yaml")
	}

	togglSecret := k.String("togglWebhooks.secret")

	http.HandleFunc("/ttwh", togglwebhook.HandleTogglWebhook(togglSecret))

	log.Println("Starting to listen on :443")
	log.Fatal(
		http.ListenAndServeTLS(":443", certFile, keyFile, nil),
	)
}
