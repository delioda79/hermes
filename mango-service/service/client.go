package service

import (
	"fmt"
	"time"

	"bitbucket.org/ConsentSystems/mango-micro/mango-service/registry"
	"nanomsg.org/go-mangos"
)

// Client is a service client
type Client interface {
	Service
	Connect(name string, version string, transport string)
}

// MangoClient is a mango service client
type MangoClient struct {
	MangoService
	conns []string
}

// NewMangoClient creates a new mango service
func NewMangoClient(sock mangos.Socket, reg registry.Registry) Client {

	service := NewMangoService(sock, reg).(*MangoService)
	return &MangoClient{
		MangoService: *service,
		conns:        []string{},
	}
}

// Connect connects to another socket
func (mgc *MangoClient) Connect(name string, version string, transport string) {
	go mgc.connect(name, version, transport)
}

// Connect connects to another socket
func (mgc *MangoClient) connect(name string, version string, transport string) error {
	for {
		urls, err := mgc.registry.Get(name, version, transport)
		if err != nil {
			time.Sleep(time.Second * 10)
			continue
		}

	OUTER:
		for _, url := range urls {
			connStr := fmt.Sprintf("%s://%s", transport, url)
			for _, conn := range mgc.conns {
				if connStr == conn {
					continue OUTER
				}
			}
			err := mgc.socket.Dial(connStr)
			if err != nil {
				fmt.Println("Conn error", err.Error())
			}
			fmt.Println("CONNECTED ", connStr)
			mgc.conns = append(mgc.conns, connStr)
		}
		//LOOK at this for connection drops
		mgc.socket.SetPortHook(func(action mangos.PortAction, port mangos.Port) bool {
			if action == mangos.PortActionRemove {
				for i, v := range mgc.conns {
					if port.Address() == v {
						mgc.conns = append(mgc.conns[:i], mgc.conns[i+1:]...)
						port.Dialer().Close()
						fmt.Println("DROPPED ", port.Address())
					}
				}
			}

			return true
		})
		time.Sleep(time.Second * 10)
	}
}
