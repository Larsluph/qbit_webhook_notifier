package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"qbit_webhook/qbit"
	"strings"
)

func TriggerWebhook(triggerType string, hash string) {
	prefix := os.Getenv("QBIT_BASE") + "/api/v2"
	authURL := prefix + "/auth/login"
	propsURL := prefix + "/torrents/properties"
	webhookURL := "https://discord.com/api/webhooks/" + os.Getenv("WEBHOOK_ID") + "/" + os.Getenv("WEBHOOK_TOKEN")

	loginForm := url.Values{}
	loginForm.Set("username", os.Getenv("QBIT_USERNAME"))
	loginForm.Set("password", os.Getenv("QBIT_PASSWORD"))

	// Encode the form data as URL-encoded form data
	body := strings.NewReader(loginForm.Encode())

	// Send the login request using the HTTP client
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookieJar,
	}

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
			//TODO
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
		panic("Login failed.")
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
			//TODO
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("Request unsuccessful with status: %d", resp.StatusCode))
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

	reqBodyBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, webhookURL, bytes.NewBuffer(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			//TODO
		}
	}(resp.Body)
}
