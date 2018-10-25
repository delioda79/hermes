package service

import (
	"time"

	"bitbucket.org/ConsentSystems/logging"
	"bitbucket.org/ConsentSystems/mango-micro/logger"
	"bitbucket.org/ConsentSystems/mango-micro/mango-service/registry"
	"nanomsg.org/go-mangos"
)

// Service represents a microservice base
type Service interface {
	Send(msg []byte) error
	Receive() ([]byte, error)
	Sock() mangos.Socket
	Registry() registry.Registry
	Logger() logging.Logger
	SetLogger(lgr logging.Logger)
}

// MangoService represents a service using mangos sockets
type MangoService struct {
	socket        mangos.Socket
	registry      registry.Registry
	ID            string
	logger        logging.Logger
	keepAliveTime time.Duration
}

// NewMangoService creates a new mango service
func NewMangoService(
	sock mangos.Socket,
	reg registry.Registry,
) Service {

	return &MangoService{
		socket:        sock,
		registry:      reg,
		logger:        logger.NewBasicLogger(),
		keepAliveTime: time.Second * 20,
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
func (mgs *MangoService) Sock() mangos.Socket {
	return mgs.socket
}

// Registry returns theinternal registry
func (mgs *MangoService) Registry() registry.Registry {
	return mgs.registry
}

// Logger returns the internal logger
func (mgs *MangoService) Logger() logging.Logger {
	return mgs.logger
}

// SetLogger sets the internl logger
func (mgs *MangoService) SetLogger(lgr logging.Logger) {
	mgs.logger = lgr
}

func (mgs *MangoService) SetKeepAlive(kplt time.Duration) {
	mgs.keepAliveTime = kplt
}
