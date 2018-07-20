package replier

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"bitbucket.org/ConsentSystems/mango-micro/handler"
	"bitbucket.org/ConsentSystems/mango-micro/mango-service/registry/consul"
	"bitbucket.org/ConsentSystems/mango-micro/mango-service/service"
	"bitbucket.org/ConsentSystems/mango-micro/messages"
	"github.com/hashicorp/consul/api"
	mangos "nanomsg.org/go-mangos"
	"nanomsg.org/go-mangos/protocol/rep"
)

//var mx sync.Mutex

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
	mx       sync.Mutex
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
		//fmt.Println("WAITING FOR MSG")
		rawMsg, err := reps.server.Sock().RecvMsg()
		if err != nil {
			fmt.Println(err)
		}

		for _, hdl := range reps.handlers {
			go func(hdl handler.Handler, origMsg *mangos.Message) {
				reps.mx.Lock()
				defer func() {
					if r := recover(); r != nil {
						fmt.Println("Recovered from: ", r)
					}
				}()

				rsp := []*[]byte{&[]byte{}, &[]byte{}}
				bts := rawMsg.Body
				msg := &messages.Trigger{}
				err = json.Unmarshal(bts, msg)
				if err != nil {
					*rsp[1] = []byte(err.Error())
					fmt.Println("error unmsrshaling", err)
				} else {
					hdl.Run(msg.Name, msg.Params, rsp...)
				}
				body, _ := json.Marshal(rsp)
				msg.Params = body
				bts, _ = json.Marshal(msg)
				response := mangos.NewMessage(len(bts))
				//response := mangos.Message(*origMsg)
				response.Body = bts
				response.Header = (*origMsg).Header
				time.Sleep(time.Microsecond * 5)
				err := reps.server.Sock().SendMsg(response)
				if err != nil {
					fmt.Println("SOmething went wrong: ", err)
				}
				reps.mx.Unlock()
			}(hdl, rawMsg)
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

	repSock, err := rep.NewSocket()
	if err != nil {
		return nil, err
	}

	if err = repSock.SetOption(mangos.OptionRaw, true); err != nil {
		return nil, fmt.Errorf("can't set raw mode: %s", err.Error())
	}
	server := service.NewMangoServer(repSock, registry)

	return &defaultServer{
		server:   server,
		name:     serviceName,
		version:  version,
		handlers: []handler.Handler{},
	}, nil

	return nil, nil
}
