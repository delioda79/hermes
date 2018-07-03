package service

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"bitbucket.org/ConsentSystems/mango-micro/mango-service/registry"
	"nanomsg.org/go-mangos"
)

// Server represents a microservice base
type Server interface {
	Service
	Run(name, addr string, port int, transport string, version string)
	Deregister()
	GetID() string
}

// MangoServer represents a service using mangos sockets
type MangoServer struct {
	MangoService
}

// NewMangoServer creates a new mango service
func NewMangoServer(sock mangos.Socket, reg registry.Registry) Server {
	service := NewMangoService(sock, reg).(*MangoService)
	return &MangoServer{
		MangoService: *service,
	}
}

// Run runs the service
func (mgs *MangoServer) Run(name, addr string, port int, transport string, version string) {

	url := fmt.Sprintf("%s://%s:%d", transport, addr, port)
	if err := mgs.socket.Listen(url); err != nil {
		log.Fatal("can't listen on rep socket:", err.Error())
	}
	if version == "" {
		version = "1"
	}

	tags := []string{"v=" + version}
	sID, err := mgs.registry.Register(name, "", port, tags)
	if err != nil {
		log.Fatal("Wrong registration ", err.Error())
	}
	fmt.Println("Setting ID: ", sID)
	mgs.ID = sID

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGKILL)

	go mgs.exit(sigChan)
}

func (mgs MangoServer) exit(sigChan chan os.Signal) {
	fmt.Println("Entering loop")
	for {
		<-sigChan
		mgs.Deregister()
		os.Exit(1)
	}
}

// Deregister deregisters the service
func (mgs MangoServer) Deregister() {
	fmt.Println("Deregistering ", mgs.ID)
	err := mgs.registry.Deregister(mgs.ID)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (mgs MangoServer) GetID() string {
	return mgs.ID
}
