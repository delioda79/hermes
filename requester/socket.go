package requester

import (
	"encoding/json"

	"bitbucket.org/ddanna79/mango-micro/messages"
	mangos "nanomsg.org/go-mangos"
)

type Requester interface {
	mangos.Socket
	Request(name string, message []byte) ([]byte, error)
}

type defaultRequester struct {
	mangos.Socket
	ch       chan mangos.Message
	uid      string
	suppSock mangos.Socket
}

func (rqs *defaultRequester) Request(name string, message []byte) ([]byte, error) {
	trg := &messages.Trigger{
		Name:   name,
		Params: message,
		UID:    rqs.uid,
	}

	bts, err := json.Marshal(trg)
	if err != nil {
		return nil, err
	}

	msg := mangos.NewMessage(len(bts))
	msg.Body = bts
	encoder := rqs.suppSock.GetProtocol().(mangos.ProtocolSendHook)
	encoder.SendHook(msg)
	err = rqs.Socket.SendMsg(msg)
	if err != nil {
		return nil, err
	}

	resp := <-rqs.ch
	trigger := &messages.Trigger{}
	err = json.Unmarshal(resp.Body, trigger)
	if err != nil {
		return nil, err
	}
	return trigger.Params, err
}

func NewRequester(uid string, sock, supSock mangos.Socket) (Requester, chan mangos.Message) {
	ch := make(chan mangos.Message, 1)
	return &defaultRequester{
		Socket:   sock,
		ch:       ch,
		uid:      uid,
		suppSock: supSock,
	}, ch
}
