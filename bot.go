package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"

	"github.com/Patrolavia/botgoram/telegram"
)

var (
	key         string
	gpiocmdPipe *io.PipeWriter
	api         telegram.API
	bot         *telegram.User
	doorManager = map[string]bool{
		"rsghost": true,
		"xatier":  true,
		"ronmi":   true,
	}
)

func init() {
	keyBytes, err := ioutil.ReadFile("key")
	if err != nil {
		log.Fatalf("Cannot load bot token from key file: %s\n", err)
	}
	key = strings.TrimSpace(string(keyBytes))

	api = telegram.New(key)
	if bot, err = api.Me(); err != nil {
		log.Fatalf("Error validating bot token: %s", err)
	}

	// XXX: change this to a Unix domain socekt conneting to /tmp/doorctl
	gpiocmd := exec.Command("python", "door.py")
	gpiocmd.Stdin, gpiocmdPipe = io.Pipe()
	if err := gpiocmd.Start(); err != nil {
		log.Fatalf("Can't run door.py: %s", err)
	}
}

func processMessage(message *telegram.Message) (passed bool) {
	defer fmt.Printf("[%s]: %s -> %s]\n",
		message.Chat.Title, message.Sender.Username, message.Text)

	if doorManager[message.Sender.Username] {
		switch message.Text {
		case "/ping":
			api.SendMessage(message.Chat,
				"pong, "+message.Sender.FirstName+"!", nil)
		case "/up":
			fmt.Fprintln(gpiocmdPipe, "up")
			api.SendMessage(message.Chat, "door up!", nil)

		case "/down":
			fmt.Fprintln(gpiocmdPipe, "down")
			api.SendMessage(message.Chat, "door down!", nil)

		case "/stop":
			fmt.Fprintln(gpiocmdPipe, "stop")
			api.SendMessage(message.Chat, "door stop!", nil)
		default:
			return true
		}
	}

	return
}

func main() {
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
		processMessage(message)
	}
}
