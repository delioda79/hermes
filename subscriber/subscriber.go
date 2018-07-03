package subscriber

import (
	"flag"
	"log"

	"bitbucket.org/ConsentSystems/mango-micro/mango-service/registry/consul"
	"bitbucket.org/ConsentSystems/mango-micro/mango-service/service"
	mangos "nanomsg.org/go-mangos"
	"nanomsg.org/go-mangos/protocol/sub"
	"nanomsg.org/go-mangos/transport/tcp"
)

type Subscriber interface {
	AddTransport(tr mangos.Transport)
}
type defaultSubscriber struct {
	receiveSocket mangos.Socket
	client        service.Client
}

// AddTransport adds a transport to the workers pool manager
func (sub *defaultSubscriber) AddTransport(tr mangos.Transport) {
	sub.client.Sock().AddTransport(tr)
}

func NewSubscriber() service.Client {
	var subSock mangos.Socket
	var err error

	flag.Parse()

	if subSock, err = sub.NewSocket(); err != nil {
		log.Fatal("can't get new req socket: ", err.Error())
	}
	subSock.AddTransport(tcp.NewTransport())

	registry := consul.NewRegistry(nil)

	return service.NewMangoClient(subSock, registry)
}
