package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/Patrolavia/botgoram/telegram"
)

func main() {
	var doorctl string
	doorManager := map[string]bool{
		"rsghost": true,
		"xatier":  true,
		"ronmi":   true,
	}

	// get named pipe to control door
	flag.StringVar(&doorctl, "s", "/tmp/doorctl", "path to unix socket for controlling door")
	flag.Parse()
	if _, err := os.Stat(doorctl); err != nil {
		log.Fatalf("doorctl %s does not exists!")
	}

	keyBytes, err := ioutil.ReadFile("key")
	if err != nil {
		log.Fatalf("Cannot load bot token from key file: %s\n", err)
	}
	key := strings.TrimSpace(string(keyBytes))

	api := telegram.New(key)
	if _, err = api.Me(); err != nil {
		log.Fatalf("Error validating bot token: %s", err)
	}


	processer := &CommandProcesser{
		Control: DoorControl(doorctl),
		Telegram: api,
		DoorManager: doorManager,
	}

	messages := make(chan *telegram.Message)
	go func(messages chan *telegram.Message) {
		offset := 0
		for {
			updates, err := api.GetUpdates(offset, 0, 30) // 30s timeout for long-polling
			if err != nil {
				fmt.Printf("Cannot fetch new messages: %s", err)
				continue
			}
			for _, update := range updates {
				offset++
				messages <- update.Message
			}
		}
	}(messages)

	fmt.Println("Waiting for commands")

	for message := range messages {
		processer.Handle(message)
	}
}
