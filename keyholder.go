package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"

	"github.com/Patrolavia/botgoram/telegram"
)

type KeyHolderManager interface {
	Add(user *telegram.User) (err error)
	Remove(user *telegram.User) (err error)
	Has(user *telegram.User) bool
}

type keyholder struct {
	filename string
	users    map[int]bool
	*sync.Mutex
}

func export(filename string, users map[int]bool) (err error) {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	validKeyholder := make([]int, 0, len(users))
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

func (a *keyholder) Add(user *telegram.User) (err error) {
	a.Lock()
	defer a.Unlock()

	a.users[user.ID] = true
	return export(a.filename, a.users)
}

func (a *keyholder) Remove(user *telegram.User) (err error) {
	a.Lock()
	defer a.Unlock()

	delete(a.users, user.ID)
	return export(a.filename, a.users)
}

func (a *keyholder) Has(user *telegram.User) bool {
	return a.users[user.ID]
}

func LoadKeyholders(filename string) (keyholders KeyHolderManager, err error) {
	users := make(map[int]bool)
	keyholders = &keyholder{filename, users, &sync.Mutex{}}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return keyholders, nil
	}

	items := []int{}
	if err = json.Unmarshal(data, &items); err != nil {
		return
	}

	for _, item := range items {
		keyholders.(*keyholder).users[item] = true
	}

	return
}
