package main

import (
	"flag"
	"fmt"
	"log"

	"bitbucket.org/ConsentSystems/mango-micro/mango-service/registry/consul"
	"bitbucket.org/ConsentSystems/mango-micro/mango-service/service"
	"github.com/nanomsg/mangos"
	"github.com/nanomsg/mangos/protocol/pull"
	"github.com/nanomsg/mangos/protocol/push"
	"github.com/nanomsg/mangos/transport/tcp"
)

func main() {
	// CREATING THE REP SOCKET
	var pullSock, pushSock mangos.Socket
	var err error
	var msg []byte

	var port = flag.Int("port", 9008, "port to run the server")
	flag.Parse()

	if pullSock, err = pull.NewSocket(); err != nil {
		log.Fatal("can't get new req socket: ", err.Error())
	}
	pullSock.AddTransport(tcp.NewTransport())

	if pushSock, err = push.NewSocket(); err != nil {
		log.Fatal("can't get new req socket: ", err.Error())
	}
	pushSock.AddTransport(tcp.NewTransport())

	registry := consul.NewRegistry(nil)

	server := service.NewMangoServer(pullSock, registry)
	client := service.NewMangoClient(pushSock, registry)
	client.Connect("subscription-client", "1", "tcp")

	server.Run("subscription-server", "", *port, "tcp")

	for {
		// Could also use sock.RecvMsg to get header
		msg, err = server.Receive()
		//if string(msg) == "DATE" { // no need to terminate
		fmt.Println("Received", string(msg))
		err = client.Send(msg)
		if err != nil {
			log.Fatal("can't send reply: ", err.Error())
		}
		//}
	}

}
