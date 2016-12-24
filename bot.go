package main

import (
	"encoding/base32"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"log/syslog"
	"net/url"
	"os"
	"strings"

	"github.com/Patrolavia/botgoram"
	"github.com/Patrolavia/telegram"
	"github.com/dgryski/dgoogauth"
)

var logger *log.Logger

func init() {
	var err error
	// log to syslog
	logger, err = syslog.NewLogger(syslog.LOG_DAEMON|syslog.LOG_ERR, log.LstdFlags)

	if err != nil {
		log.Fatalf("cannot log to syslog: %s", err)
	}
}

func main() {
	var (
		doorctl   string
		tokenfile string
		adminfile string
		khfile    string
		secretkey string
		dumpmap   string
	)
	flag.StringVar(&doorctl, "s", "/tmp/doorctl", "path to unix socket for controlling door")
	flag.StringVar(&tokenfile, "t", "token", "file contains telegram bot token")
	flag.StringVar(&adminfile, "a", "admins", "file stores administrator lists")
	flag.StringVar(&khfile, "h", "keyholders", "file stores keyholder lists")
	flag.StringVar(&secretkey, "k", "", "10bytes secret key in hexdecimal")
	flag.StringVar(&dumpmap, "m", "", "dump state map to this file")
	flag.Parse()
	// get named pipe to control door
	if _, err := os.Stat(doorctl); err != nil {
		logger.Fatalf("doorctl %s does not exists!", doorctl)
	}

	keyBytes, err := ioutil.ReadFile(tokenfile)
	if err != nil {
		logger.Fatalf("Cannot load bot token from key file: %s\n", err)
	}
	key := strings.TrimSpace(string(keyBytes))

	api := telegram.New(key, nil)
	if _, err = api.GetMe(); err != nil {
		logger.Fatalf("Error validating bot token: %s", err)
	}

	admins, err := LoadKeyholders(adminfile)
	if err != nil {
		logger.Fatalf("Cannot load admins from %s: %s", adminfile, err)
	}
	logger.Printf("%#v", admins)
	khs, err := LoadKeyholders(khfile)
	if err != nil {
		logger.Fatalf("Cannot load keyholders from %s: %s", khfile, err)
	}

	secretBytes, err := hex.DecodeString(secretkey)
	if err != nil || len(secretBytes) != 10 {
		logger.Fatalf("OTP secret is not 10bytes hexdecimal string!")
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

	cp := &CommandProcesser{
		Control:  DoorControl(doorctl),
		Telegram: api,
		Admins:   admins,
		Members:  khs,
		In:       make(chan *telegram.Message),
		Out:      make(chan *telegram.Message),
	}
	go cp.Process()

	// long-polling
	lp := &telegram.LongPollFetcher{
		Message: cp.In,
		API:     api,
	}
	go lp.Fetch(1, 30)

	fsm := botgoram.NewBySender(
		api,
		NewSL(),
		1,
		cp.Out,
	)

	registerAuthStates(fsm, otpcfg, admins)
	registerAddKHStates(fsm, admins, khs)
	registerDelKHStates(fsm, admins, khs)

	// register fallback
	initState, _ := fsm.State("")
	initState.RegisterFallback(func(msg *telegram.Message, state botgoram.State) (next string, err error) {
		// ignore invalid command
		return
	})

	if dumpmap != "" {
		dot := fsm.StateMap("doorbot")
		ioutil.WriteFile(dumpmap, []byte(dot), 0644)
	}

	fmt.Println("Waiting messages.")
	err = fsm.Start(30)
	for err != nil {
		logger.Printf("Error happened: %s", err)
		err = fsm.Resume()
	}
}
