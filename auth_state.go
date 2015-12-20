package main

import (
	"strings"

	"github.com/dgryski/dgoogauth"
	"github.com/Patrolavia/botgoram"
	"github.com/Patrolavia/botgoram/telegram"
)

type AuthAskPass string

func (a AuthAskPass) Name() string {
	return "auth:askpass"
}

func (a AuthAskPass) Actions() (enter botgoram.Action, leave botgoram.Action) {
	enter = func(msg *telegram.Message, current botgoram.State, api telegram.API) error {
		api.SendMessage(current.User(), "Please input your password (6 digits)", nil)
		return nil
	}
	return
}

func (a AuthAskPass) Transitors() []botgoram.TransitorMap {
	return []botgoram.TransitorMap{
		botgoram.TransitorMap{
			Transitor: func(msg *telegram.Message, state botgoram.State) (next string, err error) {
				next = "auth:askpass"
				return
			},
			State:   "",
			Type:    telegram.TEXT,
			Command: string(a),
			Desc:    `Accepted /auth command, prompt for password.`,
		},
	}
}

type AuthValidate struct {
	StateName string
	Config    *dgoogauth.OTPConfig
	Admins    KeyHolderManager
}

func (a *AuthValidate) Name() string {
	return a.StateName
}

func (a *AuthValidate) Actions() (enter botgoram.Action, leave botgoram.Action) {
	enter = func(msg *telegram.Message, current botgoram.State, api telegram.API) error {
		ok, err := a.Config.Authenticate(strings.TrimSpace(msg.Text))
		if !ok || err != nil {
			api.SendMessage(current.User(), "Password incorrect.", nil)
			current.Transit("")
			return nil
		}

		a.Admins.Add(current.User().ToUser())
		api.SendMessage(current.User(), "You are administrator now.", nil)

		return nil
	}
	return
}

func (a *AuthValidate) Transitors() []botgoram.TransitorMap {
	return []botgoram.TransitorMap{
		botgoram.TransitorMap{
			Transitor: func(msg *telegram.Message, state botgoram.State) (next string, err error) {
				next = a.StateName
				return
			},
			State: "auth:askpass",
			Type:  telegram.TEXT,
			Desc:  `Validate the password`,
		},
	}
}