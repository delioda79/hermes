package service

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"bitbucket.org/ConsentSystems/logging"
	"bitbucket.org/ConsentSystems/mango-micro/mango-service/registry"
	"nanomsg.org/go-mangos"
)

// Server represents a microservice base
type Server interface {
	Service
	Run(name, addr string, port int, transport string, version string)
	Deregister()
	GetID() string
	Stop()
}

// MangoServer represents a service using mangos sockets
type MangoServer struct {
	MangoService
	logger logging.Logger
}

// NewMangoServer creates a new mango service
func NewMangoServer(
	sock mangos.Socket,
	reg registry.Registry,
) Server {
	service := NewMangoService(sock, reg).(*MangoService)
	return &MangoServer{
		MangoService: *service,
	}
}

// Run runs the service
func (mgs *MangoServer) Run(name, addr string, port int, transport string, version string) {

	url := fmt.Sprintf("%s://%s:%d", transport, addr, port)
	fmt.Println("URL IS ", url)

	lstnr, err := mgs.socket.NewListener(url, map[string]interface{}{})
	if err != nil {
		os.Exit(1)
	}
	lstnr.SetOption(mangos.OptionKeepAlive, true)
	lstnr.SetOption(mangos.OptionKeepAliveTime, mgs.keepAliveTime)

	if err := lstnr.Listen(); err != nil {
		log.Fatal("can't listen on socket:", err.Error())
	}

	if version == "" {
		version = "1"
	}

	tags := []string{
		"v=" + version,
		"transport=" + transport,
	}
	sID, err := mgs.registry.Register(name, addr, port, tags)
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

func (mgs MangoServer) Stop() {
	mgs.Sock().Close()
	mgs.Deregister()
	return
}
