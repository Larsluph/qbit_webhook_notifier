package discord

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"qbit_webhook/qbit"
	"strings"
)

var cookieJar, _ = cookiejar.New(nil)
var client = &http.Client{
	Jar: cookieJar,
}

func TriggerWebhook(triggerType string, hash string) {
	loginForm := url.Values{}
	loginForm.Set("username", os.Getenv("QBIT_USERNAME"))
	loginForm.Set("password", os.Getenv("QBIT_PASSWORD"))

	// Encode the form data as URL-encoded form data
	body := strings.NewReader(loginForm.Encode())

	var prefix = os.Getenv("QBIT_BASE") + "/api/v2"
	var authURL = prefix + "/auth/login"
	var propsURL = prefix + "/torrents/properties"

	// Send the login request using the HTTP client
	rLogin, err := http.NewRequest(http.MethodPost, authURL, body)
	if err != nil {
		panic(err)
	}
	rLogin.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(rLogin)
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Login success!")

		// Read the request body
		body, err := io.ReadAll(rLogin.Body)
		if err != nil {
			panic(err)
		}

		// Print the request body to the console
		fmt.Println(string(body))
	} else {
		panic(errors.New("login failed"))
	}

	propsForm := url.Values{}
	propsForm.Set("hash", hash)

	// Send the props request using the HTTP client
	rTorrent, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?%s", propsURL, propsForm.Encode()), nil)
	if err != nil {
		panic(err)
	}
	resp, err = client.Do(rTorrent)
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		panic(errors.New(fmt.Sprintf("request unsuccessful with status: %d", resp.StatusCode)))
	}

	var torrentProps qbit.TorrentProps
	err = json.NewDecoder(resp.Body).Decode(&torrentProps)
	if err != nil {
		panic(err)
	}

	var payload WebhookPayload
	if triggerType == "added" {
		payload = GenerateAddedEmbed(torrentProps)
	} else if triggerType == "completed" {
		payload = GenerateCompletedEmbed(torrentProps)
	}

	SendWebhook(payload)
}

func SendWebhook(payload WebhookPayload) {
	var webhookURL = "https://discord.com/api/webhooks/" + os.Getenv("WEBHOOK_ID") + "/" + os.Getenv("WEBHOOK_TOKEN")

	reqBodyBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, webhookURL, bytes.NewBuffer(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			//TODO: Add logger when unable to send payload to Discord
		}
	}(resp.Body)
}
