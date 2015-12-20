package main

import (
	"log"
	"strings"

	"github.com/Patrolavia/botgoram"
	"github.com/Patrolavia/botgoram/telegram"
	"github.com/dgryski/dgoogauth"
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
			return nil
		}

		a.Admins.Add(current.User().ToUser())
		api.SendMessage(current.User(), "You are administrator now.", nil)

		current.Transit("")
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
		botgoram.TransitorMap{
			Transitor: func(msg *telegram.Message, state botgoram.State) (next string, err error) {
				next = ""
				return
			},
			IsHidden: true,
			State:    "",
			Type:     telegram.CONTACT,
			Desc:     `Done validating`,
		},
	}
}

func registerAuthStates(fsm botgoram.FSM, otpcfg *dgoogauth.OTPConfig, admins KeyHolderManager) {
	if _, err := fsm.MakeState(AuthAskPass("/auth")); err != nil {
		log.Fatalf("Error registering state askpass: %s", err)
	}
	if _, err := fsm.MakeState(&AuthValidate{
		StateName: "auth:validate",
		Config:    otpcfg,
		Admins:    admins,
	}); err != nil {
		log.Fatalf("Error registering state askpass: %s", err)
	}
}
