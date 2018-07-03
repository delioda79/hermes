package publisher

import (
	"bitbucket.org/ConsentSystems/mango-micro/mango-service/registry/consul"
	"bitbucket.org/ConsentSystems/mango-micro/mango-service/service"
	"github.com/hashicorp/consul/api"
	mangos "nanomsg.org/go-mangos"
)

// Server is the server
type Server interface {
	AddTransport(tr mangos.Transport)
	Run(port int, transport, addr string)
	Sock() Publisher
}
type defaultServer struct {
	server  service.Server
	name    string
	version string
}

// AddTransport adds a transport to the subscriber's socket
func (pubs *defaultServer) AddTransport(tr mangos.Transport) {
	pubs.server.Sock().AddTransport(tr)
}

func (pubs *defaultServer) Sock() Publisher {
	sock := pubs.server.Sock().(Publisher)
	return sock
}

// Run runs the subscriber
func (pubs *defaultServer) Run(port int, transport, addr string) {
	pubs.server.Run(
		pubs.name,
		addr,
		port,
		transport,
		pubs.version,
	)
}

func NewServer(
	regAddr string,
	serviceName string,
	version string,
) (Server, error) {
	pubSock, err := NewPublisher()
	if err != nil {
		return nil, err
	}
	registry := consul.NewRegistry(&api.Config{
		Address: regAddr,
		Scheme:  "http",
	})

	server := service.NewMangoServer(pubSock, registry)

	return &defaultServer{
		server:  server,
		name:    serviceName,
		version: version,
	}, nil
}
