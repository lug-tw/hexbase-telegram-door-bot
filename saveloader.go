package main

import "sync"

type NoDataSaveLoader struct {
	*sync.Mutex
	state map[string]string
}

func NewSL() *NoDataSaveLoader {
	return &NoDataSaveLoader{
		&sync.Mutex{},
		make(map[string]string),
	}
}

func (s *NoDataSaveLoader) Save(uid string, sid string, data interface{}) error {
	s.Lock()
	defer s.Unlock()

	s.state[uid] = sid
	if sid == "" {
		delete(s.state, uid)
	}
	return nil
}

func (s *NoDataSaveLoader) Load(uid string) (sid string, data interface{}, err error) {
	return s.state[uid], nil, nil
}
