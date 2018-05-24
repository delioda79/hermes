package service

import (
	"fmt"
	"time"

	"bitbucket.org/ConsentSystems/subscription-service/mango-service/registry"
	"github.com/go-mangos/mangos"
)

// Client is a service client
type Client interface {
	Service
	Connect(name string, version string, transport string)
}

// MangoClient is a mango service client
type MangoClient struct {
	MangoService
}

// NewMangoClient creates a new mango service
func NewMangoClient(sock mangos.Socket, reg registry.Registry) Client {

	service := NewMangoService(sock, reg).(*MangoService)
	return &MangoClient{
		MangoService: *service,
	}
}

// Connect connects to another socket
func (mgc MangoClient) Connect(name string, version string, transport string) {
	go mgc.connect(name, version, transport)
}

// Connect connects to another socket
func (mgc MangoClient) connect(name string, version string, transport string) error {
	for {
		urls, err := mgc.registry.Get(name, version)
		if err != nil {
			continue
		}
		for _, url := range urls {
			connStr := fmt.Sprintf("%s://%s", transport, url)
			fmt.Println("CONN", connStr)
			err := mgc.socket.Dial(connStr)
			if err != nil {
				fmt.Println("Conn error", err.Error())
			}
		}
		//LOOK at this for connection drops
		//mgc.socket.SetPortHook()
		time.Sleep(time.Minute)
	}
}
