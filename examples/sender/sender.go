package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	"bitbucket.org/ddanna79/mango-micro/mango-service/registry/consul"
	"bitbucket.org/ddanna79/mango-micro/mango-service/service"
	"nanomsg.org/go-mangos"
	"nanomsg.org/go-mangos/protocol/pull"
	"nanomsg.org/go-mangos/protocol/push"
	"nanomsg.org/go-mangos/transport/tcp"
)

func main() {

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
	fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAaID", server.GetID())
	client := service.NewMangoClient(pushSock, registry)
	client.Connect("subscription-server", "1", "tcp")

	server.Run("subscription-client", "", *port, "tcp")

	c := make(chan int, 1)

	go func() {
		for i := 0; i < 20; i++ {
			//time.Sleep(time.Second * 3)
			go func(c chan int, val int) {
				time.Sleep(time.Second * 3)
				message := "DATE" + strconv.Itoa(val)
				fmt.Println("Sending", message)
				if err = client.Send([]byte(message)); err != nil {
					log.Fatal("can't send message on push socket: ", err.Error())
				}
			}(c, i)
		}

	}()

	cnt := 0
	for {
		fmt.Println(server.GetID())
		if msg, err = server.Receive(); err != nil {
			log.Fatal("can't receive date: ", err.Error())
		}
		fmt.Println("Received ", string(msg))
		cnt++
		if cnt == 20 {
			server.Deregister()
			return
		}
		//}
	}
}
