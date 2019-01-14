package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	"bitbucket.org/ddanna79/mango-micro/handler"
	"bitbucket.org/ddanna79/mango-micro/replier"
	"bitbucket.org/ddanna79/mango-micro/requester"
	"nanomsg.org/go-mangos/transport/inproc"
)

var mx sync.Mutex

type ourStruct struct {
	Message string
}

func CreateReplier(port int, name string, ch chan bool) {
	fmt.Println("Creating ", port)
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
			Message: "HELLO " + rqMsg.Message + " from " + name,
		}
		bts, _ := json.Marshal(response)
		*rsp[0] = bts
		ch <- true
		fl := rand.Int63n(10000)%10000 == 0
		if fl {
			fmt.Println("Yes")
			duration := time.Second * 1
			time.Sleep(duration)
		} else {
			fmt.Println("NO")
		}

		return nil
	})
	rep.AddHandler(hdl)
	rep.Run(port, "inproc", "reptest")
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

func SendMessage(rqs requester.Server, name string, ch chan bool) {
	reqStuff := ourStruct{
		Message: name,
	}

	bts, _ := json.Marshal(reqStuff)
	//duration := time.Millisecond * (time.Duration(rand.Int63n(100)) + 1)
	//fmt.Println("Waiting ", strconv.FormatInt(int64(duration), 10))
	//time.Sleep(duration)
	//fmt.Println("Sending ", name)
	//mx.Lock()
	sck := rqs.Sock()
	sck.SetDeadline(0)
	rsp, err := sck.Request("test1", bts)
	if err != nil {
		fmt.Println("An error occurred: ", err)
		ch <- false
		return
	}
	//mx.Unlock()
	rspMsg := &ourStruct{}
	json.Unmarshal(rsp, rspMsg)
	fmt.Println(name+" received teh repsonse ", rspMsg, string(rsp))
	time.Sleep(time.Second * 3)
	ch <- true
}

func main() {
	max := 1000000
	ch := make(chan bool, 1)
	ch2 := make(chan bool, 1)
	go CreateReplier(900, "Replier1", ch2)
	go CreateReplier(901, "Replier2", ch2)
	time.Sleep(time.Second * 5)
	rqs := CreateRequester()
	for i := 0; i < max; i++ {
		go SendMessage(rqs, "James"+strconv.FormatInt(int64(i), 10), ch)
	}

	cnt := 0
	yes := 0
	no := 0
	snt := 0
	for {
		select {
		case vl := <-ch:
			cnt++
			if vl {
				yes++
			} else {
				no++
			}
		case <-ch2:
			snt++
		}
		if cnt%100 == 0 || snt%100 == 0 {
			fmt.Println("Received Back: ", cnt, " OK: ", yes, " NO: ", no, " tot: ", yes+no)
			fmt.Println("Received by th ehandler: ", snt)
		}
		if cnt == max && snt == max {
			fmt.Println("We are out")
			break
		}
	}

	go func() {

		time.Sleep(time.Second * 5)
		for i := 0; i < max; i++ {
			time.Sleep(time.Second)
			fmt.Println("Open Sockets: ", rqs.OpenSockets())
			fmt.Println("Totl received: ", rqs.TotRcv(), " OK: ", yes, " NO: ", no, " tot: ", yes+no)
		}

	}()

	// time.Sleep(time.Second * 10)
	// go SendMessage(rqs, "James"+strconv.FormatInt(int64(max+1), 10), ch)
	select {}
}
