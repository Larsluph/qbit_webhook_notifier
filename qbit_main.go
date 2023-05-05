package main

import (
	"os"

	"github.com/joho/godotenv"
)

func main() {
	hash := os.Args[2]

	godotenv.Load()

	if os.Args[1] == "added" {
		webhook_added(hash)
	} else if os.Args[1] == "completed" {
		webhook_completed(hash)
	}
}
