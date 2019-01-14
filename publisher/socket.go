package publisher

import (
	"encoding/json"

	"nanomsg.org/go-mangos/protocol/pub"

	"bitbucket.org/ddanna79/mango-micro/messages"
	mangos "nanomsg.org/go-mangos"
)

type Publisher interface {
	mangos.Socket
	Publish(name string, message []byte) error
}

type defaultPubisher struct {
	mangos.Socket
}

func (pubs *defaultPubisher) Publish(name string, message []byte) error {
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

func NewPublisher() (Publisher, error) {
	sock, err := pub.NewSocket()
	if err != nil {
		return nil, err
	}
	return &defaultPubisher{
		Socket: sock,
	}, nil
}
