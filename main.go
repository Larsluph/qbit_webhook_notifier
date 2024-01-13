package main

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"golang.org/x/exp/slices"
	"os"
	"qbit_webhook/discord"
	"qbit_webhook/helpers"
	"strings"
)

func main() {
	// Recover from panics and send them to Discord
	defer func() {
		if r := recover(); r != nil {
			// Check if r is already an error type
			if err, ok := r.(error); ok {
				// If it's an error type, directly use it in GenerateErrorEmbed
				discord.SendWebhook(discord.GenerateErrorEmbed(*helpers.NewErrorPayload(err)))
				fmt.Println(err)
			} else {
				// If r is not an error type, wrap it into an error using fmt.Errorf
				err := fmt.Errorf("recovered panic: %v", r)
				discord.SendWebhook(discord.GenerateErrorEmbed(*helpers.NewErrorPayload(err)))
				fmt.Println(err)
			}
		}
	}()

	if len(os.Args) != 4 {
		panic(errors.New(fmt.Sprintf("Invalid syntax\nSyntax: ./qbit_webhook added|completed ENV_PATH TORRENT_HASH\n   was: %s", strings.Join(os.Args, " "))))
	}

	triggerType := os.Args[1]
	config := os.Args[2]
	hash := os.Args[3]

	var availableTriggers = []string{"added", "completed"}
	if !slices.Contains(availableTriggers, triggerType) {
		panic(errors.New("unknown trigger"))
	}

	err := godotenv.Load(config)
	if err != nil {
		payload := helpers.NewErrorPayload(err)
		fmt.Println(payload)
		return
	}

	discord.TriggerWebhook(triggerType, hash)
}
