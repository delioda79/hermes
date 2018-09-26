package logger

import (
	"fmt"

	"bitbucket.org/ConsentSystems/logging"
)

type basicLogger struct {
}

// NewBasicLogger returns a new basic logger
func NewBasicLogger() logging.Logger {
	return &basicLogger{}
}

func (bl *basicLogger) Info(msg logging.Log) error {
	fmt.Printf("%v", msg)
	return nil
}
func (bl *basicLogger) Debug(msg logging.Log) error {
	fmt.Printf("%v", msg)
	return nil
}
func (bl *basicLogger) Warning(msg logging.Log) error {
	fmt.Printf("%v", msg)
	return nil
}
func (bl *basicLogger) Error(msg logging.Log) error {
	fmt.Printf("%v", msg)
	return nil
}
func (bl *basicLogger) Fatal(msg logging.Log) error {
	fmt.Printf("%v", msg)
	return nil
}

func (bl *basicLogger) SetLevel(lvl int) {
}
