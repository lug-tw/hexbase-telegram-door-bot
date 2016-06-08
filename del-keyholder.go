package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/Patrolavia/botgoram"
	"github.com/Patrolavia/telegram"
)

type DelAskContact struct {
	Admins     KeyHolderManager
	KeyHolders KeyHolderManager
	Command    string
}

func (a *DelAskContact) Name() string {
	return "delkh:askcontact"
}

func (a *DelAskContact) Actions() (enter botgoram.Action, leave botgoram.Action) {
	enter = func(msg *telegram.Message, current botgoram.State, api telegram.API) error {
		u := msg.From
		if !a.Admins.Has(u) {
			return fmt.Errorf("%s(%s-%d) is not admin", u.FirstName, u.Username, u.ID)
		}

		// make up custom keyboard
		khs := a.KeyHolders.List()
		khb := make([]telegram.KeyboardButton, 0, len(khs))
		for _, u := range khs {
			khb = append(khb, telegram.KeyboardButton{Text: u})
		}

		api.SendMessage(u.Identifier(), "Please send me a username (without @)", &telegram.Options{
			ReplyMarkup: &telegram.ReplyKeyboardMarkup{
				Keyboard: telegram.FormatKeyboard(khb, 2),
				Resize:   true,
				Once:     true,
			},
		})
		return nil
	}
	return
}

func (a *DelAskContact) Transitors() []botgoram.TransitorMap {
	return []botgoram.TransitorMap{
		botgoram.TransitorMap{
			Transitor: func(msg *telegram.Message, state botgoram.State) (next string, err error) {
				next = "delkh:askcontact"
				return
			},
			State:   "",
			Type:    botgoram.TextMsg,
			Command: a.Command,
			Desc:    `Accepted /del command, prompt for contact info.`,
		},
	}
}

type DelValidate struct {
	KeyHolders KeyHolderManager
}

func (a *DelValidate) Name() string {
	return "delkh:validate"
}

func (a *DelValidate) Actions() (enter botgoram.Action, leave botgoram.Action) {
	enter = func(msg *telegram.Message, current botgoram.State, api telegram.API) error {
		contact := strings.TrimSpace(msg.Text)

		a.KeyHolders.Remove(&telegram.Victim{Username: contact})
		api.SendMessage(current.User().Identifier(), "Key holder deleted.", &telegram.Options{
			ReplyMarkup: &telegram.ReplyKeyboardHide{
				Hide: true,
			},
		})
		current.Transit("")
		return nil
	}
	return
}

func (a *DelValidate) Transitors() []botgoram.TransitorMap {
	return []botgoram.TransitorMap{
		botgoram.TransitorMap{
			Transitor: func(msg *telegram.Message, state botgoram.State) (next string, err error) {
				next = "delkh:validate"
				return
			},
			State: "delkh:askcontact",
			Type:  botgoram.TextMsg,
			Desc:  `Validate the contact info`,
		},
		botgoram.TransitorMap{
			Transitor: func(msg *telegram.Message, state botgoram.State) (next string, err error) {
				next = ""
				return
			},
			IsHidden: true,
			State:    "",
			Type:     botgoram.TextMsg,
			Desc:     `Done validating`,
		},
	}
}

func registerDelKHStates(fsm botgoram.FSM, admins KeyHolderManager, kh KeyHolderManager) {
	_, err := fsm.MakeState(&DelAskContact{
		Admins:     admins,
		KeyHolders: kh,
		Command:    "/del",
	})
	if err != nil {
		log.Fatalf("Error registering state delkh:askcontact: %s", err)
	}

	_, err = fsm.MakeState(&DelValidate{
		KeyHolders: kh,
	})
	if err != nil {
		log.Fatalf("Error registering state delkh:validate: %s", err)
	}
}
