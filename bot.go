package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/tucnak/telebot"
)

var (
	key         string
	gpiocmdPipe *io.PipeWriter
	bot         *telebot.Bot
)

func init() {
	keyBytes, err := ioutil.ReadFile("key")
	if err != nil {
		log.Fatalf("Cannot load bot token from key file: %s\n", err)
	}
	key = strings.TrimSpace(string(keyBytes))

	if bot, err = telebot.NewBot(key); err != nil {
		log.Fatalf("Cannot initialize bot: %s", err)
	}

	gpiocmd := exec.Command("python", "door.py")
	gpiocmd.Stdin, gpiocmdPipe = io.Pipe()
	if err := gpiocmd.Start(); err != nil {
		log.Fatalf("Can't run door.py: %s", err)
	}
}

func main() {
	doorManager := map[string]bool{
		"rsghost": true,
		"xatier":  true,
		"ronmi":   true,
	}

	messages := make(chan telebot.Message)
	bot.Listen(messages, 1*time.Second)
	fmt.Println("Waiting for commands")

	for message := range messages {
		fmt.Printf("[%s]: %s -> %s]\n",
			message.Chat.Title, message.Sender.Username, message.Text)

		if doorManager[message.Sender.Username] {
			switch message.Text {
			case "/ping":
				bot.SendMessage(message.Chat,
					"pong, "+message.Sender.FirstName+"!", nil)
			case "/up":
				fmt.Fprintln(gpiocmdPipe, "up")
				bot.SendMessage(message.Chat, "door up!", nil)

			case "/down":
				fmt.Fprintln(gpiocmdPipe, "down")
				bot.SendMessage(message.Chat, "door down!", nil)

			case "/stop":
				fmt.Fprintln(gpiocmdPipe, "stop")
				bot.SendMessage(message.Chat, "door stop!", nil)
			}
		}
	}
}
