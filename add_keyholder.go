package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/Patrolavia/botgoram"
	"github.com/Patrolavia/botgoram/telegram"
)

type AddAskContact struct {
	Admins  KeyHolderManager
	Command string
}

func (a *AddAskContact) Name() string {
	return "addkh:askcontact"
}

func (a *AddAskContact) Actions() (enter botgoram.Action, leave botgoram.Action) {
	enter = func(msg *telegram.Message, current botgoram.State, api telegram.API) error {
		u := msg.Sender
		if !a.Admins.Has(u) {
			return fmt.Errorf("%s(%s-%d) is not admin", u.FirstName, u.Username, u.ID)
		}
		api.SendMessage(u, "Please send me a username (without @)", nil)
		return nil
	}
	return
}

func (a *AddAskContact) Transitors() []botgoram.TransitorMap {
	return []botgoram.TransitorMap{
		botgoram.TransitorMap{
			Transitor: func(msg *telegram.Message, state botgoram.State) (next string, err error) {
				next = "addkh:askcontact"
				return
			},
			State:   "",
			Type:    telegram.TEXT,
			Command: a.Command,
			Desc:    `Accepted /add command, prompt for contact info.`,
		},
	}
}

type AddValidate struct {
	KeyHolders KeyHolderManager
}

func (a *AddValidate) Name() string {
	return "addkh:validate"
}

func (a *AddValidate) Actions() (enter botgoram.Action, leave botgoram.Action) {
	enter = func(msg *telegram.Message, current botgoram.State, api telegram.API) error {
		contact := strings.TrimSpace(msg.Text)

		a.KeyHolders.Add(&telegram.User{Username: contact})
		api.SendMessage(current.User(), "Key holder added.", nil)
		current.Transit("")
		return nil
	}
	return
}

func (a *AddValidate) Transitors() []botgoram.TransitorMap {
	return []botgoram.TransitorMap{
		botgoram.TransitorMap{
			Transitor: func(msg *telegram.Message, state botgoram.State) (next string, err error) {
				next = "addkh:validate"
				return
			},
			State: "addkh:askcontact",
			Type:  telegram.TEXT,
			Desc:  `Validate the contact info`,
		},
		botgoram.TransitorMap{
			Transitor: func(msg *telegram.Message, state botgoram.State) (next string, err error) {
				next = ""
				return
			},
			IsHidden: true,
			State:    "",
			Type:     telegram.TEXT,
			Desc:     `Done validating`,
		},
	}
}

func registerAddKHStates(fsm botgoram.FSM, admins KeyHolderManager, kh KeyHolderManager) {
	_, err := fsm.MakeState(&AddAskContact{
		Admins:  admins,
		Command: "/add",
	})
	if err != nil {
		log.Fatalf("Error registering state addkh:askcontact: %s", err)
	}

	_, err = fsm.MakeState(&AddValidate{
		KeyHolders: kh,
	})
	if err != nil {
		log.Fatalf("Error registering state addkh:validate: %s", err)
	}
}
