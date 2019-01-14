package subscriber

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"bitbucket.org/ddanna79/mango-micro/handler"
	"bitbucket.org/ddanna79/mango-micro/mango-service/registry/consul"
	"bitbucket.org/ddanna79/mango-micro/mango-service/service"
	"bitbucket.org/ddanna79/mango-micro/messages"
	"github.com/hashicorp/consul/api"
	mangos "nanomsg.org/go-mangos"
	"nanomsg.org/go-mangos/protocol/sub"
)

// Subscriber is the subscriber
type Subscriber interface {
	AddTransport(tr mangos.Transport)
	Run(pbs ...Publisher)
	AddHandler(handler handler.Handler)
}
type defaultSubscriber struct {
	client   service.Client
	handlers []handler.Handler
}

// AddTransport adds a transport to the subscriber's socket
func (sub *defaultSubscriber) AddTransport(tr mangos.Transport) {
	sub.client.Sock().AddTransport(tr)
}

// Run runs the subscriber
func (sub *defaultSubscriber) Run(pbs ...Publisher) {
	for _, p := range pbs {
		sub.client.Connect(p.Name, p.Version, p.Protocol)
	}

	err := sub.client.Sock().SetOption(mangos.OptionSubscribe, []byte(""))
	if err != nil {
		fmt.Println("Impossible to subscribe", err)
		os.Exit(1)
	}
	for {
		bts, err := sub.client.Receive()
		if err != nil {
			fmt.Println(err)
		}

		msg := &messages.Trigger{}
		err = json.Unmarshal(bts, msg)
		if err != nil {
			fmt.Println("error unmsrshaling", err)
			continue
		}

		for _, hdl := range sub.handlers {
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

// AddHandler adds a hand;er to the subscriber
func (sub *defaultSubscriber) AddHandler(handler handler.Handler) {
	sub.handlers = append(sub.handlers, handler)
}

// NewSubscriber returns a new Subscriber
func NewSubscriber(regAddr string) Subscriber {
	var subSock mangos.Socket
	var err error

	if subSock, err = sub.NewSocket(); err != nil {
		log.Fatal("can't get new req socket: ", err.Error())
	}

	registry := consul.NewRegistry(&api.Config{
		Address: regAddr,
		Scheme:  "http",
	})

	client := service.NewMangoClient(subSock, registry)

	return &defaultSubscriber{
		client:   client,
		handlers: []handler.Handler{},
	}
}

type Publisher struct {
	Name     string
	Protocol string
	Version  string
}
