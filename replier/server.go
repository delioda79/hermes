package replier

import (
	"encoding/json"
	"fmt"

	"bitbucket.org/ConsentSystems/mango-micro/handler"
	"bitbucket.org/ConsentSystems/mango-micro/mango-service/registry/consul"
	"bitbucket.org/ConsentSystems/mango-micro/mango-service/service"
	"bitbucket.org/ConsentSystems/mango-micro/messages"
	"github.com/hashicorp/consul/api"
	mangos "nanomsg.org/go-mangos"
	"nanomsg.org/go-mangos/protocol/rep"
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

// AddHandler adds a hand;er to the subscriber
func (reps *defaultServer) AddHandler(handler handler.Handler) {
	reps.handlers = append(reps.handlers, handler)
}

// AddTransport adds a transport to the subscriber's socket
func (reps *defaultServer) AddTransport(tr mangos.Transport) {
	reps.server.Sock().AddTransport(tr)
}

// Run runs the subscriber
func (reps *defaultServer) Run(port int, transport, addr string) {
	reps.server.Run(
		reps.name,
		addr,
		port,
		transport,
		reps.version,
	)
	for {
		bts, err := reps.server.Sock().Recv()
		if err != nil {
			fmt.Println(err)
		}

		msg := &messages.Trigger{}
		err = json.Unmarshal(bts, msg)
		if err != nil {
			fmt.Println("error unmsrshaling", err)
			continue
		}

		for _, hdl := range reps.handlers {
			go func(hdl handler.Handler) {
				defer func() {
					if r := recover(); r != nil {
						fmt.Println("Recovered from: ", r)
					}
				}()
				rsp := []byte{}
				hdl.Run(msg.Name, msg.Params, rsp)
				reps.server.Sock().Send(rsp)
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

	pullSock, err := rep.NewSocket()
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

	return nil, nil
}
