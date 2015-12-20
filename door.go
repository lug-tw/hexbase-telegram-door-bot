package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Patrolavia/botgoram/telegram"
)

type DoorControl string

func (d DoorControl) Send(cmd string) (err error) {
	f, err := os.Create(string(d))
	if err != nil {
		return
	}
	defer f.Close()

	_, err = fmt.Fprintln(f, cmd)
	return
}

type CommandProcesser struct {
	Control  DoorControl
	Telegram telegram.API
	Admins   KeyHolderManager
	Members  KeyHolderManager
}

func (c *CommandProcesser) chatCommand(cmd string, chat *telegram.Chat) {
	reply := "door " + cmd + "!"
	if err := c.Control.Send(cmd); err != nil {
		reply = fmt.Sprintf("Error sending command %s: %s", cmd, err)
		fmt.Println(reply)
	}
	c.Telegram.SendMessage(chat, reply, nil)
}

func (c *CommandProcesser) Handle(message *telegram.Message) (pass bool) {
	defer fmt.Printf("[%s]: %s -> %s]\n",
		message.Chat.Title, message.Sender.Username, message.Text)

	if !c.Admins.Has(message.Sender) && !c.Members.Has(message.Sender) {
		return true
	}
	switch message.Text {
	case "/ping":
		c.Telegram.SendMessage(message.Chat,
			"pong, "+message.Sender.FirstName+"!", nil)
	case "/up":
		c.chatCommand("up", message.Chat)
	case "/down":
		c.chatCommand("down", message.Chat)
	case "/stop":
		c.chatCommand("stop", message.Chat)
	case "/list":
		c.listHolders(message.Sender)
	default:
		return true
	}

	return false
}

func (c *CommandProcesser) listHolders(u *telegram.User) {
	reply := fmt.Sprintf(
		`Administrators:
%s

Key holders:
%s`,
		"@"+strings.Join(c.Admins.List(), ", @"),
		"@"+strings.Join(c.Members.List(), ", @"),
	)
	c.Telegram.SendMessage(u, reply, nil)
}
