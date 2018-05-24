package service

import (
	"bitbucket.org/ConsentSystems/subscription-service/mango-service/registry"
	"github.com/go-mangos/mangos"
)

// Service represents a microservice base
type Service interface {
	Send(msg []byte) error
	Receive() ([]byte, error)
	Sock() mangos.Socket
	Registry(msg []byte) registry.Registry
}

// MangoService represents a service using mangos sockets
type MangoService struct {
	socket   mangos.Socket
	registry registry.Registry
	ID       string
}

// NewMangoService creates a new mango service
func NewMangoService(sock mangos.Socket, reg registry.Registry) Service {

	return &MangoService{
		socket:   sock,
		registry: reg,
	}
}

// Send sends a message to the socket
func (mgs MangoService) Send(msg []byte) error {
	return mgs.socket.Send(msg)
}

// Receive receives a message to the socket
func (mgs MangoService) Receive() ([]byte, error) {
	return mgs.socket.Recv()
}

// Sock returns theinternal socket
func (mgs MangoService) Sock() mangos.Socket {
	return mgs.socket
}

// Registry returns theinternal registry
func (mgs MangoService) Registry(msg []byte) registry.Registry {
	return mgs.registry
}
