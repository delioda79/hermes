package publisher

import (
	"encoding/json"
	"log"

	"bitbucket.org/ConsentSystems/mango-micro/mango-service/registry/consul"
	"bitbucket.org/ConsentSystems/mango-micro/mango-service/service"
	"bitbucket.org/ConsentSystems/mango-micro/messages"
	"github.com/hashicorp/consul/api"
	mangos "nanomsg.org/go-mangos"
	"nanomsg.org/go-mangos/protocol/sub"
)

// Publisher is the publisher
type Publisher interface {
	AddTransport(tr mangos.Transport)
	Run(port int, transport, addr string)
	Publish(name string, message []byte) error
}
type defaultPublisher struct {
	server  service.Server
	name    string
	version string
}

// AddTransport adds a transport to the subscriber's socket
func (pubs *defaultPublisher) AddTransport(tr mangos.Transport) {
	pubs.server.Sock().AddTransport(tr)
}

// Run runs the subscriber
func (pubs *defaultPublisher) Run(port int, transport, addr string) {
	pubs.server.Run(
		pubs.name,
		addr,
		port,
		transport,
		pubs.version,
	)
}

func (pubs *defaultPublisher) Publish(name string, message []byte) error {
	trg := &messages.Trigger{
		Name:   name,
		Params: message,
	}

	bts, err := json.Marshal(trg)
	if err != nil {
		return err
	}
	pubs.server.Sock().Send(bts)
	return nil
}

func NewPublisher(
	regAddr string,
	serviceName string,
	version string,
) Publisher {
	var pubSock mangos.Socket
	var err error

	if pubSock, err = sub.NewSocket(); err != nil {
		log.Fatal("can't get new req socket: ", err.Error())
	}

	registry := consul.NewRegistry(&api.Config{
		Address: regAddr,
		Scheme:  "http",
	})

	server := service.NewMangoServer(pubSock, registry)

	return &defaultPublisher{
		server: server,
	}
}
