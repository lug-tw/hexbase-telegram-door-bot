package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"testing"

	"github.com/Patrolavia/telegram"
)

func initFile(unames []string) (fn string, err error) {
	f, err := ioutil.TempFile("", "KHM")
	if err != nil {
		return
	}
	defer f.Close()
	fn = f.Name()

	data, err := json.Marshal(unames)
	if err != nil {
		return
	}
	_, err = f.Write(data)
	return
}

func initUser(uid int) *telegram.Victim {
	return &telegram.Victim{ID: int64(uid), Username: strconv.FormatInt(int64(uid), 10)}
}

func loadKHM(uids []int) (kh KeyHolderManager, fn string, err error) {
	unames := make([]string, len(uids))
	for idx, id := range uids {
		unames[idx] = strconv.FormatInt(int64(id), 10)
	}
	fn, err = initFile(unames)
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
		if !kh.Has(u) {
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
	if !kh.Has(u) {
		t.Errorf("User#%D should be a keyholder, but KHM return false", u.ID)
	}
	kh2, err := LoadKeyholders(fn)
	if err != nil {
		t.Fatalf("Failed to load keyholder again: %s", err)
	}
	if !kh2.Has(u) {
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
	if kh.Has(u) {
		t.Errorf("User#%D should not be a keyholder, but KHM return true", u.ID)
	}
	kh2, err := LoadKeyholders(fn)
	if err != nil {
		t.Fatalf("Failed to load keyholder again: %s", err)
	}
	if kh2.Has(u) {
		t.Errorf("User#%D should not be a keyholder, but second KHM return true", u.ID)
	}
}
