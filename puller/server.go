package puller

import (
	"encoding/json"
	"fmt"

	"bitbucket.org/ConsentSystems/mango-micro/handler"
	"bitbucket.org/ConsentSystems/mango-micro/mango-service/registry/consul"
	"bitbucket.org/ConsentSystems/mango-micro/mango-service/service"
	"bitbucket.org/ConsentSystems/mango-micro/messages"
	"github.com/hashicorp/consul/api"
	mangos "nanomsg.org/go-mangos"
	"nanomsg.org/go-mangos/protocol/pull"
)

// Server is the server
type Server interface {
	AddTransport(tr mangos.Transport)
	AddHandler(handler handler.Handler)
	Run(port int, transport, addr string)
}
type defaultServer struct {
	server   service.Server
	name     string
	version  string
	handlers []handler.Handler
}

// AddTransport adds a transport to the subscriber's socket
func (pubs *defaultServer) AddTransport(tr mangos.Transport) {
	pubs.server.Sock().AddTransport(tr)
}

// AddHandler adds a hand;er to the subscriber
func (pubs *defaultServer) AddHandler(handler handler.Handler) {
	pubs.handlers = append(pubs.handlers, handler)
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
	for {
		bts, err := pubs.server.Sock().Recv()
		if err != nil {
			fmt.Println(err)
		}

		msg := &messages.Trigger{}
		err = json.Unmarshal(bts, msg)
		if err != nil {
			fmt.Println("error unmsrshaling", err)
			continue
		}

		for _, hdl := range pubs.handlers {
			go func(hdl handler.Handler) {
				defer func() {
					if r := recover(); r != nil {
						fmt.Println("Recovered from: ", r)
					}
				}()
				hdl.Run(msg.Name, msg.Params)
			}(hdl)
		}
	}

}

func NewServer(
	regAddr string,
	serviceName string,
	version string,
) (Server, error) {
	registry := consul.NewRegistry(&api.Config{
		Address: regAddr,
		Scheme:  "http",
	})

	pullSock, err := pull.NewSocket()
	if err != nil {
		return nil, err
	}
	server := service.NewMangoServer(pullSock, registry)

	return &defaultServer{
		server:   server,
		name:     serviceName,
		version:  version,
		handlers: []handler.Handler{},
	}, nil
}
