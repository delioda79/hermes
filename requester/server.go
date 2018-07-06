package requester

import (
	"encoding/json"
	"fmt"
	"sync"

	"bitbucket.org/ConsentSystems/mango-micro/mango-service/registry/consul"
	"bitbucket.org/ConsentSystems/mango-micro/mango-service/service"
	"bitbucket.org/ConsentSystems/mango-micro/messages"
	"github.com/hashicorp/consul/api"
	"github.com/pborman/uuid"
	mangos "nanomsg.org/go-mangos"
	"nanomsg.org/go-mangos/protocol/req"
)

// Server is the pusher
type Server interface {
	AddTransport(tr mangos.Transport)
	Sock() Requester
	SockBare() mangos.Socket
	Run(pbs ...Responder)
}
type defaultServer struct {
	client        service.Client
	requesters    map[string]chan mangos.Message
	supportSocket mangos.Socket
	mutex         *sync.Mutex
}

// AddTransport adds a transport to the subscriber's socket
func (reqs *defaultServer) AddTransport(tr mangos.Transport) {
	reqs.client.Sock().AddTransport(tr)
}

func (reqs *defaultServer) Sock() Requester {
	uid := uuid.NewUUID().String()
	sock, ch := NewRequester(uid, reqs.client.Sock(), reqs.supportSocket)
	reqs.mutex.Lock()
	reqs.requesters[uid] = ch
	reqs.mutex.Unlock()
	return sock
}

func (reqs *defaultServer) SockBare() mangos.Socket {
	return reqs.client.Sock()
}

// Run runs the subscriber
func (reqs *defaultServer) Run(pbs ...Responder) {
	for _, p := range pbs {
		fmt.Println("Connecting")
		reqs.client.Connect(p.Name, p.Version, p.Protocol)
	}
	for {
		fmt.Println("Receiver is waiting")
		mgMsg, err := reqs.client.Sock().RecvMsg()
		if err != nil {
			fmt.Println("Error while retrieving msg", err)
		}
		trigger := &messages.Trigger{}
		json.Unmarshal(mgMsg.Body, trigger)
		reqs.mutex.Lock()
		_, ok := reqs.requesters[trigger.UID]
		if ok {
			reqs.requesters[trigger.UID] <- *mgMsg
			delete(reqs.requesters, trigger.UID)
		}
		reqs.mutex.Unlock()
	}
}

// NewServer returns a new Subscriber
func NewServer(regAddr string) (Server, error) {
	pushSock, err := req.NewSocket()
	if err != nil {
		return nil, err
	}

	supportSock, err := req.NewSocket()
	if err != nil {
		return nil, err
	}

	if err = pushSock.SetOption(mangos.OptionRaw, true); err != nil {
		return nil, fmt.Errorf("can't set raw mode: %s", err.Error())
	}

	registry := consul.NewRegistry(&api.Config{
		Address: regAddr,
		Scheme:  "http",
	})

	client := service.NewMangoClient(pushSock, registry)

	return &defaultServer{
		client:        client,
		requesters:    map[string]chan mangos.Message{},
		supportSocket: supportSock,
		mutex:         &sync.Mutex{},
	}, nil
}

type Responder struct {
	Name     string
	Protocol string
	Version  string
}
