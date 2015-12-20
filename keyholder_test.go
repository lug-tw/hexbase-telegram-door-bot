package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/Patrolavia/botgoram/telegram"
)

func initFile(uids []int) (fn string, err error) {
	f, err := ioutil.TempFile("", "KHM")
	if err != nil {
		return
	}
	defer f.Close()
	fn = f.Name()

	data, err := json.Marshal(uids)
	if err != nil {
		return
	}
	_, err = f.Write(data)
	return
}

func initUser(uid int) *telegram.User {
	return &telegram.User{ID: uid}
}

func loadKHM(uids []int) (kh KeyHolderManager, fn string, err error) {
	fn, err = initFile(uids)
	if err != nil {
		return
	}
	kh, err = LoadKeyholders(fn)
	return
}

func TestLoadKHM(t *testing.T) {
	uids := []int{1, 2, 3}
	kh, fn, err := loadKHM(uids)
	if err != nil {
		t.Fatalf("Cannot load keyholder from file: ", err)
	}
	defer os.Remove(fn)

	for _, uid := range uids {
		u := initUser(uid)
		if !kh.Is(u) {
			t.Errorf("User#%D should be a keyholder, but KHM return false", uid)
		}
	}

}

func TestAddKH(t *testing.T) {
	kh, fn, err := loadKHM([]int{1, 2, 3})
	if err != nil {
		t.Fatalf("Cannot load keyholder from file: ", err)
	}
	defer os.Remove(fn)

	u := initUser(4)
	kh.Add(u)
	if !kh.Is(u) {
		t.Errorf("User#%D should be a keyholder, but KHM return false", u.ID)
	}
	kh2, err := LoadKeyholders(fn)
	if err != nil {
		t.Fatalf("Failed to load keyholder again: %s", err)
	}
	if !kh2.Is(u) {
		t.Errorf("User#%D should be a keyholder, but second KHM return false", u.ID)
	}
}

func TestRemoveKH(t *testing.T) {
	kh, fn, err := loadKHM([]int{1, 2, 3})
	if err != nil {
		t.Fatalf("Cannot load keyholder from file: ", err)
	}
	defer os.Remove(fn)

	u := initUser(3)
	kh.Remove(u)
	if kh.Is(u) {
		t.Errorf("User#%D should not be a keyholder, but KHM return true", u.ID)
	}
	kh2, err := LoadKeyholders(fn)
	if err != nil {
		t.Fatalf("Failed to load keyholder again: %s", err)
	}
	if kh2.Is(u) {
		t.Errorf("User#%D should not be a keyholder, but second KHM return true", u.ID)
	}
}
