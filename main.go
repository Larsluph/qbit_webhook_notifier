package main

import (
	"github.com/joho/godotenv"
	"golang.org/x/exp/slices"
	"os"
	discord2 "qbit_webhook/discord"
)

func main() {
	triggerType := os.Args[1]
	config := os.Args[2]
	hash := os.Args[3]

	var availableTriggers = []string{"added", "completed"}
	if !slices.Contains(availableTriggers, triggerType) {
		panic("Unknown trigger")
	}

	err := godotenv.Load(config)
	if err != nil {
		//TODO
		return
	}

	discord2.TriggerWebhook(triggerType, hash)
}
