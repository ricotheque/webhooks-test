// main.go
package main

import (
	"log"
	"net/http"

	"github.com/ricotheque/webhooks-test/togglwebhook"
)

func main() {
	http.HandleFunc("/ttwh", togglwebhook.HandleTogglWebhook())

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
