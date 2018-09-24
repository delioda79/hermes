package example1

import (
	"encoding/json"
	"fmt"

	"bitbucket.org/ConsentSystems/mango-micro/handler"
	"bitbucket.org/ConsentSystems/mango-micro/puller"
	"nanomsg.org/go-mangos/transport/inproc"
	"nanomsg.org/go-mangos/transport/tcp"
)
import "bitbucket.org/ConsentSystems/mango-micro/mango-service/registry"

// NewAPICallsHandlerServer returns a new puller server
func NewAPICallsHandlerServer(
	registry registry.Registry,
	portStr string,
	hdl APICallsHandler,
	serviceName string,
) (puller.Server, error) {

	serviceNmsp := serviceName
	if serviceName == "" {
		serviceNmsp = "APICallsHandler"
	}

	serviceName = serviceNmsp + "Server"

	pullserver, _ := puller.NewServer(
		registry,
		serviceName+"-puller",
		"1",
	)

	pullserver.AddTransport(inproc.NewTransport())
	pullserver.AddTransport(tcp.NewTransport())

	handler := handler.NewHandler()

	handler.Add(serviceNmsp+".RegisterCall", func(msg interface{}, rsp ...*[]byte) error {
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

	handler.Add(serviceNmsp+".NoParamsCall ", func(msg interface{}, rsp ...*[]byte) error {
		hdl.NoParamsCall()
		return nil
	})

	pullserver.AddHandler(handler)

	return pullserver, nil
}
