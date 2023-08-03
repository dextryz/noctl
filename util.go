package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/ffiat/nostr"
)

func StringEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("address env variable \"%s\" not set, usual", key)
	}
	return value
}

func ShowEvent(e *nostr.Event, author Author) {

	t := e.CreatedAt.Time()
	date := t.Format("2006-01-02")
	clock := t.Format("15:04")

	fmt.Printf("\nPosted by [%s] on [%s] at [%s]\n\n", author.Name, date, clock)
	fmt.Printf(" \033[1;32m*\033[0m %s\n", e.Content)
}

func PrintJson(s any) {
	data, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
	}
	fmt.Printf("%s\n", data)
}
