package example1

import (
	"encoding/json"
	"fmt"
	"strconv"

	"bitbucket.org/ddanna79/mango-micro/handler"
	"bitbucket.org/ddanna79/mango-micro/puller"
	"nanomsg.org/go-mangos/transport/inproc"
	"nanomsg.org/go-mangos/transport/tcp"
)
import "gopkg.in/mgo.v2/bson"

// NewAPICallsHandlerServer returns a new puller server
func NewAPICallsHandlerServer (
	discoveryAddr string,
	portStr string,
	hdl  APICallsHandler,
) (puller.Server, error) {
	pullserver, _ := puller.NewServer(
		discoveryAddr,
		"APICallsHandlerServer-puller",
		"1",
	)

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("wrong port %s", portStr)
	}
	pullserver.AddTransport(inproc.NewTransport())
	pullserver.AddTransport(tcp.NewTransport())
	go pullserver.Run(port, "inproc", "APICallsHandlerServer-puller")
	go pullserver.Run(port, "tcp", "")

	handler := handler.NewHandler()

	
	handler.Add("APICallsHandler.RegisterCall ", func(msg interface{}, rsp ...*[]byte) error {
		inParam := &APICallMessage{}
		arg, ok := msg.([]byte)

		if !ok {
			fmt.Printf("Wrong message sent %v", arg)
			return fmt.Errorf("Wrong message sent %v", arg)
		}

		err := json.Unmarshal(arg, inParam)
		if err != nil {
			fmt.Println("Error unmarshaling: ", err)
			return err
		}

		hdl.RegisterCall(inParam)
		return nil
	})


	handler.Add("APICallsHandler.NoParamsCall ", func(msg interface{}, rsp ...*[]byte) error {
		hdl.NoParamsCall()
		return nil
	})


	pullserver.AddHandler(handler)

	return pullserver, nil
}

