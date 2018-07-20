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
	OpenSockets() int
	TotRcv() int
}
type defaultServer struct {
	client        service.Client
	requesters    map[string]chan mangos.Message
	supportSocket mangos.Socket
	srvMutex      *sync.Mutex
	rcvChan       chan string
	ttrcv         int
	sckMutex      *sync.Mutex
}

// AddTransport adds a transport to the subscriber's socket
func (reqs *defaultServer) AddTransport(tr mangos.Transport) {
	reqs.client.Sock().AddTransport(tr)
}

func (reqs *defaultServer) OpenSockets() int {
	return len(reqs.requesters)
}

func (reqs *defaultServer) TotRcv() int {
	reqs.srvMutex.Lock()
	res := reqs.ttrcv
	reqs.srvMutex.Unlock()
	return res
}

func (reqs *defaultServer) Sock() Requester {
	uid := uuid.NewUUID().String()
	sock, ch := NewRequester(uid, reqs.client.Sock(), reqs.supportSocket, reqs.rcvChan, reqs.sckMutex)
	reqs.srvMutex.Lock()
	reqs.requesters[uid] = ch
	reqs.srvMutex.Unlock()
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

	go func(ch chan string) {
		for {
			uid := <-ch
			reqs.srvMutex.Lock()
			_, ok := reqs.requesters[uid]
			if ok {
				delete(reqs.requesters, uid)
			}
			reqs.srvMutex.Unlock()
		}
	}(reqs.rcvChan)
	for {
		fmt.Println("Receiver is waiting")
		mgMsg, err := reqs.client.Sock().RecvMsg()
		reqs.srvMutex.Lock()
		reqs.ttrcv++
		if err != nil {
			fmt.Println("Error while retrieving msg", err)
		}
		trigger := &messages.Trigger{}
		json.Unmarshal(mgMsg.Body, trigger)
		_, ok := reqs.requesters[trigger.UID]
		if ok {
			reqs.requesters[trigger.UID] <- *mgMsg
		}
		reqs.srvMutex.Unlock()
	}
}

// NewServer returns a new Subscriber
func NewServer(regAddr string) (Server, error) {
	rcvChan := make(chan string, 1)
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
		srvMutex:      &sync.Mutex{},
		rcvChan:       rcvChan,
		sckMutex:      &sync.Mutex{},
		ttrcv:         0,
	}, nil
}

type Responder struct {
	Name     string
	Protocol string
	Version  string
}
