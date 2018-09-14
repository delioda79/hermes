package example2

import (
	"encoding/json"

	"bitbucket.org/ConsentSystems/mango-micro/handler"
	"bitbucket.org/ConsentSystems/mango-micro/replier"
	"nanomsg.org/go-mangos/transport/inproc"
	"nanomsg.org/go-mangos/transport/tcp"
)

// NewAPICallsHandlerServer returns a new replier server
func NewAPICallsHandlerServer(
	regAddr string,
	serverPort int,
	hdl APICallsHandler,
) {

	replier, _ := replier.NewServer(regAddr, "APICallsHandlerServer-replier", "1")
	handler := handler.NewHandler()

	handler.Add("APICallsHandler.TestBool", func(in interface{}, out ...*[]byte) error {
		*out[0] = []byte{}
		*out[1] = []byte{}
		req := &APICallMessage{}
		err := json.Unmarshal(in.([]byte), req)
		if err != nil {
			*out[1] = []byte(err.Error())
			return err
		}

		rsp, err := hdl.TestBool(req)
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
	go replier.Run(serverPort, "inproc", "APICallsHandlerServer-replier")
	go replier.Run(serverPort, "tcp", "")
}
