package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"strings"
	"time"

	"github.com/tucnak/telebot"
)

func main() {

	doorManager := map[string]bool{
		"rsghost": true,
		"xatier":  true,
		"ronmi":   true,
	}

	fmt.Println("Read telegram bot secret key")
	key, err := ioutil.ReadFile("key")
	if err != nil {
		return
	}

	bot, err := telebot.NewBot(strings.TrimSpace(string(key)))
	if err != nil {
		return
	}
	fmt.Println("Connedted to telegram server")

	gpiocmdStdin, gpiocmdPipe := io.Pipe()
	gpiocmd := exec.Command("python", "door.py")
	gpiocmd.Stdin = gpiocmdStdin
	if err := gpiocmd.Start(); err != nil {
		fmt.Println("Can't run door.py")
		return
	}

	messages := make(chan telebot.Message)
	bot.Listen(messages, 1*time.Second)
	fmt.Println("Waiting for commands")

	for message := range messages {
		fmt.Println("[" + message.Chat.Title + "]: " +
			message.Sender.Username + " -> " +
			message.Text)
		if doorManager[message.Sender.Username] {
			switch message.Text {
			case "/ping":
				bot.SendMessage(message.Chat,
					"pong, "+message.Sender.FirstName+"!", nil)

			case "/up":
				fmt.Fprint(gpiocmdPipe, "up\n")
				bot.SendMessage(message.Chat, "door up!", nil)

			case "/down":
				fmt.Fprint(gpiocmdPipe, "down\n")
				bot.SendMessage(message.Chat, "door down!", nil)

			case "/stop":
				fmt.Fprint(gpiocmdPipe, "stop\n")
				bot.SendMessage(message.Chat, "door stop!", nil)
			}
		}
	}
}
