package handler

import (
	"fmt"
)

// HandleFunc represents the action for a psecific event
type HandleFunc func(msg interface{}, rsp ...*[]byte) error

// Handler implements a base subscriber
type Handler interface {
	Add(name string, action HandleFunc) error
	Run(name string, param interface{}, rsp ...*[]byte) error
}

type defaultHandler struct {
	actions map[string]HandleFunc
}

func (ds *defaultHandler) Add(name string, action HandleFunc) error {
	if _, ok := ds.actions[name]; ok {
		return fmt.Errorf("Action already set for subscriber %s", name)
	}
	fmt.Println("We have ", name, action)
	ds.actions[name] = action
	return nil
}

func (ds *defaultHandler) Run(name string, param interface{}, rsp ...*[]byte) error {
	if action, ok := ds.actions[name]; ok {
		hdlf := action(param, rsp[0])
		return hdlf
	}

	return fmt.Errorf("No action set for %s", name)
}

// NewHandler returns a new default subscriber
func NewHandler() Handler {
	return &defaultHandler{
		actions: map[string]HandleFunc{},
	}
}
