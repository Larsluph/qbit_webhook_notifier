package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"qbit_webhook/qbit"
)

var client = &http.Client{}

func TriggerWebhook(triggerType string, hash string) {
	qbit.Login()
	torrentProps := qbit.GetTorrentProperties(hash)

	var payload WebhookPayload
	if triggerType == "added" {
		payload = GenerateAddedEmbed(torrentProps)
	} else if triggerType == "completed" {
		payload = GenerateCompletedEmbed(torrentProps)
	}

	SendWebhook(payload)
}

func SendWebhook(payload WebhookPayload) {
	var webhookURL = os.Getenv("WEBHOOK_URL")

	reqBodyBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, webhookURL, bytes.NewBuffer(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Webhook status code:", resp.StatusCode)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)
}
