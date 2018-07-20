package requester

import (
	"encoding/json"
	"errors"
	"sync"
	"time"

	"bitbucket.org/ConsentSystems/mango-micro/messages"
	mangos "nanomsg.org/go-mangos"
)

type Requester interface {
	mangos.Socket
	Request(name string, message []byte) ([]byte, error)
	SetDeadline(deadline time.Duration)
}

type defaultRequester struct {
	mangos.Socket
	ch       chan mangos.Message
	uid      string
	suppSock mangos.Socket
	srvCh    chan string
	deadline time.Duration
	ddlnChan chan error
	mx       *sync.Mutex
}

func (rqs *defaultRequester) Request(name string, message []byte) ([]byte, error) {
	//rqs.mx.Lock()
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
	rqs.mx.Lock()
	err = rqs.Socket.SendMsg(msg)
	rqs.mx.Unlock()

	if rqs.deadline > 0 {
		go func(deadline time.Duration) {
			time.Sleep(deadline)
			rqs.mx.Lock()
			// Terminating the request from server
			rqs.srvCh <- rqs.uid
			rqs.ddlnChan <- errors.New("Deadline eceeded")
			rqs.mx.Unlock()
		}(rqs.deadline)
	}
	if err != nil {
		return nil, err
	}

	select {
	case resp := <-rqs.ch:
		rqs.mx.Lock()
		// Terminating the request from server
		rqs.srvCh <- rqs.uid
		// Returning the message
		trigger := &messages.Trigger{}
		err = json.Unmarshal(resp.Body, trigger)
		rqs.mx.Unlock()
		if err != nil {
			return nil, err
		}
		return trigger.Params, err
	case err := <-rqs.ddlnChan:
		//rqs.mx.Unlock()
		return nil, err
	}

}

func (rqs *defaultRequester) SetDeadline(deadline time.Duration) {
	rqs.deadline = deadline
}

func NewRequester(
	uid string,
	sock,
	supSock mangos.Socket,
	srvChan chan string,
	mx *sync.Mutex,
) (Requester, chan mangos.Message) {
	ch := make(chan mangos.Message, 1)
	ddlnChan := make(chan error, 1)
	return &defaultRequester{
		Socket:   sock,
		ch:       ch,
		uid:      uid,
		suppSock: supSock,
		srvCh:    srvChan,
		deadline: time.Second * 10,
		ddlnChan: ddlnChan,
		mx:       mx,
	}, ch
}
