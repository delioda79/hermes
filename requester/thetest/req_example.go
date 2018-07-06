package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"bitbucket.org/ConsentSystems/mango-micro/handler"
	"bitbucket.org/ConsentSystems/mango-micro/replier"
	"bitbucket.org/ConsentSystems/mango-micro/requester"
	"nanomsg.org/go-mangos/transport/inproc"
)

type ourStruct struct {
	Message string
}

func CreateReplier() {
	rep, err := replier.NewServer(":8500", "reptest", "1")
	if err != nil {
		fmt.Println("Error", err)
		os.Exit(1)
	}
	rep.AddTransport(inproc.NewTransport())
	hdl := handler.NewHandler()
	hdl.Add("test1", func(msg interface{}, rsp ...*[]byte) error {
		rcvd := msg.([]byte)
		rqMsg := &ourStruct{}
		_ = json.Unmarshal(rcvd, rqMsg)
		response := &ourStruct{
			Message: "HELLO " + rqMsg.Message,
		}
		bts, _ := json.Marshal(response)
		*rsp[0] = bts

		return nil
	})
	rep.AddHandler(hdl)
	rep.Run(900, "inproc", "reptest")
}

func CreateRequester() requester.Server {
	rqs, _ := requester.NewServer(":8500")

	rqs.AddTransport(inproc.NewTransport())
	go rqs.Run(requester.Responder{
		Name:     "reptest",
		Protocol: "inproc",
		Version:  "1",
	})
	time.Sleep(time.Second * 5)
	fmt.Println("OK")
	return rqs
}

func SendMessage(rqs requester.Server, name string) {
	reqStuff := ourStruct{
		Message: name,
	}

	bts, _ := json.Marshal(reqStuff)
	duration := time.Millisecond * time.Duration(rand.Int63n(100))
	//fmt.Println("Waiting ", strconv.FormatInt(int64(duration), 10))
	time.Sleep(duration)
	fmt.Println("Sending ", name)
	rsp, err := rqs.Sock().Request("test1", bts)
	if err != nil {
		fmt.Println(err)
		return
	}
	rspMsg := &ourStruct{}
	json.Unmarshal(rsp, rspMsg)
	fmt.Println(name+" received teh repsonse ", rspMsg, string(rsp))
}

func main() {
	go CreateReplier()
	time.Sleep(time.Second * 5)
	rqs := CreateRequester()
	for i := 0; i < 10000; i++ {
		go SendMessage(rqs, "James"+strconv.FormatInt(int64(i), 10))
	}
	select {}
}
