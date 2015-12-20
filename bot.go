package main

import (
	"encoding/base32"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/Patrolavia/botgoram"
	"github.com/Patrolavia/botgoram/telegram"
	"github.com/dgryski/dgoogauth"
)

type MockTelegram struct {
	telegram.API
	Processer *CommandProcesser
	max       int
}

func (m *MockTelegram) GetUpdates(offset, limit, timeout int) (updates []telegram.Update, err error) {
	if offset < m.max {
		offset = m.max
	}
	if updates, err = m.API.GetUpdates(offset, limit, timeout); err != nil {
		return
	}

	ret := make([]telegram.Update, 0, len(updates))
	for _, update := range updates {
		if m.max <= update.ID {
			m.max = update.ID + 1
		}
		if m.Processer.Handle(update.Message) {
			ret = append(ret, update)
		}
	}
	return ret, err
}

func main() {
	var (
		doorctl   string
		tokenfile string
		adminfile string
		khfile    string
		secretkey string
	)
	flag.StringVar(&doorctl, "s", "/tmp/doorctl", "path to unix socket for controlling door")
	flag.StringVar(&tokenfile, "t", "token", "file contains telegram bot token")
	flag.StringVar(&adminfile, "a", "admins", "file stores administrator lists")
	flag.StringVar(&khfile, "h", "keyholders", "file stores keyholder lists")
	flag.StringVar(&secretkey, "k", "", "10bytes secret key in hexdecimal")
	flag.Parse()
	// get named pipe to control door
	if _, err := os.Stat(doorctl); err != nil {
		log.Fatalf("doorctl %s does not exists!", doorctl)
	}

	keyBytes, err := ioutil.ReadFile(tokenfile)
	if err != nil {
		log.Fatalf("Cannot load bot token from key file: %s\n", err)
	}
	key := strings.TrimSpace(string(keyBytes))

	api := telegram.New(key)
	if _, err = api.Me(); err != nil {
		log.Fatalf("Error validating bot token: %s", err)
	}

	admins, err := LoadKeyholders(adminfile)
	if err != nil {
		log.Fatalf("Cannot load admins from %s: %s", adminfile, err)
	}
	khs, err := LoadKeyholders(khfile)
	if err != nil {
		log.Fatalf("Cannot load keyholders from %s: %s", khfile, err)
	}

	secretBytes, err := hex.DecodeString(secretkey)
	if err != nil || len(secretBytes) != 10 {
		log.Fatalf("OTP secret is not 10bytes hexdecimal string!")
	}
	otpcfg := &dgoogauth.OTPConfig{
		Secret:     base32.StdEncoding.EncodeToString(secretBytes),
		WindowSize: 5,
	}
	otpuri := otpcfg.ProvisionURIWithIssuer("DoorControl", "Hexbase")
	fmt.Println("OTP uri:", otpuri)
	fmt.Println(
		"QRCode url: https://chart.googleapis.com/chart?cht=qr&chs=256x256&chl=" +
			url.QueryEscape(otpuri),
	)

	mock := &MockTelegram{
		API: api,
		Processer: &CommandProcesser{
			Control:  DoorControl(doorctl),
			Telegram: api,
			Admins:   admins,
			Members:  khs,
		},
	}

	fsm := botgoram.NewBySender(
		mock,
		botgoram.MemoryStore(func(uid string) interface{} {
			return true
		}),
		1,
	)

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

	// register fallback
	initState, _ := fsm.State("")
	initState.RegisterFallback(func(msg *telegram.Message, state botgoram.State) (next string, err error) {
		// ignore invalid command
		return
	})

	fmt.Println("Waiting messages.")
	err = fsm.Start(30)
	for err != nil {
		log.Printf("Error happened: %s", err)
		err = fsm.Resume()
	}
}
