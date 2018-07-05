package pusher

import (
	"encoding/json"

	"nanomsg.org/go-mangos/protocol/pub"

	"bitbucket.org/ConsentSystems/mango-micro/messages"
	mangos "nanomsg.org/go-mangos"
)

type Pusher interface {
	mangos.Socket
	Push(name string, message []byte) error
}

type defaultPusher struct {
	mangos.Socket
}

func (pubs *defaultPusher) Push(name string, message []byte) error {
	trg := &messages.Trigger{
		Name:   name,
		Params: message,
	}

	bts, err := json.Marshal(trg)
	if err != nil {
		return err
	}
	pubs.Send(bts)
	return nil
}

func NewPusher() (Pusher, error) {
	sock, err := pub.NewSocket()
	if err != nil {
		return nil, err
	}
	return &defaultPusher{
		Socket: sock,
	}, nil
}
