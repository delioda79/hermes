package subscriber

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"bitbucket.org/ConsentSystems/logging"
	"bitbucket.org/ConsentSystems/mango-micro/handler"
	"bitbucket.org/ConsentSystems/mango-micro/mango-service/registry"
	"bitbucket.org/ConsentSystems/mango-micro/mango-service/service"
	"bitbucket.org/ConsentSystems/mango-micro/messages"
	mangos "nanomsg.org/go-mangos"
	"nanomsg.org/go-mangos/protocol/sub"
)

// Subscriber is the subscriber
type Subscriber interface {
	AddTransport(tr mangos.Transport)
	Run(pbs ...Publisher)
	AddHandler(handler handler.Handler)
	Client() service.Client
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
		if sub.client.Logger() != nil {
			sub.client.Logger().Error(logging.Log{
				Code:   701,
				Status: "404",
				Detail: fmt.Sprintf(
					"Impossible to subscribe: %v",
					err,
				),
			})
		}
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
			if sub.client.Logger() != nil {
				sub.client.Logger().Error(logging.Log{
					Code:   704,
					Status: "404",
					Detail: fmt.Sprintf(
						"error unmarshaling: %v",
						err,
					),
				})
			}
			continue
		}

		for _, hdl := range sub.handlers {
			go func(hdl handler.Handler) {
				defer func() {
					if r := recover(); r != nil {
						if sub.client.Logger() != nil {
							sub.client.Logger().Fatal(logging.Log{
								Code:   700,
								Status: "500",
								Detail: fmt.Sprintf("Recovered from: %v", r),
							})
						}
					}
				}()
				err := hdl.Run(msg.Name, msg.Params)
				if sub.client.Logger() != nil {
					sub.client.Logger().Error(logging.Log{
						Code:   701,
						Status: "404",
						Detail: fmt.Sprintf(
							"Error while calling: %s with params %v: %v",
							msg.Name,
							msg.Params,
							err,
						),
					})
				}
			}(hdl)
		}
	}
}

// AddHandler adds a hand;er to the subscriber
func (sub *defaultSubscriber) AddHandler(handler handler.Handler) {
	sub.handlers = append(sub.handlers, handler)
}

// Client returns teh client
func (sub *defaultSubscriber) Client() service.Client {
	return sub.client
}

// NewSubscriber returns a new Subscriber
func NewSubscriber(registry registry.Registry) Subscriber {
	var subSock mangos.Socket
	var err error

	if subSock, err = sub.NewSocket(); err != nil {
		log.Fatal("can't get new req socket: ", err.Error())
	}

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
