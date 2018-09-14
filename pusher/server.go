package pusher

import (
	"bitbucket.org/ConsentSystems/mango-micro/mango-service/registry"
	"bitbucket.org/ConsentSystems/mango-micro/mango-service/service"
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

// Run runs the subscriber
func (sub *defaultServer) Run(pbs ...Puller) {
	for _, p := range pbs {
		sub.client.Connect(p.Name, p.Version, p.Protocol)
	}
}

// NewServer returns a new Subscriber
func NewServer(registry registry.Registry) (Server, error) {
	pushSock, err := NewPusher()
	if err != nil {
		return nil, err
	}

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
