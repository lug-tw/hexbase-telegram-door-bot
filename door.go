package main

import (
	"fmt"
	"os"

	"github.com/Patrolavia/botgoram/telegram"
)

func sendDoorCommand(cmd string) (err error) {
	f, err := os.Open(doorctl)
	if err != nil {
		return
	}
	defer f.Close()

	_, err = fmt.Fprintln(f, cmd)
	return
}

func processCommand(cmd string, chat *telegram.Chat) {
	reply := "door " + cmd + "!"
	if err := sendDoorCommand(cmd); err != nil {
		reply = fmt.Sprintf("Error sending command %s: %s", cmd, err)
		fmt.Println(reply)
	}
	api.SendMessage(chat, reply, nil)
}

func processMessage(message *telegram.Message) (pass bool) {
	defer fmt.Printf("[%s]: %s -> %s]\n",
		message.Chat.Title, message.Sender.Username, message.Text)

	if doorManager[message.Sender.Username] {
		switch message.Text {
		case "/ping":
			api.SendMessage(message.Chat,
				"pong, "+message.Sender.FirstName+"!", nil)
		case "/up":
			processCommand("up", message.Chat)
		case "/down":
			processCommand("down", message.Chat)
		case "/stop":
			processCommand("stop", message.Chat)
		default:
			return true
		}
	}

	return
}
