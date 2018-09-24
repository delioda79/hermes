package example2

import (
	"encoding/json"
	"fmt"

	"bitbucket.org/ConsentSystems/mango-micro/handler"
	"bitbucket.org/ConsentSystems/mango-micro/mango-service/registry"
	"bitbucket.org/ConsentSystems/mango-micro/messages"
	"bitbucket.org/ConsentSystems/mango-micro/replier"
	"nanomsg.org/go-mangos/transport/inproc"
	"nanomsg.org/go-mangos/transport/tcp"
)

// NewAPICallsHandlerServer returns a new replier server
func NewAPICallsHandlerServer(
	registry registry.Registry,
	serverPort int,
	hdl APICallsHandler,
	serviceName string,
) (replier.Server, error) {
	serviceNmsp := serviceName
	if serviceName == "" {
		serviceNmsp = "APICallsHandler"
	}

	serviceName = serviceNmsp + "Server"

	replier, _ := replier.NewServer(registry, serviceName+"-replier", "1")
	handler := handler.NewHandler()

	handler.Add(serviceNmsp+".RegisterCall", func(in interface{}, out ...*[]byte) error {
		*out[0] = []byte{}
		*out[1] = []byte{}
		req := &APICallMessage{}
		err := json.Unmarshal(in.([]byte), req)
		if err != nil {
			*out[1] = []byte(err.Error())
			return err
		}

		rsp, err := hdl.RegisterCall(req)
		if err != nil {
			*out[1] = []byte(err.Error())
			return err
		}

		bts, err := json.Marshal(rsp)
		if err != nil {
			*out[1] = []byte(err.Error())
			return err
		}

		*out[0] = bts
		return nil
	})

	handler.Add(serviceNmsp+".External", func(in interface{}, out ...*[]byte) error {
		*out[0] = []byte{}
		*out[1] = []byte{}
		req := &messages.Trigger{}
		err := json.Unmarshal(in.([]byte), req)
		if err != nil {
			*out[1] = []byte(err.Error())
			return err
		}

		rsp, err := hdl.External(req)
		if err != nil {
			*out[1] = []byte(err.Error())
			return err
		}

		bts, err := json.Marshal(rsp)
		if err != nil {
			*out[1] = []byte(err.Error())
			return err
		}

		*out[0] = bts
		return nil
	})

	handler.Add(serviceNmsp+".NoParams", func(in interface{}, out ...*[]byte) error {
		*out[0] = []byte{}
		*out[1] = []byte{}
		fmt.Println("RECEIVED HOOK")

		rsp, err := hdl.NoParams()
		if err != nil {
			*out[1] = []byte(err.Error())
			return err
		}

		bts, err := json.Marshal(rsp)
		if err != nil {
			*out[1] = []byte(err.Error())
			return err
		}

		*out[0] = bts
		return nil
	})

	replier.AddHandler(handler)
	replier.AddTransport(tcp.NewTransport())
	replier.AddTransport(inproc.NewTransport())

	return replier, nil
}
