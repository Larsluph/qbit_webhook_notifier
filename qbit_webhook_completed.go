package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
)

func webhook_completed(hash string) {
	prefix := os.Getenv("QBIT_BASE") + "/api/v2"
	authURL := prefix + "/auth/login"
	propsURL := prefix + "/torrents/properties"
	webhookURL := "https://discord.com/api/webhooks/" + os.Getenv("WEBHOOK_ID") + "/" + os.Getenv("WEBHOOK_TOKEN")

	login_form := url.Values{}
	login_form.Set("username", os.Getenv("QBIT_USERNAME"))
	login_form.Set("password", os.Getenv("QBIT_PASSWORD"))

	// Encode the form data as URL-encoded form data
	body := strings.NewReader(login_form.Encode())

	// Send the login request using the HTTP client
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookieJar,
	}

	r_login, err := http.NewRequest("POST", authURL, body)
	if err != nil {
		panic(err)
	}
	r_login.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(r_login)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Login success!")

		// Read the request body
		body, err := io.ReadAll(r_login.Body)
		if err != nil {
			panic(err)
		}

		// Print the request body to the console
		fmt.Println(string(body))
	} else {
		panic("Login failed.")
	}

	props_form := url.Values{}
	props_form.Set("hash", hash)

	// Send the props request using the HTTP client
	r_torrent, err := http.NewRequest("GET", fmt.Sprintf("%s?%s", propsURL, props_form.Encode()), nil)
	if err != nil {
		panic(err)
	}
	resp, err = client.Do(r_torrent)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("Request unsuccessful with status: %d", resp.StatusCode))
	}

	var torrent_props TorrentProps
	err = json.NewDecoder(resp.Body).Decode(&torrent_props)
	if err != nil {
		panic(err)
	}

	payload := DiscordWebhookPayload{
		Content: "",
		Embeds: []Embed{
			{
				Title:       "Torrent completed",
				Description: torrent_props.Name,
				Color:       0x00FF00,
				URL:         torrent_props.Comment,
				Fields: []Field{
					{
						Name:   "Total Size",
						Value:  FormatByteSize(torrent_props.TotalSize),
						Inline: true,
					},
					{
						Name:   "Torrent finished in",
						Value:  RelativeTimeElapsed(torrent_props.AdditionDate, torrent_props.CompletionDate),
						Inline: true,
					},
					{
						Name:   "Average DL Speed",
						Value:  FormatByteSpeed(torrent_props.DlSpeedAvg),
						Inline: true,
					},
				},
			},
		},
	}

	reqBodyBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", webhookURL, bytes.NewBuffer(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}
