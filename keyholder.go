package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"

	"github.com/Patrolavia/telegram"
)

type KeyHolderManager interface {
	Add(user *telegram.Victim) (err error)
	Remove(user *telegram.Victim) (err error)
	Has(user *telegram.Victim) bool
	List() []string
}

type keyholder struct {
	filename string
	users    map[string]bool
	*sync.Mutex
}

func export(filename string, users map[string]bool) (err error) {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	validKeyholder := make([]string, 0, len(users))
	for user, valid := range users {
		if valid {
			validKeyholder = append(validKeyholder, user)
		}
	}

	data, err := json.Marshal(validKeyholder)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	return
}

func (a *keyholder) List() []string {
	ret := make([]string, 0, len(a.users))
	for u, valid := range a.users {
		if valid {
			ret = append(ret, u)
		}
	}
	return ret
}

func (a *keyholder) Add(user *telegram.Victim) (err error) {
	if user.Username == "" {
		return
	}
	a.Lock()
	defer a.Unlock()

	a.users[user.Username] = true
	return export(a.filename, a.users)
}

func (a *keyholder) Remove(user *telegram.Victim) (err error) {
	if user.Username == "" {
		return
	}
	a.Lock()
	defer a.Unlock()

	delete(a.users, user.Username)
	return export(a.filename, a.users)
}

func (a *keyholder) Has(user *telegram.Victim) bool {
	if user.Username == "" {
		return false
	}
	return a.users[user.Username]
}

func LoadKeyholders(filename string) (keyholders KeyHolderManager, err error) {
	users := make(map[string]bool)
	keyholders = &keyholder{filename, users, &sync.Mutex{}}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return keyholders, nil
	}

	items := []string{}
	if err = json.Unmarshal(data, &items); err != nil {
		return
	}

	for _, item := range items {
		keyholders.(*keyholder).users[item] = true
	}

	return
}
