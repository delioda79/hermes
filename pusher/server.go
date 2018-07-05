package pusher

import (
	"encoding/json"

	"bitbucket.org/ConsentSystems/mango-micro/mango-service/registry/consul"
	"bitbucket.org/ConsentSystems/mango-micro/mango-service/service"
	"bitbucket.org/ConsentSystems/mango-micro/messages"
	"github.com/hashicorp/consul/api"
	mangos "nanomsg.org/go-mangos"
)

// Server is the pusher
type Server interface {
	AddTransport(tr mangos.Transport)
	Sock() Pusher
	Run(pbs ...Puller)
}
type defaultServer struct {
	client service.Client
}

// AddTransport adds a transport to the subscriber's socket
func (sub *defaultServer) AddTransport(tr mangos.Transport) {
	sub.client.Sock().AddTransport(tr)
}

func (sub *defaultServer) Sock() Pusher {
	sock := sub.client.Sock().(Pusher)
	return sock
}

func (sub *defaultServer) Push(name string, message []byte) error {
	trg := &messages.Trigger{
		Name:   name,
		Params: message,
	}

	bts, err := json.Marshal(trg)
	if err != nil {
		return err
	}
	sub.client.Send(bts)
	return nil
}

// Run runs the subscriber
func (sub *defaultServer) Run(pbs ...Puller) {
	for _, p := range pbs {
		sub.client.Connect(p.Name, p.Version, p.Protocol)
	}
}

// NewServer returns a new Subscriber
func NewServer(regAddr string) (Server, error) {
	pushSock, err := NewPusher()
	if err != nil {
		return nil, err
	}

	registry := consul.NewRegistry(&api.Config{
		Address: regAddr,
		Scheme:  "http",
	})

	client := service.NewMangoClient(pushSock, registry)

	return &defaultServer{
		client: client,
	}, nil
}

type Puller struct {
	Name     string
	Protocol string
	Version  string
}
