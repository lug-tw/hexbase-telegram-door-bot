package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/Patrolavia/telegram"
)

type DoorControl string

func (d DoorControl) Send(cmd string) (err error) {
	conn, err := net.Dial("unix", string(d))
	if err != nil {
		return
	}
	defer conn.Close()

	_, err = fmt.Fprintln(conn, cmd)
	return
}

type CommandProcesser struct {
	Control  DoorControl
	Telegram telegram.API
	Admins   KeyHolderManager
	Members  KeyHolderManager
}

func (c *CommandProcesser) chatCommand(cmd string, chat *telegram.Victim) {
	reply := "door " + cmd + "!"
	if err := c.Control.Send(cmd); err != nil {
		reply = fmt.Sprintf("Error sending command %s: %s", cmd, err)
		fmt.Println(reply)
	}
	c.Telegram.SendMessage(chat.Identifier(), reply, nil)
}

func (c *CommandProcesser) Handle(message *telegram.Message) (pass bool) {
	defer fmt.Printf("[%s]: %s -> %s]\n",
		message.Chat.Title, message.From.Username, message.Text)

	if !c.Admins.Has(message.From) && !c.Members.Has(message.From) {
		return true
	}
	switch message.Text {
	case "/ping":
		c.Telegram.SendMessage(message.Chat.Identifier(),
			"pong, "+message.From.FirstName+"!", nil)
	case "/up":
		c.chatCommand("up", message.Chat)
	case "/down":
		c.chatCommand("down", message.Chat)
	case "/stop":
		c.chatCommand("stop", message.Chat)
	case "/list":
		c.listHolders(message.From)
	default:
		return true
	}

	return false
}

func (c *CommandProcesser) listHolders(u *telegram.Victim) {
	mkstr := func(kl KeyHolderManager) (ret string) {
		arr := kl.List()
		if len(arr) > 0 {
			ret = "@" + strings.Join(arr, ", @")
		}
		return
	}

	reply := fmt.Sprintf(
		`Administrators:
%s

Key holders:
%s`,
		mkstr(c.Admins),
		mkstr(c.Members),
	)
	c.Telegram.SendMessage(u.Identifier(), reply, nil)
}
