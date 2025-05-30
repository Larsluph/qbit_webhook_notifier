package qbit

import (
	//"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
)

var cookieJar, _ = cookiejar.New(nil)
var client = &http.Client{
	Jar: cookieJar,
	//Transport: &http.Transport{
	//	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	//},
}

func Login() {
	var authURL = os.Getenv("QBIT_BASE") + "/api/v2/auth/login"

	loginForm := url.Values{}
	loginForm.Set("username", os.Getenv("QBIT_USERNAME"))
	loginForm.Set("password", os.Getenv("QBIT_PASSWORD"))

	// Encode the form data as URL-encoded form data
	body := strings.NewReader(loginForm.Encode())

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

	fmt.Println("Login status code:", resp.StatusCode)
	if resp.StatusCode == http.StatusOK {
		// Read the request body
		body, err := io.ReadAll(rLogin.Body)
		if err != nil {
			panic(err)
		}

		// Print the request body to the console
		fmt.Println("Login body:", string(body))
	} else {
		panic(errors.New("login failed"))
	}
}

func GetTorrentProperties(hash string) TorrentProps {
	var propsURL = os.Getenv("QBIT_BASE") + "/api/v2/torrents/properties"

	propsForm := url.Values{}
	propsForm.Set("hash", hash)

	// Send the props request using the HTTP client
	rTorrent, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?%s", propsURL, propsForm.Encode()), nil)
	if err != nil {
		panic(err)
	}
	resp, err := client.Do(rTorrent)
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	fmt.Println("Torrent props status code:", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		panic(errors.New(fmt.Sprintf("request unsuccessful with status: %d", resp.StatusCode)))
	}

	var torrentProps TorrentProps
	err = json.NewDecoder(resp.Body).Decode(&torrentProps)
	if err != nil {
		panic(err)
	}

	return torrentProps
}
